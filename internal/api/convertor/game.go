package convertor

import v1 "github.com/makasim/gogame/protogen/gogame/v1"

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

func NextPlayer(g *v1.Game) *v1.Player {
	if g.CurrentMove.PlayerId == g.Player1.Id {
		return g.Player2
	}
	return g.Player1
}
