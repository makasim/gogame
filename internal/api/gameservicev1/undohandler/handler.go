package undohandler

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
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

func (h *Handler) Undo(_ context.Context, req *connect.Request[v1.UndoRequest]) (*connect.Response[v1.UndoResponse], error) {
	if req.Msg.GameId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game id is required"))
	}
	if req.Msg.GameRev == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game rev is required"))
	}
	if req.Msg.PlayerId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("move player id is required"))
	}

	g, stateCtx, d, err := convertor.FindGame(h.e, req.Msg.GameId, req.Msg.GameRev)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if stateCtx.Current.Transition.ToID != moveflow.ID {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("state is not move"))
	}

	if req.Msg.Request {
		if g.CurrentMove.PlayerId == req.Msg.PlayerId {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("must be another's player turn"))
		}

		g.Undo = &v1.Game_Undo{
			Requested:         true,
			RequesteePlayerId: req.Msg.PlayerId,
			Decision:          0,
		}
	} else if req.Msg.Decision != 0 {
		if g.CurrentMove.PlayerId != req.Msg.PlayerId {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("must decide current move player"))
		}

		g.Undo = &v1.Game_Undo{
			Requested:         true,
			RequesteePlayerId: req.Msg.PlayerId,
			Decision:          0,
		}
	}

	if err = convertor.GameToData(g, d); err != nil {

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := h.e.Do(flowstate.Commit(
		flowstate.StoreData(d),
		flowstate.ReferenceData(stateCtx, d, `game`),
		flowstate.Pause(stateCtx).WithTransit(moveflow.ID),
		flowstate.Delay(stateCtx, time.Second*30).WithCommit(true),
	)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	g.Rev = int32(stateCtx.Current.Rev)

	return connect.NewResponse(&v1.PassResponse{
		Game: g,
	}), nil
}
