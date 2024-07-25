package main

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
	"github.com/makasim/gogame/protogen/gogame/v1/gogamev1connect"
)

func main() {
	gsc := gogamev1connect.NewGameServiceClient(&http.Client{}, `http://127.0.0.1:8181`)

	cgr, err := gsc.CreateGame(context.Background(), connect.NewRequest(&v1.CreateGameRequest{
		Name: "aGame",
		Player1: &v1.Player{
			Name: "aPlayer1",
		},
	}))
	if err != nil {
		panic(fmt.Errorf("create game: %w", err))
	}
	g := cgr.Msg.Game
	fmt.Printf("Game created {Name: %s, ID: %s}\n", g.Name, g.Id)

	jgr, err := gsc.JoinGame(context.Background(), connect.NewRequest(&v1.JoinGameRequest{
		GameId: g.Id,
		Player2: &v1.Player{
			Name: "aPlayer2",
		},
	}))
	if err != nil {
		panic(fmt.Errorf("join game: %w", err))
	}
	g = jgr.Msg.Game

	var movePlayer *v1.Player
	if g.Player1.Color == v1.Color_COLOR_BLACK {
		movePlayer = g.Player1
	} else {
		movePlayer = g.Player2
	}

	fmt.Printf("Game started. Player %s plays black\n", movePlayer.Name)
}
