package convertor

import (
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
	"github.com/otrego/clamshell/go/color"
	"github.com/otrego/clamshell/go/move"
	"github.com/otrego/clamshell/go/point"
)

func AnotherPlayer(g *v1.Game, pID string) *v1.Player {
	if g.Player1.Id == pID {
		return g.Player2
	}
	return g.Player1
}

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
