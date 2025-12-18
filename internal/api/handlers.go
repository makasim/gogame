package api

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/gogame/internal/creategameflow"
	"github.com/makasim/gogame/internal/joingameflow"
	"github.com/makasim/gogame/internal/makemoveflow"
	"github.com/makasim/gogame/internal/passflow"
	"github.com/makasim/gogame/internal/promutil"
	"github.com/makasim/gogame/internal/resignflow"
	"github.com/makasim/gogame/internal/undoflow"
	"github.com/oklog/ulid/v2"
)

func HandleAll(rw http.ResponseWriter, r *http.Request, e *flowstate.Engine, l *slog.Logger) bool {
	if HandleCreateGame(rw, r, e, l) {
		return true
	}
	if HandleJoinGame(rw, r, e, l) {
		return true
	}
	if HandleMakeMove(rw, r, e, l) {
		return true
	}
	if HandleResign(rw, r, e, l) {
		return true
	}
	if HandlePass(rw, r, e, l) {
		return true
	}
	if HandleUndo(rw, r, e, l) {
		return true
	}

	return false
}

func HandleCreateGame(rw http.ResponseWriter, r *http.Request, e *flowstate.Engine, l *slog.Logger) bool {
	if r.URL.Path != "/gogame.v1.GameService/CreateGame" {
		return false
	}

	res := handleFlow(rw, r, e, creategameflow.ID, l)
	return res
}

func HandleJoinGame(rw http.ResponseWriter, r *http.Request, e *flowstate.Engine, l *slog.Logger) bool {
	if r.URL.Path != "/gogame.v1.GameService/JoinGame" {
		return false
	}

	res := handleFlow(rw, r, e, joingameflow.ID, l)
	return res
}

func HandleMakeMove(rw http.ResponseWriter, r *http.Request, e *flowstate.Engine, l *slog.Logger) bool {
	if r.URL.Path != "/gogame.v1.GameService/MakeMove" {
		return false
	}

	res := handleFlow(rw, r, e, makemoveflow.ID, l)
	return res
}

func HandleResign(rw http.ResponseWriter, r *http.Request, e *flowstate.Engine, l *slog.Logger) bool {
	if r.URL.Path != "/gogame.v1.GameService/Resign" {
		return false
	}

	res := handleFlow(rw, r, e, resignflow.ID, l)
	return res
}

func HandlePass(rw http.ResponseWriter, r *http.Request, e *flowstate.Engine, l *slog.Logger) bool {
	if r.URL.Path != "/gogame.v1.GameService/Pass" {
		return false
	}

	res := handleFlow(rw, r, e, passflow.ID, l)
	return res
}

func HandleUndo(rw http.ResponseWriter, r *http.Request, e *flowstate.Engine, l *slog.Logger) bool {
	if r.URL.Path != "/gogame.v1.GameService/Undo" {
		return false
	}

	res := handleFlow(rw, r, e, undoflow.ID, l)
	return res
}

func handleFlow(rw http.ResponseWriter, r *http.Request, e *flowstate.Engine, fID flowstate.FlowID, l *slog.Logger) bool {
	proto := r.Header.Get("Content-Type") != "application/json"

	b, err := io.ReadAll(r.Body)
	if err != nil {
		promutil.WriteInvalidArgumentError(rw, "failed to read request body: "+err.Error(), proto)
		return true
	}

	stateCtx := &flowstate.StateCtx{
		Current: flowstate.State{
			ID: flowstate.StateID(ulid.Make().String()),
			Transition: flowstate.Transition{
				To: fID,
			},
		},
		Datas: map[string]*flowstate.Data{
			"request": {
				Annotations: map[string]string{
					"content-type": r.Header.Get("Content-Type"),
				},
				Blob: b,
			},
		},
	}

	if err := e.Execute(stateCtx); err != nil {
		var connErr *connect.Error
		if errors.As(err, &connErr) {
			promutil.WriteConnectError(rw, connErr, proto)
			return true
		}

		l.Error("engine execute", "flow", stateCtx.Current.Transition.To, "error", err)
		promutil.WriteUnknownError(rw, err.Error(), proto)
		return true
	}

	respData, err := stateCtx.Data("response")
	if err != nil {
		promutil.WriteUnknownError(rw, err.Error(), proto)
		return true
	}

	promutil.WriteOK(rw, respData)
	return true
}
