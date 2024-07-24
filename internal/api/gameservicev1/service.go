package gameservicev1

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/createdflow"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
)

type Service struct {
	e *flowstate.Engine
}

func New(e *flowstate.Engine) *Service {
	return &Service{
		e: e,
	}
}

func (s *Service) CreateGame(_ context.Context, req *connect.Request[v1.CreateGameRequest]) (*connect.Response[v1.CreateGameResponse], error) {
	if req.Msg.Name == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game name is required"))
	}
	if req.Msg.Player1 != nil && req.Msg.Player1.Name == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("player1 name is required"))
	}

	g := &v1.Game{
		Id:      strconv.FormatInt(time.Now().UnixNano(), 10),
		Name:    req.Msg.Name,
		Player1: req.Msg.Player1,
	}

	d := &flowstate.Data{}
	if err := convertor.GameToData(g, d); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	stateCtx := &flowstate.StateCtx{
		Current: flowstate.State{
			ID: flowstate.StateID(g.Id),
		},
	}

	if err := s.e.Do(flowstate.Commit(
		flowstate.StoreData(d),
		flowstate.ReferenceData(stateCtx, d, `game`),
		flowstate.Pause(stateCtx).WithTransit(createdflow.ID),
	)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.CreateGameResponse{
		Game: g,
	}), nil
}

func (s *Service) JoinGame(_ context.Context, req *connect.Request[v1.JoinGameRequest]) (*connect.Response[v1.JoinGameResponse], error) {
	if req.Msg.GameId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game id is required"))
	}
	if req.Msg.Player2.Name == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("player2 name is required"))
	}

	stateCtx := &flowstate.StateCtx{}
	d := &flowstate.Data{}

	if err := s.e.Do(
		flowstate.GetByID(stateCtx, flowstate.StateID(req.Msg.GameId), 0),
		flowstate.DereferenceData(stateCtx, d, `game`),
		flowstate.GetData(d),
	); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if stateCtx.Current.Transition.ToID != createdflow.ID {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game is not joinable"))
	}

	g, err := convertor.DataToGame(d)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	g.Player2 = req.Msg.Player2
	if time.Now().UnixNano()%2 == 0 {
		g.Player1.Color = v1.Color_COLOR_BLACK
		g.Player2.Color = v1.Color_COLOR_WHITE
	} else {
		g.Player1.Color = v1.Color_COLOR_WHITE
		g.Player2.Color = v1.Color_COLOR_BLACK
	}

	if err = convertor.GameToData(g, d); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := s.e.Do(flowstate.Commit(
		flowstate.StoreData(d),
		flowstate.ReferenceData(stateCtx, d, `game`),
		flowstate.Pause(stateCtx).WithTransit(createdflow.ID),
	)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.JoinGameResponse{
		Game: g,
	}), nil
}
