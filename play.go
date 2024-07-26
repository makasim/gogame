package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"connectrpc.com/connect"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
	"github.com/makasim/gogame/protogen/gogame/v1/gogamev1connect"
)

func main() {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		playPlayer1()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		playPlayer2()
	}()

	wg.Wait()
}

func playPlayer1() {
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
	fmt.Printf("player1: game created {Name: %s, ID: %s}\n", g.Name, g.Id)
}

func playPlayer2() {
	gsc := gogamev1connect.NewGameServiceClient(&http.Client{}, `http://127.0.0.1:8181`)

	fvgStream, err := gsc.FindVacantGames(context.Background(), connect.NewRequest(&v1.FindVacantGamesRequest{}))
	if err != nil {
		panic(fmt.Errorf("find vacant games: %w", err))
	}
	defer fvgStream.Close()

	for fvgStream.Receive() {
		if !fvgStream.Msg().Joinable {
			continue
		}

		gID := fvgStream.Msg().Game.Id

		log.Printf("player2: joining game %s", gID)
		_, err := gsc.JoinGame(context.Background(), connect.NewRequest(&v1.JoinGameRequest{
			GameId: gID,
			Player2: &v1.Player{
				Name: "aPlayer2",
			},
		}))
		if err != nil {
			log.Printf("player2: cannot join game: %s: %s", gID, err)
			continue
		}

		log.Printf("player2: game started")
		break
	}
}
