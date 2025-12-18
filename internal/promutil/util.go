package promutil

import (
	"encoding/json"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/VictoriaMetrics/easyproto"
	"github.com/makasim/flowstate"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func WriteOK(rw http.ResponseWriter, d *flowstate.Data) {
	rw.Header().Set("Content-Type", d.Annotations["content-type"])
	rw.WriteHeader(http.StatusOK)

	_, _ = rw.Write(d.Blob)
}

func WriteConnectError(rw http.ResponseWriter, err *connect.Error, proto bool) {
	switch err.Code() {
	case connect.CodeInvalidArgument:
		WriteInvalidArgumentError(rw, err.Message(), proto)
	case connect.CodeUnknown:
		WriteUnknownError(rw, err.Message(), proto)
	//case connect.CodeInternal:
	//	WriteInternalError(rw, connErr.Message(), proto)
	//case connect.CodeUnimplemented:
	//	WriteUnimplementedError(rw, connErr.Message(), proto)
	case connect.CodeNotFound:
		WriteNotFoundError(rw, err.Message(), proto)
	//case connect.CodeAlreadyExists:
	//	WriteAlreadyExistsError(rw, connErr.Message(), proto)
	//case connect.CodeUnauthenticated:
	//	WriteUnauthenticatedError(rw, connErr.Message(), proto)
	//case connect.CodePermissionDenied:
	//	WritePermissionDeniedError(rw, connErr.Message(), proto)
	default:
		WriteUnknownError(rw, err.Message(), proto)
	}
}

func WriteUnknownError(rw http.ResponseWriter, message string, proto bool) {
	if proto {
		rw.Header().Set("Content-Type", "application/proto")
	} else {
		rw.Header().Set("Content-Type", "application/json")
	}

	rw.WriteHeader(http.StatusInternalServerError)

	if proto {
		_, _ = rw.Write(MarshalError("unknown", message))
	} else {
		_, _ = rw.Write(MarshalJSONError("unknown", message))
	}
}

func WriteInvalidArgumentError(rw http.ResponseWriter, message string, proto bool) {
	if proto {
		rw.Header().Set("Content-Type", "application/proto")
	} else {
		rw.Header().Set("Content-Type", "application/json")
	}

	rw.WriteHeader(http.StatusBadRequest)

	if proto {
		_, _ = rw.Write(MarshalError("invalid_argument", message))
	} else {
		_, _ = rw.Write(MarshalJSONError("invalid_argument", message))
	}
}

func WriteNotFoundError(rw http.ResponseWriter, message string, proto bool) {
	if proto {
		rw.Header().Set("Content-Type", "application/proto")
	} else {
		rw.Header().Set("Content-Type", "application/json")
	}

	rw.WriteHeader(http.StatusNotFound)

	if proto {
		_, _ = rw.Write(MarshalError("not_found", message))
	} else {
		_, _ = rw.Write(MarshalJSONError("not_found", message))
	}
}

func MarshalJSONError(code, message string) []byte {
	b, _ := json.Marshal(map[string]string{
		"code":    code,
		"message": message,
	})
	return b
}

func MarshalError(code, message string) []byte {
	m := &easyproto.Marshaler{}
	mm := m.MessageMarshaler()

	if code != "" {
		mm.AppendString(1, code)
	}
	if message != "" {
		mm.AppendString(2, message)
	}

	return m.Marshal(nil)
}

func UnmarshalRequest(stateCtx *flowstate.StateCtx, msg proto.Message) error {
	reqData, err := stateCtx.Data("request")
	if err != nil {
		return fmt.Errorf("failed to get request data: %w", err)
	}

	// Unmarshal based on content-type annotation
	contentType := reqData.Annotations["content-type"]
	switch contentType {
	case "application/proto", "application/protobuf", "application/x-protobuf":
		if err := proto.Unmarshal(reqData.Blob, msg); err != nil {
			return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("failed to unmarshal protobuf request: %w", err))
		}
	case "application/json", "application/protojson":
		if err := protojson.Unmarshal(reqData.Blob, msg); err != nil {
			return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("failed to unmarshal protojson request: %w", err))
		}
	default:
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unsupported content-type: %s", contentType))
	}

	return nil
}

func MarshalResponse(stateCtx *flowstate.StateCtx, msg proto.Message) error {
	reqData := stateCtx.MustData("request")

	respData := &flowstate.Data{
		Annotations: map[string]string{
			"content-type": reqData.Annotations["content-type"],
		},
	}

	contentType := respData.Annotations["content-type"]
	switch contentType {
	case "application/proto", "application/protobuf", "application/x-protobuf":
		b, err := proto.Marshal(msg)
		if err != nil {
			return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("failed to unmarshal protobuf request: %w", err))
		}
		respData.Blob = b
	case "application/json", "application/protojson":
		b, err := protojson.Marshal(msg)
		if err != nil {
			return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("failed to unmarshal protojson request: %w", err))
		}
		respData.Blob = b
	default:
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unsupported content-type: %s", contentType))
	}

	stateCtx.Datas["response"] = respData

	return nil
}
