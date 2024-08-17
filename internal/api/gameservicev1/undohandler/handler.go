package undohandler

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/moveflow"
	"github.com/makasim/gogame/internal/undoflow"
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

func (h *Handler) Undo(_ context.Context, req *connect.Request[v1.UndoRequest]) (*connect.Response[v1.UndoResponse], error) {
	if req.Msg.GameId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game id is required"))
	}
	if req.Msg.GameRev <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game rev is required"))
	}

	g, stateCtx, d, err := convertor.FindGame(h.e, req.Msg.GameId, req.Msg.GameRev)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if len(g.PreviousMoves) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("no moves to undo"))
	}

	m := g.PreviousMoves[len(g.PreviousMoves)-1]

	switch {
	case req.Msg.GetRequest() != nil:
		undoReq := req.Msg.GetRequest()

		if undoReq.PlayerId != m.PlayerId {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("cannot undo other player's move"))
		}

		u := &v1.Undo{
			GameId:   g.Id,
			GameRev:  g.Rev,
			PlayerId: m.PlayerId,
			Move:     int32(len(g.PreviousMoves)),
		}

		undoD := &flowstate.Data{}
		if err := convertor.UndoToData(u, undoD); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		undoStateCtx := &flowstate.StateCtx{
			Current: flowstate.State{
				ID: flowstate.StateID(fmt.Sprintf(`undo-%s-%d`, g.Id, g.Rev)),
				Labels: map[string]string{
					`undo.game.id`: g.Id,
				},
				Annotations: map[string]string{
					`game.id`:  g.Id,
					`game.rev`: fmt.Sprintf(`%d`, g.Rev),
				},
			},
		}

		if err := h.e.Do(flowstate.Commit(
			flowstate.StoreData(undoD),
			flowstate.ReferenceData(undoStateCtx, undoD, `undo`),
			flowstate.Pause(undoStateCtx).WithTransit(undoflow.ID),
		)); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		return connect.NewResponse(&v1.UndoResponse{
			Game: g,
			Undo: u,
		}), nil
	case req.Msg.GetDecision() != nil:
		undoDecision := req.Msg.GetDecision()

		if undoDecision.PlayerId == m.PlayerId {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("cannot decide on own undo"))
		}

		undoD := &flowstate.Data{}
		undoStateCtx := &flowstate.StateCtx{}
		if err := h.e.Do(
			flowstate.GetByID(undoStateCtx, flowstate.StateID(fmt.Sprintf(`undo-%s-%d`, g.Id, g.Rev)), 0),
			flowstate.DereferenceData(undoStateCtx, undoD, `undo`),
			flowstate.GetData(undoD),
		); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		undo, err := convertor.DataToUndo(undoD)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		if undo.Decided {
			return connect.NewResponse(&v1.UndoResponse{
				Undo: undo,
			}), nil
		}

		undo.Accepted = undoDecision.Accepted
		undo.Decided = true
		if err := convertor.UndoToData(undo, undoD); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		if !undo.Accepted {
			if err := h.e.Do(flowstate.Commit(
				flowstate.StoreData(undoD),
				flowstate.ReferenceData(undoStateCtx, undoD, `undo`),
				flowstate.Pause(undoStateCtx).WithTransit(undoflow.ID),
			)); err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}

			return connect.NewResponse(&v1.UndoResponse{
				Undo: undo,
			}), nil

		}

		m.Undone = true
		g.CurrentMove = &v1.Move{
			PlayerId: m.PlayerId,
			Color:    m.Color,
			EndAt:    time.Now().Add(time.Duration(g.MoveDurationSec) * time.Second).Unix(),
		}
		b, err := convertor.GameToBoard(g)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		g.Board = convertor.FromClamBoard(b)

		if err := convertor.GameToData(g, d); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		if err := h.e.Do(flowstate.Commit(
			flowstate.StoreData(undoD),
			flowstate.StoreData(d),
			flowstate.ReferenceData(undoStateCtx, undoD, `undo`),
			flowstate.ReferenceData(stateCtx, d, `game`),
			flowstate.Pause(undoStateCtx).WithTransit(undoflow.ID),
			flowstate.Pause(stateCtx).WithTransit(moveflow.ID),
		)); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		return connect.NewResponse(&v1.UndoResponse{
			Game: g,
			Undo: undo,
		}), nil
	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid request"))
	}
}
