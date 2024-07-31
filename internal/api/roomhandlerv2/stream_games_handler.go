package roomhandlerv2

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	v2 "github.com/makasim/gogame/protogen/gogame/v2"
)

type StreamGamesHandler struct {
	e *flowstate.Engine
}

func NewStreamGamesHandler(e *flowstate.Engine) *StreamGamesHandler {
	return &StreamGamesHandler{
		e: e,
	}
}

func (h *StreamGamesHandler) StreamGames(ctx context.Context, req *connect.Request[v2.StreamGamesRequest], stream *connect.ServerStream[v2.StreamGamesResponse]) error {
	wCmd := flowstate.Watch(map[string]string{
		`game.state`: `created`,
	}).WithORLabels(map[string]string{
		`game.state`: `started`,
	}).WithSinceTime(time.Now().Add(-time.Minute * 5))

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

			if err := stream.Send(&v2.StreamGamesResponse{
				Game: g,
			}); err != nil {
				return connect.NewError(connect.CodeInternal, err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
