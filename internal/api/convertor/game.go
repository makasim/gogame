package convertor

import (
	v2 "github.com/makasim/gogame/protogen/gogame/v2"
	"github.com/otrego/clamshell/go/board"
	"github.com/otrego/clamshell/go/color"
	"github.com/otrego/clamshell/go/move"
	"github.com/otrego/clamshell/go/point"
)

func NextColor(g *v2.Game) v2.Color {
	lm := LastMove(g)
	if lm.Color == v2.Color_COLOR_BLACK {
		return v2.Color_COLOR_WHITE
	}
	return v2.Color_COLOR_BLACK
}

func LastMove(g *v2.Game) *v2.Change_Move {
	for i := len(g.Changes) - 1; i >= 0; i-- {
		if g.Changes[i].GetMove() != nil {
			return g.Changes[i].GetMove()
		}
	}

	return nil
}

func LastPlayer(g *v2.Game) *v2.Player {
	for i := len(g.Changes) - 1; i >= 0; i-- {
		if g.Changes[i].GetMove() != nil {
			if g.Changes[i].GetMove().PlayerId == g.Player1.Id {
				return g.Player2
			}
			return g.Player1
		}
		if g.Changes[i].GetPass() != nil {
			if g.Changes[i].GetPass().PlayerId == g.Player1.Id {
				return g.Player2
			}
			return g.Player1
		}
	}

	return nil
}

func AnotherPlayer(g *v2.Game, pID string) *v2.Player {
	if g.Player1.Id == pID {
		return g.Player2
	}
	return g.Player1
}

func ToClamMove(m *v2.Change_Move) *move.Move {
	clamColor := color.Black
	if m.Color == v2.Color_COLOR_WHITE {
		clamColor = color.White
	}

	return move.New(clamColor, point.New(int(m.X), int(m.Y)))
}

func FromClamBoard(clamBoard *board.Board) *v2.Board {
	b := &v2.Board{
		Size: 19,
		Rows: make([]*v2.Row, 19),
	}
	for i := 0; i < 19; i++ {
		b.Rows[i] = &v2.Row{
			Colors: make([]v2.Color, 19),
		}
	}

	clamBoardState := clamBoard.FullBoardState()

	for i, clamRow := range clamBoardState {
		for j, clamColor := range clamRow {
			switch clamColor {
			case color.Black:
				b.Rows[i].Colors[j] = v2.Color_COLOR_BLACK
			case color.White:
				b.Rows[i].Colors[j] = v2.Color_COLOR_WHITE
			}
		}
	}

	return b
}
