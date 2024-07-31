package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/makasim/flowstate"
	"github.com/makasim/flowstate/memdriver"
	"github.com/makasim/gogame/internal/api/corsmiddleware"
	"github.com/makasim/gogame/internal/api/gamehandlerv2"
	"github.com/makasim/gogame/internal/api/roomhandlerv2"
	"github.com/makasim/gogame/internal/createdflow"
	"github.com/makasim/gogame/internal/endedflow"
	"github.com/makasim/gogame/internal/moveflow"
	"github.com/makasim/gogame/protogen/gogame/v2/gogamev2connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Config struct {
}

type App struct {
	cfg Config
}

func New(cfg Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run(ctx context.Context) error {
	log.Println("app starting")

	d := memdriver.New()
	d.SetFlow(createdflow.New())
	d.SetFlow(moveflow.New())
	d.SetFlow(endedflow.New())

	e, err := flowstate.NewEngine(d)
	if err != nil {
		return fmt.Errorf("new engine: %w", err)
	}

	corsMW := corsmiddleware.New(os.Getenv(`CORS_ENABLED`) == `true`)

	mux := http.NewServeMux()
	mux.Handle(corsMW.WrapPath(gogamev2connect.NewRoomServiceHandler(roomhandlerv2.New(
		roomhandlerv2.NewCreateGameHandler(e),
		roomhandlerv2.NewJoinGameHandler(e),
		roomhandlerv2.NewStreamGamesHandler(e),
	))))
	mux.Handle(corsMW.WrapPath(gogamev2connect.NewGameServiceHandler(gamehandlerv2.New(
		gamehandlerv2.NewStreamEventsHandler(e),
		gamehandlerv2.NewMoveHandler(e),
		gamehandlerv2.NewPassHandler(e),
		gamehandlerv2.NewUndoHandler(e),
		gamehandlerv2.NewResignHandler(e),
	))))
	mux.Handle("/", corsMW.Wrap(http.FileServer(http.Dir("ui/public"))))

	srv := &http.Server{
		Addr:    `127.0.0.1:8181`,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("WARN: http server: listen and serve: %s", err)
		}
	}()

	log.Println("app started")
	<-ctx.Done()
	log.Println("app stopping")
	defer log.Println("app stopped")

	var shutdownRes error
	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), time.Second*30)
	defer shutdownCtxCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		shutdownRes = errors.Join(shutdownRes, fmt.Errorf("http server: shutdown: %w", err))
	}

	if err := e.Shutdown(shutdownCtx); err != nil {
		shutdownRes = errors.Join(shutdownRes, fmt.Errorf("engine: shutdown: %w", err))
	}

	return shutdownRes
}
