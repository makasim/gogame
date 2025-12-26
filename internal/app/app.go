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
	"github.com/makasim/gogame/internal/api"
	"github.com/makasim/gogame/internal/api/gameservicev1"
	"github.com/makasim/gogame/internal/api/gameservicev1/streamgameeventshandler"
	"github.com/makasim/gogame/internal/api/gameservicev1/streamvacantgameshandler"
	"github.com/makasim/gogame/internal/creategameflow"
	"github.com/makasim/gogame/internal/joingameflow"
	"github.com/makasim/gogame/internal/makemoveflow"
	"github.com/makasim/gogame/internal/movetimeoutflow"
	"github.com/makasim/gogame/internal/passflow"
	"github.com/makasim/gogame/internal/resignflow"
	"github.com/makasim/gogame/internal/staleflow"
	"github.com/makasim/gogame/internal/undoflow"
	"github.com/makasim/gogame/protogen/gogame/v1/gogamev1connect"
	"github.com/makasim/gogame/ui"
	"github.com/rs/cors"
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

	//d := memdriver.New(a.l)
	d := netdriver.New(flowstateHttpHost)

	httpHost := os.Getenv("HTTP_HOST")
	if httpHost == "" {
		httpHost = "http://localhost:8181"
	}

	a.l.Info("flow execute server at: " + httpHost)

	fr := netflow.NewRegistry(httpHost, d, a.l)
	defer fr.Close()

	if err := fr.SetFlow(creategameflow.New()); err != nil {
		return fmt.Errorf("flow registry: set flow: creategameflow: %w", err)
	}
	if err := fr.SetFlow(joingameflow.New()); err != nil {
		return fmt.Errorf("flow registry: set flow: joingameflow: %w", err)
	}
	if err := fr.SetFlow(makemoveflow.New()); err != nil {
		return fmt.Errorf("flow registry: set flow: makemoveflow: %w", err)
	}
	if err := fr.SetFlow(resignflow.New()); err != nil {
		return fmt.Errorf("flow registry: set flow: resignflow: %w", err)
	}
	if err := fr.SetFlow(passflow.New()); err != nil {
		return fmt.Errorf("flow registry: set flow: passflow: %w", err)
	}
	if err := fr.SetFlow(undoflow.New()); err != nil {
		return fmt.Errorf("flow registry: set flow: undoflow: %w", err)
	}
	if err := fr.SetFlow(movetimeoutflow.New()); err != nil {
		return fmt.Errorf("set flow move: %w", err)
	}
	if err := fr.SetFlow(staleflow.New()); err != nil {
		return fmt.Errorf("set flow stale: %w", err)
	}

	e, err := flowstate.NewEngine(d, fr, a.l)
	if err != nil {
		return fmt.Errorf("new engine: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle(gogamev1connect.NewGameServiceHandler(gameservicev1.New(
		streamvacantgameshandler.New(e),
		streamgameeventshandler.New(e),
	)))

	mux.Handle("/", http.FileServerFS(ui.PublicFS()))

	srv := &http.Server{
		Addr: `0.0.0.0:8181`,
		Handler: h2c.NewHandler(handleCORS(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if netflow.HandleExecute(rw, r, e) {
				return
			}
			if api.HandleAll(rw, r, e, a.l) {
				return
			}

			mux.ServeHTTP(rw, r)
		})), &http2.Server{}),
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

func handleCORS(h http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{`*`},
		AllowedMethods:   []string{`POST`, `GET`},
		AllowedHeaders:   []string{`*`},
		AllowCredentials: true,
		MaxAge:           600,
	}).Handler(h)
}
