package streamgameeventshandler

import (
	"context"
	"log"
	"strconv"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
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

func (h *Handler) StreamGameEvents(ctx context.Context, req *connect.Request[v1.StreamGameEventsRequest], stream *connect.ServerStream[v1.StreamGameEventsResponse]) error {
	_, _, _, err := convertor.FindGame(h.e, req.Msg.GameId, 0)
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	wCmd := flowstate.Watch(map[string]string{
		`game.id`: req.Msg.GameId,
	}).WithORLabels(map[string]string{
		`undo.game.id`: req.Msg.GameId,
	}).WithSinceLatest()

	if err := h.e.Do(wCmd); err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	lis := wCmd.Listener
	defer lis.Close()

	for {
		select {
		case state := <-lis.Listen():
			if state.Labels[`undo.game.id`] != `` {
				gID := state.Annotations[`game.id`]
				gRev, _ := strconv.ParseInt(state.Annotations[`game.rev`], 10, 0)

				d := &flowstate.Data{}
				stateCtx := &flowstate.StateCtx{}

				undoStateCtx := state.CopyToCtx(&flowstate.StateCtx{})
				undoD := &flowstate.Data{}
				if err := h.e.Do(
					flowstate.GetByID(stateCtx, flowstate.StateID(gID), gRev),
					flowstate.DereferenceData(stateCtx, d, `game`),
					flowstate.GetData(d),
					flowstate.DereferenceData(undoStateCtx, undoD, `undo`),
					flowstate.GetData(undoD),
				); err != nil {
					log.Printf("stream: 1: %s", err)
					continue
				}

				u, err := convertor.DataToUndo(undoD)
				if err != nil {
					log.Printf("stream: 2: %s", err)
					continue
				}
				g, err := convertor.DataToGame(d)
				if err != nil {
					log.Printf("stream: 3: %s", err)
					continue
				}

				if err := stream.Send(&v1.StreamGameEventsResponse{
					Game: g,
					Undo: u,
				}); err != nil {
					log.Printf("stream: 4: %s", err)
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
