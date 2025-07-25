package creategamehandler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/staleflow"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
)

type Handler struct {
	e flowstate.Engine
}

func New(e flowstate.Engine) *Handler {
	return &Handler{
		e: e,
	}
}

func (h *Handler) CreateGame(_ context.Context, req *connect.Request[v1.CreateGameRequest]) (*connect.Response[v1.CreateGameResponse], error) {
	if req.Msg.Name == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game name is required"))
	}
	if req.Msg.Player1 != nil && req.Msg.Player1.Name == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("player1 name is required"))
	}

	g := &v1.Game{
		Id:              strconv.FormatInt(time.Now().UnixNano(), 10),
		Name:            req.Msg.Name,
		Player1:         req.Msg.Player1,
		State:           v1.State_STATE_CREATED,
		MoveDurationSec: req.Msg.MoveDurationSec,
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
	}

	if err := h.e.Do(flowstate.Commit(
		flowstate.AttachData(stateCtx, d, `game`),
		flowstate.Park(stateCtx),
		flowstate.Delay(stateCtx, staleflow.ID, time.Minute),
	)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	g.Rev = int32(stateCtx.Current.Rev)

	return connect.NewResponse(&v1.CreateGameResponse{
		Game: g,
	}), nil
}
