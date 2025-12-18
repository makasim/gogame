package passflow

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

var ID flowstate.FlowID = `gogame.pass`

type Flow struct {
}

func New() (flowstate.FlowID, *Flow) {
	return ID, &Flow{}
}

func (f *Flow) Execute(reqStateCtx *flowstate.StateCtx, e *flowstate.Engine) (flowstate.Command, error) {
	msg := &v1.PassRequest{}
	if err := promutil.UnmarshalRequest(reqStateCtx, msg); err != nil {
		return nil, err
	}
	if msg.GameId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game id is required"))
	}
	if msg.GameRev == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game rev is required"))
	}
	if msg.PlayerId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("move player id is required"))
	}

	g, stateCtx, d, err := convertor.FindGame(e, msg.GameId, msg.GameRev)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if !(stateCtx.Current.Labels[`game.state`] == `started` || stateCtx.Current.Labels[`game.state`] == `move`) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("state is not move"))
	}
	if g.CurrentMove.PlayerId != msg.PlayerId {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("not player's turn"))
	}

	g.PreviousMoves = append(g.PreviousMoves, &v1.Move{
		PlayerId: g.CurrentMove.PlayerId,
		Color:    g.CurrentMove.Color,
		Pass:     true,
	})

	if len(g.PreviousMoves) > 1 && g.PreviousMoves[len(g.PreviousMoves)-2].Pass {
		stateCtx.Current.SetLabel(`game.state`, `ended`)
		g.State = v1.State_STATE_ENDED

		// TODO: add decide on winner algorithm
		g.Winner = convertor.CurrentPlayer(g)
		g.WonBy = `score`

		if err = convertor.GameToData(g, d); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		if err := e.Do(flowstate.Commit(
			flowstate.StoreData(stateCtx, `game`),
			flowstate.Park(stateCtx),
		)); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		g.Rev = int32(stateCtx.Current.Rev)

		return flowstate.Noop(), promutil.MarshalResponse(reqStateCtx, &v1.PassResponse{
			Game: g,
		})
	}

	g.State = v1.State_STATE_MOVE
	stateCtx.Current.SetLabel(`game.state`, `move`)
	g.CurrentMove = &v1.Move{
		PlayerId: convertor.NextPlayer(g).Id,
		Color:    convertor.NextColor(g),
		EndAt:    time.Now().Add(time.Duration(g.MoveDurationSec) * time.Second).Unix(),
	}

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

	return flowstate.Noop(), promutil.MarshalResponse(reqStateCtx, &v1.PassResponse{
		Game: g,
	})
}
