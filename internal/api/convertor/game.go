package convertor

import (
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
	"github.com/otrego/clamshell/go/board"
	"github.com/otrego/clamshell/go/color"
	"github.com/otrego/clamshell/go/move"
	"github.com/otrego/clamshell/go/point"
)

func NextColor(g *v1.Game) v1.Color {
	if g.CurrentMove.Color == v1.Color_COLOR_BLACK {
		return v1.Color_COLOR_WHITE
	}
	return v1.Color_COLOR_BLACK
}

func CurrentPlayer(g *v1.Game) *v1.Player {
	if g.CurrentMove.PlayerId == g.Player1.Id {
		return g.Player1
	}
	return g.Player2
}

func NextPlayer(g *v1.Game) *v1.Player {
	if g.CurrentMove.PlayerId == g.Player1.Id {
		return g.Player2
	}
	return g.Player1
}

func ToClamMove(m *v1.Move) *move.Move {
	clamColor := color.Black
	if m.Color == v1.Color_COLOR_WHITE {
		clamColor = color.White
	}

	return move.New(clamColor, point.New(int(m.X), int(m.Y)))
}

func FromClamBoard(clamBoard *board.Board) *v1.Board {
	b := &v1.Board{
		Size: 19,
		Rows: make([]*v1.Row, 19),
	}
	for i := 0; i < 19; i++ {
		b.Rows[i] = &v1.Row{
			Colors: make([]v1.Color, 19),
		}
	}

	clamBoardState := clamBoard.FullBoardState()

	for i, clamRow := range clamBoardState {
		for j, clamColor := range clamRow {
			switch clamColor {
			case color.Black:
				b.Rows[i].Colors[j] = v1.Color_COLOR_BLACK
			case color.White:
				b.Rows[i].Colors[j] = v1.Color_COLOR_WHITE
			}
		}
	}

	return b
}
