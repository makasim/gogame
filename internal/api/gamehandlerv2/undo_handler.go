package gamehandlerv2

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/moveflow"
	v2 "github.com/makasim/gogame/protogen/gogame/v2"
)

type UndoHandler struct {
	e *flowstate.Engine
}

func NewUndoHandler(e *flowstate.Engine) *UndoHandler {
	return &UndoHandler{
		e: e,
	}
}

func (h *UndoHandler) Undo(_ context.Context, req *connect.Request[v2.UndoRequest]) (*connect.Response[v2.UndoResponse], error) {
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
	lp := convertor.LastPlayer(g)
	if lp == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game has no previous move"))
	}
	if lp.Id == req.Msg.PlayerId {
		if !req.Msg.Accept {
			return connect.NewResponse(&v2.UndoResponse{
				Game: g,
			}), nil
		}

		if len(g.Changes) > 0 && g.Changes[len(g.Changes)-1].GetUndo() == nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("last change not undo"))
		}

		g.Changes = append(g.Changes, &v2.Change{
			Change: &v2.Change_Undo_{
				Undo: &v2.Change_Undo{
					PlayerId: req.Msg.PlayerId,
				},
			},
		})

		return nil, nil
	}

	g.Changes = append(g.Changes, &v2.Change{
		Change: &v2.Change_Undo_{
			Undo: &v2.Change_Undo{
				PlayerId: req.Msg.PlayerId,
			},
		},
	})

	//g.State = v1.State_STATE_MOVE
	stateCtx.Current.SetLabel(`game.state`, `move`)

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

	return connect.NewResponse(&v2.PassResponse{
		Game: g,
	}), nil
}
