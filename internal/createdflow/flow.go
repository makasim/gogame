package createdflow

import (
	"fmt"

	"github.com/makasim/flowstate"
)

var ID flowstate.FlowID = `created`

type Flow struct {
}

func New() (flowstate.FlowID, *Flow) {
	return ID, &Flow{}
}

func (f *Flow) Execute(stateCtx *flowstate.StateCtx, _ *flowstate.Engine) (flowstate.Command, error) {
	return nil, fmt.Errorf("a flow should not be executed")
}
