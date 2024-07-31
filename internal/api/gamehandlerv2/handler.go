package gamehandlerv2

import (
	"context"

	"connectrpc.com/connect"
	v2 "github.com/makasim/gogame/protogen/gogame/v2"
	"github.com/makasim/gogame/protogen/gogame/v2/gogamev2connect"
)

var _ gogamev2connect.GameServiceHandler = (*Service)(nil)

type resignHandler interface {
	Resign(context.Context, *connect.Request[v2.ResignRequest]) (*connect.Response[v2.ResignResponse], error)
}

type Service struct {
	seh *StreamEventsHandler
	mh  *MoveHandler
	ph  *PassHandler
	uh  *UndoHandler
	rh  *ResignHandler
}

func New(
	seh *StreamEventsHandler,
	mh *MoveHandler,
	ph *PassHandler,
	uh *UndoHandler,
	rh *ResignHandler,
) *Service {
	return &Service{
		seh: seh,
		mh:  mh,
		rh:  rh,
		uh:  uh,
		ph:  ph,
	}
}

func (s *Service) StreamEvents(ctx context.Context, req *connect.Request[v2.StreamEventsRequest], stream *connect.ServerStream[v2.StreamEventsResponse]) error {
	return s.seh.StreamEvents(ctx, req, stream)
}

func (s *Service) Move(ctx context.Context, req *connect.Request[v2.MoveRequest]) (*connect.Response[v2.MoveResponse], error) {
	return s.mh.Move(ctx, req)
}

func (s *Service) Pass(ctx context.Context, req *connect.Request[v2.PassRequest]) (*connect.Response[v2.PassResponse], error) {
	return s.ph.Pass(ctx, req)
}

func (s *Service) Undo(ctx context.Context, req *connect.Request[v2.UndoRequest]) (*connect.Response[v2.UndoResponse], error) {
	return s.uh.Undo(ctx, req)
}

func (s *Service) Resign(ctx context.Context, req *connect.Request[v2.ResignRequest]) (*connect.Response[v2.ResignResponse], error) {
	return s.rh.Resign(ctx, req)
}
