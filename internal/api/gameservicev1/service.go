package gameservicev1

import (
	"context"

	"connectrpc.com/connect"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
)

type createGameHandler interface {
	CreateGame(_ context.Context, req *connect.Request[v1.CreateGameRequest]) (*connect.Response[v1.CreateGameResponse], error)
}

type joinGameHandler interface {
	JoinGame(_ context.Context, req *connect.Request[v1.JoinGameRequest]) (*connect.Response[v1.JoinGameResponse], error)
}

type Service struct {
	cgh createGameHandler
	jgh joinGameHandler
}

func New(cgh createGameHandler, jgh joinGameHandler) *Service {
	return &Service{
		cgh: cgh,
		jgh: jgh,
	}
}

func (s *Service) CreateGame(ctx context.Context, req *connect.Request[v1.CreateGameRequest]) (*connect.Response[v1.CreateGameResponse], error) {
	return s.cgh.CreateGame(ctx, req)
}

func (s *Service) JoinGame(ctx context.Context, req *connect.Request[v1.JoinGameRequest]) (*connect.Response[v1.JoinGameResponse], error) {
	return s.jgh.JoinGame(ctx, req)
}
