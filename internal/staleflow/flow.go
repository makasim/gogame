package staleflow

import (
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/endedflow"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
)

var ID flowstate.TransitionID = `stale`

type Flow struct {
}

func New() (flowstate.TransitionID, *Flow) {
	return ID, &Flow{}
}

func (f *Flow) Execute(stateCtx *flowstate.StateCtx, e flowstate.Engine) (flowstate.Command, error) {
	d := &flowstate.Data{}
	if err := e.Do(flowstate.GetData(stateCtx, d, `game`)); err != nil {
		return nil, err
	}

	g, err := convertor.DataToGame(d)
	if err != nil {
		return nil, err
	}

	g.Rev = int32(stateCtx.Current.Rev)

	stateCtx.Current.SetLabel(`game.state`, `ended`)
	g.State = v1.State_STATE_ENDED
	g.WonBy = `not_started`

	if err = convertor.GameToData(g, d); err != nil {
		return nil, err
	}

	if err := e.Do(flowstate.Commit(
		flowstate.AttachData(stateCtx, d, `game`),
		flowstate.End(stateCtx).WithTransit(endedflow.ID),
	)); flowstate.IsErrRevMismatch(err) {
		return flowstate.Noop(stateCtx), nil
	} else if err != nil {
		return nil, err
	}

	return flowstate.Noop(stateCtx), nil
}
