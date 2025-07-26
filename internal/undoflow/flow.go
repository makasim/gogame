package undoflow

import (
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/movetimeoutflow"
	"github.com/makasim/gogame/internal/promutil"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
)

var ID flowstate.FlowID = `gogame.undo`

type Flow struct {
}

func New() (flowstate.FlowID, *Flow) {
	return ID, &Flow{}
}

func (f *Flow) Execute(reqStateCtx *flowstate.StateCtx, e flowstate.Engine) (flowstate.Command, error) {
	msg := &v1.UndoRequest{}
	if err := promutil.UnmarshalRequest(reqStateCtx, msg); err != nil {
		return nil, err
	}
	if msg.GameId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game id is required"))
	}
	if msg.GameRev <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game rev is required"))
	}

	g, stateCtx, d, err := convertor.FindGame(e, msg.GameId, msg.GameRev)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if len(g.PreviousMoves) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("no moves to undo"))
	}

	m := g.PreviousMoves[len(g.PreviousMoves)-1]

	switch {
	case msg.GetRequest() != nil:
		undoReq := msg.GetRequest()

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
			Datas: map[string]*flowstate.Data{
				"undo": undoD,
			},
		}

		if err := e.Do(flowstate.Commit(
			flowstate.StoreData(undoStateCtx, `undo`),
			flowstate.Park(undoStateCtx),
		)); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		return flowstate.Noop(), promutil.MarshalResponse(reqStateCtx, &v1.UndoResponse{
			Game: g,
			Undo: u,
		})
	case msg.GetDecision() != nil:
		undoDecision := msg.GetDecision()

		if undoDecision.PlayerId == m.PlayerId {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("cannot decide on own undo"))
		}

		undoStateCtx := &flowstate.StateCtx{}
		if err := e.Do(
			flowstate.GetStateByID(undoStateCtx, flowstate.StateID(fmt.Sprintf(`undo-%s-%d`, g.Id, g.Rev)), 0),
			flowstate.GetData(undoStateCtx, `undo`),
		); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		undoD := undoStateCtx.MustData(`undo`)
		undo, err := convertor.DataToUndo(undoD)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		if undo.Decided {
			return flowstate.Noop(), promutil.MarshalResponse(reqStateCtx, &v1.UndoResponse{
				Undo: undo,
			})
		}

		undo.Accepted = undoDecision.Accepted
		undo.Decided = true
		if err := convertor.UndoToData(undo, undoD); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		if !undo.Accepted {
			if err := e.Do(flowstate.Commit(
				flowstate.StoreData(undoStateCtx, `undo`),
				flowstate.Park(undoStateCtx),
			)); err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}

			return flowstate.Noop(), promutil.MarshalResponse(reqStateCtx, &v1.UndoResponse{
				Undo: undo,
			})
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

		if err := e.Do(flowstate.Commit(
			flowstate.StoreData(undoStateCtx, `undo`),
			flowstate.StoreData(stateCtx, `game`),
			flowstate.Park(undoStateCtx),
			flowstate.Park(stateCtx),
			flowstate.Delay(stateCtx, movetimeoutflow.ID, time.Duration(g.MoveDurationSec)*time.Second),
		)); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		return flowstate.Noop(), promutil.MarshalResponse(reqStateCtx, &v1.UndoResponse{
			Game: g,
			Undo: undo,
		})
	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid request"))
	}
}
