package undoflow

import (
	"fmt"

	"github.com/makasim/flowstate"
)

var ID flowstate.FlowID = `undo`

type Flow struct {
}

func New() (flowstate.FlowID, *Flow) {
	return ID, &Flow{}
}

func (f *Flow) Execute(_ *flowstate.StateCtx, e flowstate.Engine) (flowstate.Command, error) {
	return nil, fmt.Errorf("a flow should not be executed")
}
