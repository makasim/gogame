package joingamehandler

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/createdflow"
	"github.com/makasim/gogame/internal/moveflow"
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

func (h *Handler) JoinGame(_ context.Context, req *connect.Request[v1.JoinGameRequest]) (*connect.Response[v1.JoinGameResponse], error) {
	if req.Msg.GameId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game id is required"))
	}
	if req.Msg.Player2.Name == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("player2 name is required"))
	}

	g, stateCtx, d, err := convertor.FindGame(h.e, req.Msg.GameId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if stateCtx.Current.Transition.ToID != createdflow.ID {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game is not joinable"))
	}
	stateCtx.Current.SetLabel(`game.state`, `started`)

	g.Player2 = req.Msg.Player2
	if time.Now().UnixNano()%2 == 0 {
		g.Player1.Color = v1.Color_COLOR_BLACK
		g.Player2.Color = v1.Color_COLOR_WHITE
	} else {
		g.Player1.Color = v1.Color_COLOR_WHITE
		g.Player2.Color = v1.Color_COLOR_BLACK
	}

	if err = convertor.GameToData(g, d); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := h.e.Do(flowstate.Commit(
		flowstate.StoreData(d),
		flowstate.ReferenceData(stateCtx, d, `game`),
		flowstate.Pause(stateCtx).WithTransit(moveflow.ID),
	)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.JoinGameResponse{
		Game: g,
	}), nil
}
