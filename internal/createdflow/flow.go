package createdflow

import (
	"fmt"

	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/endedflow"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
)

var ID flowstate.TransitionID = `created`

type Flow struct {
}

func New() (flowstate.TransitionID, *Flow) {
	return ID, &Flow{}
}

func (f *Flow) Execute(stateCtx *flowstate.StateCtx, e flowstate.Engine) (flowstate.Command, error) {
	if !flowstate.Delayed(stateCtx.Current) {
		return nil, fmt.Errorf("a flow should not be executed")
	}

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

	stateCtx.Current.Transition = flowstate.Transition{
		From: stateCtx.Current.Transition.From,
		To:   endedflow.ID,
	}

	if err := e.Do(flowstate.Commit(
		flowstate.AttachData(stateCtx, d, `game`),
		flowstate.End(stateCtx),
	)); flowstate.IsErrRevMismatch(err) {
		return flowstate.Noop(stateCtx), nil
	} else if err != nil {
		return nil, err
	}

	return flowstate.Noop(stateCtx), nil
}
