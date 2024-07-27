package corsmiddleware

import (
	"net/http"

	"github.com/rs/cors"
)

type MW struct {
	enabled bool
}

func New(enabled bool) *MW {
	return &MW{
		enabled: enabled,
	}
}

func (mv *MW) Wrap(h http.Handler) http.Handler {
	if !mv.enabled {
		return h
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{`*`},
		AllowedMethods:   []string{`POST`, `GET`},
		AllowedHeaders:   []string{`*`},
		AllowCredentials: true,
		MaxAge:           600,
	})

	return c.Handler(h)
}

func (mv *MW) WrapPath(path string, h http.Handler) (string, http.Handler) {
	return path, mv.Wrap(h)
}
