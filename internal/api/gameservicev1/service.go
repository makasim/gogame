package gameservicev1

import (
	"context"

	"connectrpc.com/connect"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
	"github.com/makasim/gogame/protogen/gogame/v1/gogamev1connect"
)

var _ gogamev1connect.GameServiceHandler = (*Service)(nil)

type streamVacantGamesHandler interface {
	StreamVacantGames(ctx context.Context, req *connect.Request[v1.StreamVacantGamesRequest], stream *connect.ServerStream[v1.StreamVacantGamesResponse]) error
}

type streamGameEventsHandler interface {
	StreamGameEvents(context.Context, *connect.Request[v1.StreamGameEventsRequest], *connect.ServerStream[v1.StreamGameEventsResponse]) error
}

type Service struct {
	svgh streamVacantGamesHandler
	sgeh streamGameEventsHandler
}

func New(
	svgh streamVacantGamesHandler,
	sgeh streamGameEventsHandler,
) *Service {
	return &Service{
		svgh: svgh,
		sgeh: sgeh,
	}
}

func (s *Service) CreateGame(ctx context.Context, req *connect.Request[v1.CreateGameRequest]) (*connect.Response[v1.CreateGameResponse], error) {
	panic("BUG: CreateGame must not be called")
}

func (s *Service) JoinGame(ctx context.Context, req *connect.Request[v1.JoinGameRequest]) (*connect.Response[v1.JoinGameResponse], error) {
	panic("BUG: JoinGame must not be called")
}

func (s *Service) StreamVacantGames(ctx context.Context, req *connect.Request[v1.StreamVacantGamesRequest], stream *connect.ServerStream[v1.StreamVacantGamesResponse]) error {
	return s.svgh.StreamVacantGames(ctx, req, stream)
}

func (s *Service) StreamGameEvents(ctx context.Context, req *connect.Request[v1.StreamGameEventsRequest], stream *connect.ServerStream[v1.StreamGameEventsResponse]) error {
	return s.sgeh.StreamGameEvents(ctx, req, stream)
}

func (s *Service) MakeMove(ctx context.Context, req *connect.Request[v1.MakeMoveRequest]) (*connect.Response[v1.MakeMoveResponse], error) {
	panic("BUG: MakeMove must not be called")
}

func (s *Service) Pass(ctx context.Context, req *connect.Request[v1.PassRequest]) (*connect.Response[v1.PassResponse], error) {
	panic("BUG: Pass must not be called")
}

func (s *Service) Resign(ctx context.Context, req *connect.Request[v1.ResignRequest]) (*connect.Response[v1.ResignResponse], error) {
	panic("BUG: Resign must not be called")
}

func (s *Service) Undo(ctx context.Context, req *connect.Request[v1.UndoRequest]) (*connect.Response[v1.UndoResponse], error) {
	panic("BUG: Undo must not be called")
}
