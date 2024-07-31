package roomhandlerv2

import (
	"context"

	"connectrpc.com/connect"
	v2 "github.com/makasim/gogame/protogen/gogame/v2"
	"github.com/makasim/gogame/protogen/gogame/v2/gogamev2connect"
)

var _ gogamev2connect.RoomServiceHandler = (*Handler)(nil)

type Handler struct {
	cgh *CreateGameHandler
	jgh *JoinGameHandler
	sgh *StreamGamesHandler
}

func New(
	cgh *CreateGameHandler,
	jgh *JoinGameHandler,
	svh *StreamGamesHandler,
) *Handler {
	return &Handler{
		cgh: cgh,
		jgh: jgh,
		sgh: svh,
	}
}

func (s *Handler) CreateGame(ctx context.Context, req *connect.Request[v2.CreateGameRequest]) (*connect.Response[v2.CreateGameResponse], error) {
	return s.cgh.CreateGame(ctx, req)
}

func (s *Handler) JoinGame(ctx context.Context, req *connect.Request[v2.JoinGameRequest]) (*connect.Response[v2.JoinGameResponse], error) {
	return s.jgh.JoinGame(ctx, req)
}

func (s *Handler) StreamGames(ctx context.Context, req *connect.Request[v2.StreamGamesRequest], stream *connect.ServerStream[v2.StreamGamesResponse]) error {
	return s.sgh.StreamGames(ctx, req, stream)
}
