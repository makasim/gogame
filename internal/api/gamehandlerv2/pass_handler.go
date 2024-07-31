package gamehandlerv2

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/endedflow"
	"github.com/makasim/gogame/internal/moveflow"
	v2 "github.com/makasim/gogame/protogen/gogame/v2"
)

type PassHandler struct {
	e *flowstate.Engine
}

func NewPassHandler(e *flowstate.Engine) *PassHandler {
	return &PassHandler{
		e: e,
	}
}

func (h *PassHandler) Pass(_ context.Context, req *connect.Request[v2.PassRequest]) (*connect.Response[v2.PassResponse], error) {
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
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("not player's turn"))
	}

	g.Changes = append(g.Changes, &v2.Change{
		Change: &v2.Change_Pass_{
			Pass: &v2.Change_Pass{
				PlayerId: req.Msg.PlayerId,
			},
		},
	})

	if len(g.Changes) > 1 && g.Changes[len(g.Changes)-2].GetPass() != nil && g.Changes[len(g.Changes)-2].GetPass().PlayerId == lp.Id {
		stateCtx.Current.SetLabel(`game.state`, `ended`)
		// g.State = v1.State_STATE_ENDED

		g.Changes = append(g.Changes, &v2.Change{
			Change: &v2.Change_End_{
				End: &v2.Change_End{
					// TODO: add decide on winner algorithm
					Draw: true,
				},
			},
		})

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

		g.Rev = int32(stateCtx.Current.Rev)

		return connect.NewResponse(&v2.PassResponse{
			Game: g,
		}), nil
	}

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
