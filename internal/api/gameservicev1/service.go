package gameservicev1

import (
	"context"

	"connectrpc.com/connect"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
	"github.com/makasim/gogame/protogen/gogame/v1/gogamev1connect"
)

var _ gogamev1connect.GameServiceHandler = (*Service)(nil)

type createGameHandler interface {
	CreateGame(_ context.Context, req *connect.Request[v1.CreateGameRequest]) (*connect.Response[v1.CreateGameResponse], error)
}

type joinGameHandler interface {
	JoinGame(_ context.Context, req *connect.Request[v1.JoinGameRequest]) (*connect.Response[v1.JoinGameResponse], error)
}

type streamVacantGamesHandler interface {
	StreamVacantGames(ctx context.Context, req *connect.Request[v1.StreamVacantGamesRequest], stream *connect.ServerStream[v1.StreamVacantGamesResponse]) error
}

type streamGameEventsHandler interface {
	StreamGameEvents(context.Context, *connect.Request[v1.StreamGameEventsRequest], *connect.ServerStream[v1.StreamGameEventsResponse]) error
}

type makeMoveHandler interface {
	MakeMove(_ context.Context, req *connect.Request[v1.MakeMoveRequest]) (*connect.Response[v1.MakeMoveResponse], error)
}

type resignHandler interface {
	Resign(context.Context, *connect.Request[v1.ResignRequest]) (*connect.Response[v1.ResignResponse], error)
}

type passHandler interface {
	Pass(context.Context, *connect.Request[v1.PassRequest]) (*connect.Response[v1.PassResponse], error)
}

type Service struct {
	cgh  createGameHandler
	jgh  joinGameHandler
	svgh streamVacantGamesHandler
	sgeh streamGameEventsHandler
	mmh  makeMoveHandler
	rh   resignHandler
	ph   passHandler
}

func New(
	cgh createGameHandler,
	jgh joinGameHandler,
	svgh streamVacantGamesHandler,
	sgeh streamGameEventsHandler,
	mmh makeMoveHandler,
	rh resignHandler,
	ph passHandler,
) *Service {
	return &Service{
		cgh:  cgh,
		jgh:  jgh,
		svgh: svgh,
		sgeh: sgeh,
		mmh:  mmh,
		rh:   rh,
		ph:   ph,
	}
}

func (s *Service) CreateGame(ctx context.Context, req *connect.Request[v1.CreateGameRequest]) (*connect.Response[v1.CreateGameResponse], error) {
	return s.cgh.CreateGame(ctx, req)
}

func (s *Service) JoinGame(ctx context.Context, req *connect.Request[v1.JoinGameRequest]) (*connect.Response[v1.JoinGameResponse], error) {
	return s.jgh.JoinGame(ctx, req)
}

func (s *Service) StreamVacantGames(ctx context.Context, req *connect.Request[v1.StreamVacantGamesRequest], stream *connect.ServerStream[v1.StreamVacantGamesResponse]) error {
	return s.svgh.StreamVacantGames(ctx, req, stream)
}

func (s *Service) StreamGameEvents(ctx context.Context, req *connect.Request[v1.StreamGameEventsRequest], stream *connect.ServerStream[v1.StreamGameEventsResponse]) error {
	return s.sgeh.StreamGameEvents(ctx, req, stream)
}

func (s *Service) MakeMove(ctx context.Context, req *connect.Request[v1.MakeMoveRequest]) (*connect.Response[v1.MakeMoveResponse], error) {
	return s.mmh.MakeMove(ctx, req)
}

func (s *Service) Pass(ctx context.Context, req *connect.Request[v1.PassRequest]) (*connect.Response[v1.PassResponse], error) {
	return s.ph.Pass(ctx, req)
}

func (s *Service) Resign(ctx context.Context, req *connect.Request[v1.ResignRequest]) (*connect.Response[v1.ResignResponse], error) {
	return s.rh.Resign(ctx, req)
}
