package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"connectrpc.com/connect"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
	"github.com/makasim/gogame/protogen/gogame/v1/gogamev1connect"
)

var winByStrategy string

func main() {
	flag.StringVar(&winByStrategy, "winby", "pass", "Choose a win by strategy: pass, resign, or timeout.")
	flag.Parse()

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
	fmt.Printf("player1: created a game %s\n", g.Id)

	sgeCtx, sgeCtxCancel := context.WithCancel(context.Background())
	defer sgeCtxCancel()

	sgeStream, err := gsc.StreamGameEvents(sgeCtx, connect.NewRequest(&v1.StreamGameEventsRequest{
		GameId: g.Id,
	}))
	if err != nil {
		panic(fmt.Errorf("player1: stream game events: %w", err))
	}

	x := int32(4)
	for sgeStream.Receive() {
		g = sgeStream.Msg().Game

		if u := sgeStream.Msg().Undo; u != nil {
			if u.PlayerId == `player1` {
				continue
			}
			if u.Decided {
				continue
			}

			// accept
			if _, err := gsc.Undo(context.Background(), connect.NewRequest(&v1.UndoRequest{
				GameId:  u.GameId,
				GameRev: u.GameRev,
				Action: &v1.UndoRequest_Decision_{
					Decision: &v1.UndoRequest_Decision{
						PlayerId: `player1`,
						Accepted: true,
					},
				},
			})); err != nil {
				fmt.Printf("player2: cannot accept undo: %s\n", err)
				continue
			}
			fmt.Printf("player1: accepted undo\n")
			continue
		}

		switch g.State {
		case v1.State_STATE_CREATED:
			continue
		case v1.State_STATE_STARTED, v1.State_STATE_MOVE:
			if g.CurrentMove.PlayerId != `player1` {
				continue
			}
			if g.State == v1.State_STATE_STARTED {
				fmt.Printf("player1: plays black; first move\n")
			}

			if x > 9 {
				switch winByStrategy {
				case "pass":
					_, err := gsc.Pass(context.Background(), connect.NewRequest(&v1.PassRequest{
						GameId:   g.Id,
						GameRev:  g.Rev,
						PlayerId: "player1",
					}))
					if err != nil && strings.Contains(err.Error(), `rev mismatch`) {
						continue
					} else if err != nil {
						fmt.Printf("player1: cannot pass: %s\n", err)
					}
					fmt.Printf("player1: passed\n")
					continue
				case "resign":
					_, err := gsc.Resign(context.Background(), connect.NewRequest(&v1.ResignRequest{
						GameId:   g.Id,
						PlayerId: "player1",
					}))
					if err != nil && strings.Contains(err.Error(), `rev mismatch`) {
						continue
					} else if err != nil {
						fmt.Printf("player1: cannot resign: %s\n", err)
					}
					fmt.Printf("player1: resigned\n")
					continue
				case "timeout":
					fmt.Printf("player1: do nothing waiting for timeout\n")
					continue
				default:
					panic(fmt.Errorf("player1: unknown win by strategy: %s", winByStrategy))
				}
			}

			m := g.CurrentMove
			m.X = x
			m.Y = 4

			x++

			go func() {
				time.Sleep(time.Millisecond * 100)
				mmr, err := gsc.MakeMove(context.Background(), connect.NewRequest(&v1.MakeMoveRequest{
					GameId:  g.Id,
					GameRev: g.Rev,
					Move:    m,
				}))
				if err != nil && strings.Contains(err.Error(), `rev mismatch`) {
					return
				} else if err != nil {
					fmt.Printf("player1: cannot make move: %s\n", err)
					return
				}
				fmt.Printf("player1: move made: %d:%d\n", m.Y, m.X)
				g = mmr.Msg.Game

				if x == 6 {
					_, err := gsc.Undo(context.Background(), connect.NewRequest(&v1.UndoRequest{
						GameId:  g.Id,
						GameRev: g.Rev,
						Action: &v1.UndoRequest_Request_{
							Request: &v1.UndoRequest_Request{
								PlayerId: "player1",
							},
						},
					}))
					if err != nil && strings.Contains(err.Error(), `rev mismatch`) {
						return
					} else if err != nil {
						fmt.Printf("player1: cannot undo: %s\n", err)
						return
					}
					fmt.Printf("player1: requested undo\n")
				}
			}()
		case v1.State_STATE_ENDED:
			if g.Winner.Id == `player1` {
				fmt.Printf("player1: won by %s\n", g.WonBy)
			} else {
				fmt.Printf("player1: lost by %s\n", g.WonBy)
			}

			return
		default:
			panic(fmt.Errorf("player1: unknown game state: %s", g.State))
		}
	}
	if err := sgeStream.Err(); err != nil {
		fmt.Printf("player2: stream game events: %s\n", err)
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
vacantGamesLoop:
	for svgStream.Receive() {
		g = svgStream.Msg().Game

		switch g.State {
		case v1.State_STATE_CREATED:
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
				fmt.Printf("player2: cannot join game: %s: %s\n", g.Id, err)
				continue
			}

			g = jgr.Msg.Game

			break vacantGamesLoop
		case v1.State_STATE_STARTED:
			continue
		}
	}
	svgCtxCancel()

	fmt.Printf("player2: joined game %s\n", g.Id)

	sgeCtx, sgeCtxCancel := context.WithCancel(context.Background())
	defer sgeCtxCancel()

	sgeStream, err := gsc.StreamGameEvents(sgeCtx, connect.NewRequest(&v1.StreamGameEventsRequest{
		GameId: g.Id,
	}))
	if err != nil {
		panic(fmt.Errorf("player2: stream game events: %w", err))
	}

	x := int32(4)
	for sgeStream.Receive() {
		g = sgeStream.Msg().Game

		if u := sgeStream.Msg().Undo; u != nil {
			if u.PlayerId == `player2` {
				continue
			}
			if u.Decided {
				continue
			}

			// accept
			if _, err := gsc.Undo(context.Background(), connect.NewRequest(&v1.UndoRequest{
				GameId:  u.GameId,
				GameRev: u.GameRev,
				Action: &v1.UndoRequest_Decision_{
					Decision: &v1.UndoRequest_Decision{
						PlayerId: `player2`,
						Accepted: true,
					},
				},
			})); err != nil {
				fmt.Printf("player2: cannot accept undo: %s\n", err)
				continue
			}
			fmt.Printf("player2: accepted undo\n")
			continue
		}

		switch g.State {
		case v1.State_STATE_STARTED, v1.State_STATE_MOVE:
			if g.CurrentMove.PlayerId != `player2` {
				continue
			}
			if g.State == v1.State_STATE_STARTED {
				fmt.Printf("player2: plays black; first move\n")
			}

			if len(g.PreviousMoves) > 0 && g.PreviousMoves[len(g.PreviousMoves)-1].Pass {
				_, err := gsc.Pass(context.Background(), connect.NewRequest(&v1.PassRequest{
					GameId:   g.Id,
					GameRev:  g.Rev,
					PlayerId: "player2",
				}))
				if err != nil && strings.Contains(err.Error(), `rev mismatch`) {
					continue
				} else if err != nil {
					fmt.Printf("player2: cannot pass: %s", err)
				}
				fmt.Printf("player2: passed\n")
				continue
			}

			m := g.CurrentMove
			m.X = x
			m.Y = 14

			x++

			go func() {
				time.Sleep(time.Millisecond * 100)
				_, err := gsc.MakeMove(context.Background(), connect.NewRequest(&v1.MakeMoveRequest{
					GameId:  g.Id,
					GameRev: g.Rev,
					Move:    m,
				}))
				if err != nil && strings.Contains(err.Error(), `rev mismatch`) {
					return
				} else if err != nil {
					fmt.Printf("player2: cannot make move: %s\n", err)
					return
				}
				fmt.Printf("player2: move made: %d:%d\n", m.Y, m.X)
				// g = mmr.Msg.Game
			}()
		case v1.State_STATE_CREATED:
			continue

		case v1.State_STATE_ENDED:
			if g.Winner.Id == `player2` {
				fmt.Printf("player2: won by %s\n", g.WonBy)
			} else {
				fmt.Printf("player2: lost by %s\n", g.WonBy)
			}

			return
		default:
			panic(fmt.Errorf("player2: unknown game state: %s", g.State))
		}
	}
	if err := sgeStream.Err(); err != nil {
		fmt.Printf("player2: stream game events: %s\n", err)
	}
}
