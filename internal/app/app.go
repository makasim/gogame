package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/makasim/flowstate"
	"github.com/makasim/flowstate/netdriver"
	"github.com/makasim/flowstate/netflow"
	"github.com/makasim/gogame/internal/api/corsmiddleware"
	"github.com/makasim/gogame/internal/api/gameservicev1"
	"github.com/makasim/gogame/internal/api/gameservicev1/creategamehandler"
	"github.com/makasim/gogame/internal/api/gameservicev1/joingamehandler"
	"github.com/makasim/gogame/internal/api/gameservicev1/makemovehandler"
	"github.com/makasim/gogame/internal/api/gameservicev1/passhandler"
	"github.com/makasim/gogame/internal/api/gameservicev1/resignhandler"
	"github.com/makasim/gogame/internal/api/gameservicev1/streamgameeventshandler"
	"github.com/makasim/gogame/internal/api/gameservicev1/streamvacantgameshandler"
	"github.com/makasim/gogame/internal/api/gameservicev1/undohandler"
	"github.com/makasim/gogame/internal/createdflow"
	"github.com/makasim/gogame/internal/endedflow"
	"github.com/makasim/gogame/internal/moveflow"
	"github.com/makasim/gogame/protogen/gogame/v1/gogamev1connect"
	"github.com/makasim/gogame/ui"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Config struct {
}

type App struct {
	cfg Config
	l   *slog.Logger
}

func New(cfg Config) *App {
	return &App{
		cfg: cfg,
		l:   slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (a *App) Run(ctx context.Context) error {
	log.Println("app starting")

	flowstateHttpHost := os.Getenv("FLOWSTATE_HTTP_HOST")
	if flowstateHttpHost == "" {
		flowstateHttpHost = "http://localhost:8080"
	}

	a.l.Info("connecting to flowstate at: " + flowstateHttpHost)

	d := netdriver.New(flowstateHttpHost)

	httpHost := os.Getenv("HTTP_HOST")
	if httpHost == "" {
		httpHost = "http://localhost:8181"
	}

	a.l.Info("flow execute server at: " + httpHost)

	fr := netflow.NewRegistry(httpHost, d, a.l)
	defer fr.Close()

	if err := fr.SetFlow(createdflow.New()); err != nil {
		return fmt.Errorf("set flow created: %w", err)
	}
	if err := fr.SetFlow(moveflow.New()); err != nil {
		return fmt.Errorf("set flow move: %w", err)
	}
	if err := fr.SetFlow(endedflow.New()); err != nil {
		return fmt.Errorf("set flow ended: %w", err)
	}

	e, err := flowstate.NewEngine(d, fr, a.l)
	if err != nil {
		return fmt.Errorf("new engine: %w", err)
	}

	corsEnv := os.Getenv(`CORS_ENABLED`)
	corsMW := corsmiddleware.New(corsEnv == `true` || corsEnv == ``)

	mux := http.NewServeMux()
	mux.Handle(corsMW.WrapPath(gogamev1connect.NewGameServiceHandler(gameservicev1.New(
		creategamehandler.New(e),
		joingamehandler.New(e),
		streamvacantgameshandler.New(e),
		streamgameeventshandler.New(e),
		makemovehandler.New(e),
		resignhandler.New(e),
		passhandler.New(e),
		undohandler.New(e),
	))))

	mux.Handle("/", corsMW.Wrap(http.FileServerFS(ui.PublicFS())))

	srv := &http.Server{
		Addr: `0.0.0.0:8181`,
		Handler: h2c.NewHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if netflow.HandleExecute(rw, r, e) {
				return
			}

			mux.ServeHTTP(rw, r)
		}), &http2.Server{}),
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
