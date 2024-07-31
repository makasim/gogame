package roomhandlerv2

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/createdflow"
	v2 "github.com/makasim/gogame/protogen/gogame/v2"
)

type CreateGameHandler struct {
	e *flowstate.Engine
}

func NewCreateGameHandler(e *flowstate.Engine) *CreateGameHandler {
	return &CreateGameHandler{
		e: e,
	}
}

func (h *CreateGameHandler) CreateGame(_ context.Context, req *connect.Request[v2.CreateGameRequest]) (*connect.Response[v2.CreateGameResponse], error) {
	if req.Msg.Name == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game name is required"))
	}
	if req.Msg.Player1 != nil && req.Msg.Player1.Name == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("player1 name is required"))
	}

	g := &v2.Game{
		Id:      strconv.FormatInt(time.Now().UnixNano(), 10),
		Name:    req.Msg.Name,
		Player1: req.Msg.Player1,
		// State:   v2.State_STATE_CREATED,
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
		flowstate.StoreData(d),
		flowstate.ReferenceData(stateCtx, d, `game`),
		flowstate.Pause(stateCtx).WithTransit(createdflow.ID),
	)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	g.Rev = int32(stateCtx.Current.Rev)

	return connect.NewResponse(&v2.CreateGameResponse{
		Game: g,
	}), nil
}
