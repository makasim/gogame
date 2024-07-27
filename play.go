package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
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
			Id:   `player1`,
			Name: "John Doe",
		},
	}))
	if err != nil {
		panic(fmt.Errorf("player1: create game: %w", err))
	}

	g := cgr.Msg.Game
	log.Printf("player1: created a game %s", g.Id)

	sgeCtx, sgeCtxCancel := context.WithCancel(context.Background())
	defer sgeCtxCancel()

	sgeStream, err := gsc.StreamGameEvents(sgeCtx, connect.NewRequest(&v1.StreamGameEventsRequest{
		GameId: g.Id,
	}))
	if err != nil {
		panic(fmt.Errorf("player1: stream game events: %w", err))
	}

	x := int64(4)
	for sgeStream.Receive() {
		if x > 9 {
			_, err := gsc.Resign(context.Background(), connect.NewRequest(&v1.ResignRequest{
				GameId:   g.Id,
				PlayerId: "player1",
			}))
			if err != nil && strings.Contains(err.Error(), `rev mismatch`) {
				continue
			} else if err != nil {
				log.Printf("player1: cannot resign: %s", err)
			}
			log.Printf("player1: resigned")
			return
		}

		g = sgeStream.Msg().Game
		switch g.State {
		case `created`:
			continue
		case `started`, `move`:
			if g.CurrentMove.PlayerId != `player1` {
				continue
			}
			if g.State == `started` {
				log.Printf("player1: plays black; first move")
			}

			m := g.CurrentMove
			m.X = x
			m.Y = 4

			x++

			mmr, err := gsc.MakeMove(context.Background(), connect.NewRequest(&v1.MakeMoveRequest{
				GameId:  g.Id,
				GameRev: g.Rev,
				Move:    m,
			}))
			if err != nil && strings.Contains(err.Error(), `rev mismatch`) {
				continue
			} else if err != nil {
				log.Printf("player1: cannot make move: %s", err)
				continue
			}
			log.Printf("player1: move made: %d:%d", m.Y, m.X)

			g = mmr.Msg.Game
		case `ended`:
			if g.Winner.Id == `player1` {
				log.Printf("player1: won by %s", g.WonBy)
			} else {
				log.Printf("player1: lost by %s", g.WonBy)
			}

			return
		default:
			panic(fmt.Errorf("player1: unknown game state: %s", g.State))
		}
	}
}

func playPlayer2() {
	gsc := gogamev1connect.NewGameServiceClient(&http.Client{}, `http://127.0.0.1:8181`)

	svgCtx, svgCtxCancel := context.WithCancel(context.Background())
	defer svgCtxCancel()

	svgStream, err := gsc.StreamVacantGames(svgCtx, connect.NewRequest(&v1.StreamVacantGamesRequest{}))
	if err != nil {
		panic(fmt.Errorf("player2: stream vacant games: %w", err))
	}

	var g *v1.Game
	for svgStream.Receive() {
		switch g.State {
		case `created`:
			g = svgStream.Msg().Game

			jgr, err := gsc.JoinGame(context.Background(), connect.NewRequest(&v1.JoinGameRequest{
				GameId: g.Id,
				Player2: &v1.Player{
					Id:   `player2`,
					Name: "Tom Harry",
				},
			}))
			if err != nil && strings.Contains(err.Error(), `game is not joinable`) {
				continue
			} else if err != nil {
				log.Printf("player2: cannot join game: %s: %s", g.Id, err)
				continue
			}

			g = jgr.Msg.Game

			break
		case `started`:
			continue
		}
	}
	svgCtxCancel()
	//svgStream.Close()

	log.Printf("player2: joined game %s", g.Id)

	sgeCtx, sgeCtxCancel := context.WithCancel(context.Background())
	defer sgeCtxCancel()

	sgeStream, err := gsc.StreamGameEvents(sgeCtx, connect.NewRequest(&v1.StreamGameEventsRequest{
		GameId: g.Id,
	}))
	if err != nil {
		panic(fmt.Errorf("player2: stream game events: %w", err))
	}

	x := int64(4)
	for sgeStream.Receive() {
		g = sgeStream.Msg().Game
		switch g.State {
		case `created`:
			continue
		case `started`, `move`:
			if g.CurrentMove.PlayerId != `player2` {
				continue
			}
			if g.State == `started` {
				log.Printf("player1: plays black; first move")
			}

			m := g.CurrentMove
			m.X = x
			m.Y = 14

			x++

			mmr, err := gsc.MakeMove(context.Background(), connect.NewRequest(&v1.MakeMoveRequest{
				GameId:  g.Id,
				GameRev: g.Rev,
				Move:    m,
			}))
			if err != nil && strings.Contains(err.Error(), `rev mismatch`) {
				continue
			} else if err != nil {
				log.Printf("player2: cannot make move: %s", err)
				continue
			}
			g = mmr.Msg.Game
			log.Printf("player2: move made: %d:%d", m.Y, m.X)
		case `ended`:
			if g.Winner.Id == `player2` {
				log.Printf("player2: won by %s", g.WonBy)
			} else {
				log.Printf("player2: lost by %s", g.WonBy)
			}

			return
		default:
			panic(fmt.Errorf("player2: unknown game state: %s", g.State))
		}
	}
}
