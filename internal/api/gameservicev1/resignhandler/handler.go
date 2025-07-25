package resignhandler

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
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

	if stateCtx.Current.Labels[`game.state`] == `ended` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game has already ended"))
	}

	stateCtx.Current.SetLabel(`game.state`, `ended`)
	g.State = v1.State_STATE_ENDED
	g.Winner = convertor.NextPlayer(g)
	g.WonBy = `resign`

	if err = convertor.GameToData(g, d); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := h.e.Do(flowstate.Commit(
		flowstate.AttachData(stateCtx, d, `game`),
		flowstate.Park(stateCtx),
	)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	g.Rev = int32(stateCtx.Current.Rev)

	return connect.NewResponse(&v1.ResignResponse{
		Game: g,
	}), nil
}
