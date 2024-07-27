package streamvacantgameshandler

import (
	"context"
	"time"

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

func (h *Handler) StreamVacantGames(ctx context.Context, _ *connect.Request[v1.StreamVacantGamesRequest], stream *connect.ServerStream[v1.StreamVacantGamesResponse]) error {
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

			if err := stream.Send(&v1.StreamVacantGamesResponse{
				Game: g,
			}); err != nil {
				return connect.NewError(connect.CodeInternal, err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
