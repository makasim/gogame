package gamehandlerv2

import (
	"context"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	v2 "github.com/makasim/gogame/protogen/gogame/v2"
)

type StreamEventsHandler struct {
	e *flowstate.Engine
}

func NewStreamEventsHandler(e *flowstate.Engine) *StreamEventsHandler {
	return &StreamEventsHandler{
		e: e,
	}
}

func (h *StreamEventsHandler) StreamEvents(ctx context.Context, req *connect.Request[v2.StreamEventsRequest], stream *connect.ServerStream[v2.StreamEventsResponse]) error {
	_, _, _, err := convertor.FindGame(h.e, req.Msg.GameId, 0)
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	wCmd := flowstate.Watch(map[string]string{
		`game.id`: req.Msg.GameId,
	}).WithSinceLatest()

	if err := h.e.Do(wCmd); err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	lis := wCmd.Listener
	defer lis.Close()

	for {
		select {
		case state := <-lis.Listen():
			g, _, _, err := convertor.FindGame(h.e, state.Labels[`game.id`], int32(state.Rev))
			if err != nil {
				return connect.NewError(connect.CodeInternal, err)
			}

			if err := stream.Send(&v2.StreamEventsResponse{
				Game: g,
			}); err != nil {
				return connect.NewError(connect.CodeInternal, err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
