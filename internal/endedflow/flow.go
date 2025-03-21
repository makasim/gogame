package endedflow

import (
	"fmt"

	"github.com/makasim/flowstate"
)

var ID flowstate.FlowID = `ended`

type Flow struct {
}

func New() (flowstate.FlowID, *Flow) {
	return ID, &Flow{}
}

func (f *Flow) Execute(_ *flowstate.StateCtx, _ flowstate.Engine) (flowstate.Command, error) {
	return nil, fmt.Errorf("a flow should not be executed")
}
