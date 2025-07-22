package endedflow

import (
	"fmt"

	"github.com/makasim/flowstate"
)

var ID flowstate.TransitionID = `ended`

type Flow struct {
}

func New() (flowstate.TransitionID, *Flow) {
	return ID, &Flow{}
}

func (f *Flow) Execute(_ *flowstate.StateCtx, _ flowstate.Engine) (flowstate.Command, error) {
	return nil, fmt.Errorf("a flow should not be executed")
}
