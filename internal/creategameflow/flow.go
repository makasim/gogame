package creategameflow

import (
	"fmt"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/promutil"
	"github.com/makasim/gogame/internal/staleflow"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
)

var ID flowstate.FlowID = `gogame.create_game`

type Flow struct {
}

func New() (flowstate.FlowID, *Flow) {
	return ID, &Flow{}
}

func (f *Flow) Execute(reqStateCtx *flowstate.StateCtx, e flowstate.Engine) (flowstate.Command, error) {
	msg := &v1.CreateGameRequest{}
	if err := promutil.UnmarshalRequest(reqStateCtx, msg); err != nil {
		return nil, err
	}
	if msg.Name == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game name is required"))
	}
	if msg.Player1 != nil && msg.Player1.Name == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("player1 name is required"))
	}

	g := &v1.Game{
		Id:              strconv.FormatInt(time.Now().UnixNano(), 10),
		Name:            msg.Name,
		Player1:         msg.Player1,
		State:           v1.State_STATE_CREATED,
		MoveDurationSec: msg.MoveDurationSec,
	}
	if g.MoveDurationSec == 0 {
		g.MoveDurationSec = 60
	}

	d := &flowstate.Data{}
	if err := convertor.GameToData(g, d); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	stateCtx := &flowstate.StateCtx{
		Current: flowstate.State{
			ID: flowstate.StateID(g.Id),
			Labels: map[string]string{
				`game`:       `true`,
				`game.id`:    g.Id,
				`game.state`: `created`,
			},
		},
		Datas: map[string]*flowstate.Data{
			"game": d,
		},
	}

	if err := e.Do(flowstate.Commit(
		flowstate.StoreData(stateCtx, `game`),
		flowstate.Park(stateCtx),
		flowstate.Delay(stateCtx, staleflow.ID, time.Minute),
	)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	g.Rev = int32(stateCtx.Current.Rev)

	return flowstate.Noop(), promutil.MarshalResponse(reqStateCtx, &v1.CreateGameResponse{
		Game: g,
	})
}
