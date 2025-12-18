package resignflow

import (
	"fmt"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/promutil"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
)

var ID flowstate.FlowID = `gogame.resign`

type Flow struct {
}

func New() (flowstate.FlowID, *Flow) {
	return ID, &Flow{}
}

func (f *Flow) Execute(reqStateCtx *flowstate.StateCtx, e *flowstate.Engine) (flowstate.Command, error) {
	msg := &v1.ResignRequest{}
	if err := promutil.UnmarshalRequest(reqStateCtx, msg); err != nil {
		return nil, err
	}
	if msg.GameId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game id is required"))
	}
	if msg.PlayerId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("player id is required"))
	}

	g, stateCtx, d, err := convertor.FindGame(e, msg.GameId, 0)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if stateCtx.Current.Labels[`game.state`] == `ended` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game has already ended"))
	}

	stateCtx.Current.SetLabel(`game.state`, `ended`)
	g.State = v1.State_STATE_ENDED
	g.Winner = convertor.NextPlayer(g)
	g.WonBy = `resign`

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

	return flowstate.Noop(), promutil.MarshalResponse(reqStateCtx, &v1.ResignResponse{
		Game: g,
	})
}
