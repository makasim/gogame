package streamgameeventshandler

import (
	"context"
	"strconv"
	"time"

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

func (h *Handler) StreamGameEvents(ctx context.Context, req *connect.Request[v1.StreamGameEventsRequest], stream *connect.ServerStream[v1.StreamGameEventsResponse]) error {
	_, _, _, err := convertor.FindGame(h.e, req.Msg.GameId, 0)
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	getManyCmd := flowstate.GetStatesByLabels(map[string]string{
		`game.id`: req.Msg.GameId,
	}).WithORLabels(map[string]string{
		`undo.game.id`: req.Msg.GameId,
	}).WithSinceLatest()

	w := flowstate.NewWatcher(h.e, time.Millisecond*500, getManyCmd)
	defer w.Close()

	for {
		select {
		case state := <-w.Next():
			if state.Labels[`undo.game.id`] != `` {
				gID := state.Annotations[`game.id`]
				gRev, _ := strconv.ParseInt(state.Annotations[`game.rev`], 10, 0)

				d := &flowstate.Data{}
				stateCtx := &flowstate.StateCtx{}

				undoStateCtx := state.CopyToCtx(&flowstate.StateCtx{})
				undoD := &flowstate.Data{}
				if err := h.e.Do(
					flowstate.GetStateByID(stateCtx, flowstate.StateID(gID), gRev),
					flowstate.GetData(stateCtx, d, `game`),
					flowstate.GetData(undoStateCtx, undoD, `undo`),
				); err != nil {
					continue
				}

				u, err := convertor.DataToUndo(undoD)
				if err != nil {
					continue
				}
				g, err := convertor.DataToGame(d)
				if err != nil {
					continue
				}

				if err := stream.Send(&v1.StreamGameEventsResponse{
					Game: g,
					Undo: u,
				}); err != nil {
					continue
				}
				continue
			}

			g, _, _, err := convertor.FindGame(h.e, state.Labels[`game.id`], int32(state.Rev))
			if err != nil {
				return connect.NewError(connect.CodeInternal, err)
			}

			if err := stream.Send(&v1.StreamGameEventsResponse{
				Game: g,
			}); err != nil {
				return connect.NewError(connect.CodeInternal, err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
