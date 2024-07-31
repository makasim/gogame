package gamehandlerv2

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/endedflow"
	v2 "github.com/makasim/gogame/protogen/gogame/v2"
)

type ResignHandler struct {
	e *flowstate.Engine
}

func NewResignHandler(e *flowstate.Engine) *ResignHandler {
	return &ResignHandler{
		e: e,
	}
}

func (h *ResignHandler) Resign(_ context.Context, req *connect.Request[v2.ResignRequest]) (*connect.Response[v2.ResignResponse], error) {
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
	g.Changes = append(g.Changes, &v2.Change{
		Change: &v2.Change_End_{
			End: &v2.Change_End{
				Draw:   false,
				Winner: convertor.AnotherPlayer(g, req.Msg.PlayerId),
				WonBy:  "resign",
			},
		},
	})

	//g.State = v1.State_STATE_ENDED

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

	return connect.NewResponse(&v2.ResignResponse{
		Game: g,
	}), nil
}
