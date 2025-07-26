package joingameflow

import (
	"fmt"
	"math/rand"
	"time"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/movetimeoutflow"
	"github.com/makasim/gogame/internal/promutil"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
	"github.com/otrego/clamshell/go/board"
)

var ID flowstate.FlowID = `gogame.join_game`

type Flow struct {
}

func New() (flowstate.FlowID, *Flow) {
	return ID, &Flow{}
}

func (f *Flow) Execute(reqStateCtx *flowstate.StateCtx, e flowstate.Engine) (flowstate.Command, error) {
	msg := &v1.JoinGameRequest{}
	if err := promutil.UnmarshalRequest(reqStateCtx, msg); err != nil {
		return nil, err
	}
	if msg.GameId == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game id is required"))
	}
	if msg.Player2.Name == `` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("player2 name is required"))
	}

	g, stateCtx, d, err := convertor.FindGame(e, msg.GameId, 0)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if stateCtx.Current.Labels[`game.state`] != `created` {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game is not joinable"))
	}
	if g.Player1.Id == msg.Player2.Id {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("player1 and player2 are the same"))
	}

	stateCtx.Current.SetLabel(`game.state`, `started`)

	g.Player2 = msg.Player2
	g.State = v1.State_STATE_STARTED
	chooseFirstMove(g)

	g.Board = convertor.FromClamBoard(board.New(19))

	if err = convertor.GameToData(g, d); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := e.Do(flowstate.Commit(
		flowstate.StoreData(stateCtx, `game`),
		flowstate.Park(stateCtx),
		flowstate.Delay(stateCtx, movetimeoutflow.ID, time.Duration(g.MoveDurationSec)*time.Second),
	)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	g.Rev = int32(stateCtx.Current.Rev)

	return flowstate.Noop(), promutil.MarshalResponse(reqStateCtx, &v1.JoinGameResponse{
		Game: g,
	})
}

func chooseFirstMove(g *v1.Game) {
	rand.Seed(time.Now().UnixNano())
	players := []*v1.Player{g.Player1, g.Player2}

	i := rand.Intn(len(players))

	firstPlayer := players[i]

	g.CurrentMove = &v1.Move{
		PlayerId: firstPlayer.Id,
		Color:    v1.Color_COLOR_BLACK,
		EndAt:    time.Now().Add(time.Duration(g.MoveDurationSec) * time.Second).Unix(),
	}
}
