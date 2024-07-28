package moveflow

import (
	"fmt"

	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/endedflow"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
)

var ID flowstate.FlowID = `move`

type Flow struct {
}

func New() (flowstate.FlowID, *Flow) {
	return ID, &Flow{}
}

func (f *Flow) Execute(stateCtx *flowstate.StateCtx, e *flowstate.Engine) (flowstate.Command, error) {
	if flowstate.Delayed(stateCtx.Current) {
		d := &flowstate.Data{}

		if err := e.Do(
			flowstate.DereferenceData(stateCtx, d, `game`),
			flowstate.GetData(d),
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

		if err := e.Do(flowstate.Commit(
			flowstate.StoreData(d),
			flowstate.ReferenceData(stateCtx, d, `game`),
			flowstate.Pause(stateCtx).WithTransit(endedflow.ID),
		)); err != nil {
			return nil, err
		}

		return flowstate.Noop(stateCtx), nil
	}

	return nil, fmt.Errorf("a flow should not be executed")
}
