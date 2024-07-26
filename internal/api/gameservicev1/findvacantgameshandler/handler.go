package findvacantgameshandler

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

func (h *Handler) FindVacantGames(ctx context.Context, _ *connect.Request[v1.FindVacantGamesRequest], stream *connect.ServerStream[v1.FindVacantGamesResponse]) error {
	wCmd := flowstate.Watch(map[string]string{
		`game.state`: `created`,
	}).WithORLabels(map[string]string{
		`game.state`: `started`,
	})

	if err := h.e.Do(wCmd); err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	lis := wCmd.Listener
	defer lis.Close()

	for {
		select {
		case state := <-lis.Listen():
			g, _, _, err := convertor.FindGame(h.e, state.Labels[`game.id`])
			if err != nil {
				return connect.NewError(connect.CodeInternal, err)
			}

			if err := stream.Send(&v1.FindVacantGamesResponse{
				Game:     g,
				Joinable: state.Labels[`game.state`] == `created`,
			}); err != nil {
				return connect.NewError(connect.CodeInternal, err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
