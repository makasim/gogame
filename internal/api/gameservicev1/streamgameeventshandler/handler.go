package streamgameeventshandler

import (
	"context"

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
