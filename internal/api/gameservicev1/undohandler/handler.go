package undohandler

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/api/convertor"
	"github.com/makasim/gogame/internal/movetimeoutflow"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
)

type Handler struct {
	e flowstate.Engine
}

func New(e flowstate.Engine) *Handler {
	return &Handler{
		e: e,
	}
}

func (h *Handler) Undo(_ context.Context, req *connect.Request[v1.UndoRequest]) (*connect.Response[v1.UndoResponse], error) {

}
