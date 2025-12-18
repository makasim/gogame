package makemoveflow

import (
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/movetimeoutflow"
	"github.com/makasim/gogame/internal/promutil"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
)

var ID flowstate.FlowID = `gogame.make_move`

type Flow struct {
}

func New() (flowstate.FlowID, *Flow) {
	return ID, &Flow{}
}

func (f *Flow) Execute(reqStateCtx *flowstate.StateCtx, e *flowstate.Engine) (flowstate.Command, error) {
	msg := &v1.MakeMoveRequest{}
	if err := promutil.UnmarshalRequest(reqStateCtx, msg); err != nil {
		return nil, err
	}
	if msg.GameId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game id is required"))
	}
	if msg.GameRev == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game rev is required"))
	}
	if msg.Move == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("move is required"))
	}
	if msg.Move.PlayerId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("move player id is required"))
	}
	if msg.Move.Color <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("move color is required"))
	}
	if msg.Move.X < 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("move x is required"))
	}
	if msg.Move.Y < 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("move y is required"))
	}

	g, stateCtx, d, err := convertor.FindGame(e, msg.GameId, msg.GameRev)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if !(stateCtx.Current.Labels[`game.state`] == `started` || stateCtx.Current.Labels[`game.state`] == `move`) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("state is not move"))
	}
	if g.CurrentMove.PlayerId != msg.Move.PlayerId {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("not player's turn"))
	}

	b, err := convertor.GameToBoard(g)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	nextMove := &v1.Move{
		PlayerId: g.CurrentMove.PlayerId,
		Color:    g.CurrentMove.Color,
		X:        msg.Move.X,
		Y:        msg.Move.Y,
	}

	l, err := b.PlaceStone(convertor.ToClamMove(nextMove))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	convertor.CurrentPlayer(g).CapturedStones += int32(len(l))

	g.State = v1.State_STATE_MOVE
	stateCtx.Current.SetLabel(`game.state`, `move`)

	g.PreviousMoves = append(g.PreviousMoves, nextMove)
	g.CurrentMove = &v1.Move{
		PlayerId: convertor.NextPlayer(g).Id,
		Color:    convertor.NextColor(g),
		EndAt:    time.Now().Add(time.Duration(g.MoveDurationSec) * time.Second).Unix(),
	}
	g.Board = convertor.FromClamBoard(b)

	if err = convertor.GameToData(g, d); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := e.Do(flowstate.Commit(
		flowstate.StoreData(stateCtx, `game`),
		flowstate.Park(stateCtx),
		flowstate.Delay(stateCtx, movetimeoutflow.ID, time.Duration(g.MoveDurationSec)*time.Second),
	)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	g.Rev = int32(stateCtx.Current.Rev)

	return flowstate.Noop(), promutil.MarshalResponse(reqStateCtx, &v1.MakeMoveResponse{
		Game: g,
	})
}
