package corsmiddleware

import (
	"net/http"

	"github.com/rs/cors"
)

func Wrap(h http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowOriginVaryRequestFunc: func(r *http.Request, origin string) (bool, []string) {
			// We allow all requests to pass through, so it is equivalent to setting allowed_origins: ['*'].
			// The method is added solely for the purpose of collecting metrics.
			return true, nil
		},
		AllowedMethods:   []string{`POST`, `GET`},
		AllowedHeaders:   []string{},
		AllowCredentials: true,
		MaxAge:           600,
	})

	return c.Handler(h)
}

func WrapPath(path string, h http.Handler) (string, http.Handler) {
	return path, Wrap(h)
}
