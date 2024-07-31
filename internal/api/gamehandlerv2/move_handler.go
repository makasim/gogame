package gamehandlerv2

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/moveflow"
	v2 "github.com/makasim/gogame/protogen/gogame/v2"
	"github.com/otrego/clamshell/go/board"
)

type MoveHandler struct {
	e *flowstate.Engine
}

func NewMoveHandler(e *flowstate.Engine) *MoveHandler {
	return &MoveHandler{
		e: e,
	}
}

func (h *MoveHandler) Move(_ context.Context, req *connect.Request[v2.MoveRequest]) (*connect.Response[v2.MoveResponse], error) {
	if req.Msg.GameId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game id is required"))
	}
	if req.Msg.GameRev == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game rev is required"))
	}
	if req.Msg.PlayerId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("move player id is required"))
	}
	if req.Msg.MoveX < 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("move x is required"))
	}
	if req.Msg.MoveY < 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("move y is required"))
	}

	g, stateCtx, d, err := convertor.FindGame(h.e, req.Msg.GameId, req.Msg.GameRev)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	lp := convertor.LastPlayer(g)
	if lp == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game has no previous move"))
	}

	if stateCtx.Current.Transition.ToID != moveflow.ID {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("state is not move"))
	}
	if lp.Id == req.Msg.PlayerId {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("not player's turn"))
	}

	b := board.New(19)
	for _, c := range g.Changes {
		if c.GetMove() == nil {
			continue
		}

		_, err := b.PlaceStone(convertor.ToClamMove(c.GetMove()))
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	nextMove := &v2.Change_Move{
		PlayerId: req.Msg.PlayerId,
		Color:    convertor.NextColor(g),
		X:        req.Msg.MoveX,
		Y:        req.Msg.MoveY,
	}

	l, err := b.PlaceStone(convertor.ToClamMove(nextMove))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	nextMove.CapturedStones += int32(len(l))

	//g.State = v1.State_STATE_MOVE
	stateCtx.Current.SetLabel(`game.state`, `move`)

	g.Changes = append(g.Changes, &v2.Change{
		Change: &v2.Change_Move_{
			Move: nextMove,
		},
	})
	g.Board = convertor.FromClamBoard(b)

	if err = convertor.GameToData(g, d); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := h.e.Do(flowstate.Commit(
		flowstate.StoreData(d),
		flowstate.ReferenceData(stateCtx, d, `game`),
		flowstate.Pause(stateCtx).WithTransit(moveflow.ID),
		flowstate.Delay(stateCtx, time.Second*30).WithCommit(true),
	)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	g.Rev = int32(stateCtx.Current.Rev)

	return connect.NewResponse(&v2.MoveResponse{
		Game: g,
	}), nil
}
