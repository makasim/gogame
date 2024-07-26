package resignhandler

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/endedflow"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
)

type Handler struct {
	e *flowstate.Engine
}

func New(e *flowstate.Engine) *Handler {
	return &Handler{
		e: e,
	}
}

func (h *Handler) Resign(_ context.Context, req *connect.Request[v1.ResignRequest]) (*connect.Response[v1.ResignResponse], error) {
	if req.Msg.GameId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game id is required"))
	}
	if req.Msg.PlayerId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("player id is required"))
	}

	g, stateCtx, d, err := convertor.FindGame(h.e, req.Msg.GameId, 0)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if stateCtx.Current.Transition.ToID == endedflow.ID {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game has already ended"))
	}

	stateCtx.Current.SetLabel(`game.state`, `ended`)
	g.State = `ended`
	// TODO: reason, winner

	if err = convertor.GameToData(g, d); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := h.e.Do(flowstate.Commit(
		flowstate.StoreData(d),
		flowstate.ReferenceData(stateCtx, d, `game`),
		flowstate.Pause(stateCtx).WithTransit(endedflow.ID),
	)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	g.Rev = stateCtx.Current.Rev

	return connect.NewResponse(&v1.ResignResponse{
		Game: g,
	}), nil
}
