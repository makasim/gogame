package movetimeoutflow

import (
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
)

var ID flowstate.FlowID = `movetimeout`

type Flow struct {
}

func New() (flowstate.FlowID, *Flow) {
	return ID, &Flow{}
}

func (f *Flow) Execute(stateCtx *flowstate.StateCtx, e flowstate.Engine) (flowstate.Command, error) {
	d := &flowstate.Data{}

	if err := e.Do(
		flowstate.GetData(stateCtx, d, `game`),
	); err != nil {
		return nil, err
	}

	g, err := convertor.DataToGame(d)
	if err != nil {
		return nil, err
	}

	g.Rev = int32(stateCtx.Current.Rev)

	stateCtx.Current.SetLabel(`game.state`, `ended`)
	g.State = v1.State_STATE_ENDED
	g.Winner = convertor.NextPlayer(g)
	g.WonBy = `timeout`

	if err = convertor.GameToData(g, d); err != nil {
		return nil, err
	}

	return flowstate.Commit(
		flowstate.AttachData(stateCtx, d, `game`),
		flowstate.Park(stateCtx),
	), nil
}
