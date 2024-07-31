package roomhandlerv2

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/createdflow"
	"github.com/makasim/gogame/internal/moveflow"
	v2 "github.com/makasim/gogame/protogen/gogame/v2"
	"github.com/otrego/clamshell/go/board"
)

type JoinGameHandler struct {
	e *flowstate.Engine
}

func NewJoinGameHandler(e *flowstate.Engine) *JoinGameHandler {
	return &JoinGameHandler{
		e: e,
	}
}

func (h *JoinGameHandler) JoinGame(_ context.Context, req *connect.Request[v2.JoinGameRequest]) (*connect.Response[v2.JoinGameResponse], error) {
	if req.Msg.GameId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game id is required"))
	}
	if req.Msg.Player2.Name == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("player2 name is required"))
	}

	g, stateCtx, d, err := convertor.FindGame(h.e, req.Msg.GameId, 0)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if stateCtx.Current.Transition.ToID != createdflow.ID {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game is not joinable"))
	}
	if g.Player1.Id == req.Msg.Player2.Id {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("player1 and player2 are the same"))
	}

	stateCtx.Current.SetLabel(`game.state`, `started`)

	g.Player2 = req.Msg.Player2
	//g.State = v2.State_STATE_STARTED
	chooseFirstMove(g)

	g.Board = convertor.FromClamBoard(board.New(19))

	if err = convertor.GameToData(g, d); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := h.e.Do(flowstate.Commit(
		flowstate.StoreData(d),
		flowstate.ReferenceData(stateCtx, d, `game`),
		flowstate.Pause(stateCtx).WithTransit(moveflow.ID),
		flowstate.Delay(stateCtx, time.Second*30).WithCommit(true),
	)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	g.Rev = int32(stateCtx.Current.Rev)

	return connect.NewResponse(&v2.JoinGameResponse{
		Game: g,
	}), nil
}

func chooseFirstMove(g *v2.Game) {
	rand.Seed(time.Now().UnixNano())
	players := []*v2.Player{g.Player1, g.Player2}

	i := rand.Intn(len(players))

	firstPlayer := players[i]

	g.Changes = append(g.Changes, &v2.Change{
		Change: &v2.Change_Move_{
			Move: &v2.Change_Move{
				PlayerId: firstPlayer.Id,
				Color:    v2.Color_COLOR_BLACK,
			},
		},
	})
}
