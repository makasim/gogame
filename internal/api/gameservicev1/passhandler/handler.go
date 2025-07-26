package passhandler

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/movetimeoutflow"
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

func (h *Handler) Pass(_ context.Context, req *connect.Request[v1.PassRequest]) (*connect.Response[v1.PassResponse], error) {
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

	if !(stateCtx.Current.Labels[`game.state`] == `started` || stateCtx.Current.Labels[`game.state`] == `move`) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("state is not move"))
	}
	if g.CurrentMove.PlayerId != req.Msg.PlayerId {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("not player's turn"))
	}

	g.PreviousMoves = append(g.PreviousMoves, &v1.Move{
		PlayerId: g.CurrentMove.PlayerId,
		Color:    g.CurrentMove.Color,
		Pass:     true,
	})

	if len(g.PreviousMoves) > 1 && g.PreviousMoves[len(g.PreviousMoves)-2].Pass {
		stateCtx.Current.SetLabel(`game.state`, `ended`)
		g.State = v1.State_STATE_ENDED

		// TODO: add decide on winner algorithm
		g.Winner = convertor.CurrentPlayer(g)
		g.WonBy = `score`

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

		return connect.NewResponse(&v1.PassResponse{
			Game: g,
		}), nil
	}

	g.State = v1.State_STATE_MOVE
	stateCtx.Current.SetLabel(`game.state`, `move`)
	g.CurrentMove = &v1.Move{
		PlayerId: convertor.NextPlayer(g).Id,
		Color:    convertor.NextColor(g),
		EndAt:    time.Now().Add(time.Duration(g.MoveDurationSec) * time.Second).Unix(),
	}

	if err = convertor.GameToData(g, d); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := h.e.Do(flowstate.Commit(
		flowstate.AttachData(stateCtx, d, `game`),
		flowstate.Park(stateCtx),
		flowstate.Delay(stateCtx, movetimeoutflow.ID, time.Duration(g.MoveDurationSec)*time.Second),
	)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	g.Rev = int32(stateCtx.Current.Rev)

	return connect.NewResponse(&v1.PassResponse{
		Game: g,
	}), nil
}
