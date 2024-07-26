package makemovehandler

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/moveflow"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
)

type Handler struct {
	e *flowstate.Engine
}

func New(e *flowstate.Engine) *Handler {
	return &Handler{
		e: e,
	}
}

func (h *Handler) MakeMove(_ context.Context, req *connect.Request[v1.MakeMoveRequest]) (*connect.Response[v1.MakeMoveResponse], error) {
	if req.Msg.GameId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game id is required"))
	}
	if req.Msg.GameRev == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game rev is required"))
	}
	if req.Msg.Move == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("move is required"))
	}
	if req.Msg.Move.PlayerId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("move player id is required"))
	}
	if req.Msg.Move.Color <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("move color is required"))
	}
	if req.Msg.Move.X <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("move x is required"))
	}
	if req.Msg.Move.Y == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("move y is required"))
	}

	g, stateCtx, d, err := convertor.FindGame(h.e, req.Msg.GameId, req.Msg.GameRev)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if stateCtx.Current.Transition.ToID != moveflow.ID {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("state is not move"))
	}
	if g.CurrentMove.PlayerId != req.Msg.Move.PlayerId {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("not player's turn"))
	}

	g.CurrentMove.X = req.Msg.Move.X
	g.CurrentMove.Y = req.Msg.Move.Y
	g.State = `move`
	stateCtx.Current.SetLabel(`game.state`, `move`)

	g.PreviousMoves = append(g.PreviousMoves, g.CurrentMove)
	g.CurrentMove = &v1.Move{
		PlayerId: convertor.NextPlayer(g).Id,
		Color:    convertor.NextColor(g),
	}

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

	g.Rev = stateCtx.Current.Rev

	return connect.NewResponse(&v1.MakeMoveResponse{
		Game: g,
	}), nil
}
