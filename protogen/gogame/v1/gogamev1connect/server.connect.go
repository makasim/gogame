// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: gogame/v1/server.proto

package gogamev1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// GameServiceName is the fully-qualified name of the GameService service.
	GameServiceName = "gogame.v1.GameService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// GameServiceCreateGameProcedure is the fully-qualified name of the GameService's CreateGame RPC.
	GameServiceCreateGameProcedure = "/gogame.v1.GameService/CreateGame"
	// GameServiceStreamVacantGamesProcedure is the fully-qualified name of the GameService's
	// StreamVacantGames RPC.
	GameServiceStreamVacantGamesProcedure = "/gogame.v1.GameService/StreamVacantGames"
	// GameServiceJoinGameProcedure is the fully-qualified name of the GameService's JoinGame RPC.
	GameServiceJoinGameProcedure = "/gogame.v1.GameService/JoinGame"
	// GameServiceStreamGameEventsProcedure is the fully-qualified name of the GameService's
	// StreamGameEvents RPC.
	GameServiceStreamGameEventsProcedure = "/gogame.v1.GameService/StreamGameEvents"
	// GameServiceMakeMoveProcedure is the fully-qualified name of the GameService's MakeMove RPC.
	GameServiceMakeMoveProcedure = "/gogame.v1.GameService/MakeMove"
	// GameServiceResignProcedure is the fully-qualified name of the GameService's Resign RPC.
	GameServiceResignProcedure = "/gogame.v1.GameService/Resign"
	// GameServicePassProcedure is the fully-qualified name of the GameService's Pass RPC.
	GameServicePassProcedure = "/gogame.v1.GameService/Pass"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	gameServiceServiceDescriptor                 = v1.File_gogame_v1_server_proto.Services().ByName("GameService")
	gameServiceCreateGameMethodDescriptor        = gameServiceServiceDescriptor.Methods().ByName("CreateGame")
	gameServiceStreamVacantGamesMethodDescriptor = gameServiceServiceDescriptor.Methods().ByName("StreamVacantGames")
	gameServiceJoinGameMethodDescriptor          = gameServiceServiceDescriptor.Methods().ByName("JoinGame")
	gameServiceStreamGameEventsMethodDescriptor  = gameServiceServiceDescriptor.Methods().ByName("StreamGameEvents")
	gameServiceMakeMoveMethodDescriptor          = gameServiceServiceDescriptor.Methods().ByName("MakeMove")
	gameServiceResignMethodDescriptor            = gameServiceServiceDescriptor.Methods().ByName("Resign")
	gameServicePassMethodDescriptor              = gameServiceServiceDescriptor.Methods().ByName("Pass")
)

// GameServiceClient is a client for the gogame.v1.GameService service.
type GameServiceClient interface {
	CreateGame(context.Context, *connect.Request[v1.CreateGameRequest]) (*connect.Response[v1.CreateGameResponse], error)
	StreamVacantGames(context.Context, *connect.Request[v1.StreamVacantGamesRequest]) (*connect.ServerStreamForClient[v1.StreamVacantGamesResponse], error)
	JoinGame(context.Context, *connect.Request[v1.JoinGameRequest]) (*connect.Response[v1.JoinGameResponse], error)
	StreamGameEvents(context.Context, *connect.Request[v1.StreamGameEventsRequest]) (*connect.ServerStreamForClient[v1.StreamGameEventsResponse], error)
	MakeMove(context.Context, *connect.Request[v1.MakeMoveRequest]) (*connect.Response[v1.MakeMoveResponse], error)
	Resign(context.Context, *connect.Request[v1.ResignRequest]) (*connect.Response[v1.ResignResponse], error)
	Pass(context.Context, *connect.Request[v1.PassRequest]) (*connect.Response[v1.PassResponse], error)
}

// NewGameServiceClient constructs a client for the gogame.v1.GameService service. By default, it
// uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewGameServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) GameServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &gameServiceClient{
		createGame: connect.NewClient[v1.CreateGameRequest, v1.CreateGameResponse](
			httpClient,
			baseURL+GameServiceCreateGameProcedure,
			connect.WithSchema(gameServiceCreateGameMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		streamVacantGames: connect.NewClient[v1.StreamVacantGamesRequest, v1.StreamVacantGamesResponse](
			httpClient,
			baseURL+GameServiceStreamVacantGamesProcedure,
			connect.WithSchema(gameServiceStreamVacantGamesMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		joinGame: connect.NewClient[v1.JoinGameRequest, v1.JoinGameResponse](
			httpClient,
			baseURL+GameServiceJoinGameProcedure,
			connect.WithSchema(gameServiceJoinGameMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		streamGameEvents: connect.NewClient[v1.StreamGameEventsRequest, v1.StreamGameEventsResponse](
			httpClient,
			baseURL+GameServiceStreamGameEventsProcedure,
			connect.WithSchema(gameServiceStreamGameEventsMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		makeMove: connect.NewClient[v1.MakeMoveRequest, v1.MakeMoveResponse](
			httpClient,
			baseURL+GameServiceMakeMoveProcedure,
			connect.WithSchema(gameServiceMakeMoveMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		resign: connect.NewClient[v1.ResignRequest, v1.ResignResponse](
			httpClient,
			baseURL+GameServiceResignProcedure,
			connect.WithSchema(gameServiceResignMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		pass: connect.NewClient[v1.PassRequest, v1.PassResponse](
			httpClient,
			baseURL+GameServicePassProcedure,
			connect.WithSchema(gameServicePassMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// gameServiceClient implements GameServiceClient.
type gameServiceClient struct {
	createGame        *connect.Client[v1.CreateGameRequest, v1.CreateGameResponse]
	streamVacantGames *connect.Client[v1.StreamVacantGamesRequest, v1.StreamVacantGamesResponse]
	joinGame          *connect.Client[v1.JoinGameRequest, v1.JoinGameResponse]
	streamGameEvents  *connect.Client[v1.StreamGameEventsRequest, v1.StreamGameEventsResponse]
	makeMove          *connect.Client[v1.MakeMoveRequest, v1.MakeMoveResponse]
	resign            *connect.Client[v1.ResignRequest, v1.ResignResponse]
	pass              *connect.Client[v1.PassRequest, v1.PassResponse]
}

// CreateGame calls gogame.v1.GameService.CreateGame.
func (c *gameServiceClient) CreateGame(ctx context.Context, req *connect.Request[v1.CreateGameRequest]) (*connect.Response[v1.CreateGameResponse], error) {
	return c.createGame.CallUnary(ctx, req)
}

// StreamVacantGames calls gogame.v1.GameService.StreamVacantGames.
func (c *gameServiceClient) StreamVacantGames(ctx context.Context, req *connect.Request[v1.StreamVacantGamesRequest]) (*connect.ServerStreamForClient[v1.StreamVacantGamesResponse], error) {
	return c.streamVacantGames.CallServerStream(ctx, req)
}

// JoinGame calls gogame.v1.GameService.JoinGame.
func (c *gameServiceClient) JoinGame(ctx context.Context, req *connect.Request[v1.JoinGameRequest]) (*connect.Response[v1.JoinGameResponse], error) {
	return c.joinGame.CallUnary(ctx, req)
}

// StreamGameEvents calls gogame.v1.GameService.StreamGameEvents.
func (c *gameServiceClient) StreamGameEvents(ctx context.Context, req *connect.Request[v1.StreamGameEventsRequest]) (*connect.ServerStreamForClient[v1.StreamGameEventsResponse], error) {
	return c.streamGameEvents.CallServerStream(ctx, req)
}

// MakeMove calls gogame.v1.GameService.MakeMove.
func (c *gameServiceClient) MakeMove(ctx context.Context, req *connect.Request[v1.MakeMoveRequest]) (*connect.Response[v1.MakeMoveResponse], error) {
	return c.makeMove.CallUnary(ctx, req)
}

// Resign calls gogame.v1.GameService.Resign.
func (c *gameServiceClient) Resign(ctx context.Context, req *connect.Request[v1.ResignRequest]) (*connect.Response[v1.ResignResponse], error) {
	return c.resign.CallUnary(ctx, req)
}

// Pass calls gogame.v1.GameService.Pass.
func (c *gameServiceClient) Pass(ctx context.Context, req *connect.Request[v1.PassRequest]) (*connect.Response[v1.PassResponse], error) {
	return c.pass.CallUnary(ctx, req)
}

// GameServiceHandler is an implementation of the gogame.v1.GameService service.
type GameServiceHandler interface {
	CreateGame(context.Context, *connect.Request[v1.CreateGameRequest]) (*connect.Response[v1.CreateGameResponse], error)
	StreamVacantGames(context.Context, *connect.Request[v1.StreamVacantGamesRequest], *connect.ServerStream[v1.StreamVacantGamesResponse]) error
	JoinGame(context.Context, *connect.Request[v1.JoinGameRequest]) (*connect.Response[v1.JoinGameResponse], error)
	StreamGameEvents(context.Context, *connect.Request[v1.StreamGameEventsRequest], *connect.ServerStream[v1.StreamGameEventsResponse]) error
	MakeMove(context.Context, *connect.Request[v1.MakeMoveRequest]) (*connect.Response[v1.MakeMoveResponse], error)
	Resign(context.Context, *connect.Request[v1.ResignRequest]) (*connect.Response[v1.ResignResponse], error)
	Pass(context.Context, *connect.Request[v1.PassRequest]) (*connect.Response[v1.PassResponse], error)
}

// NewGameServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewGameServiceHandler(svc GameServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	gameServiceCreateGameHandler := connect.NewUnaryHandler(
		GameServiceCreateGameProcedure,
		svc.CreateGame,
		connect.WithSchema(gameServiceCreateGameMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	gameServiceStreamVacantGamesHandler := connect.NewServerStreamHandler(
		GameServiceStreamVacantGamesProcedure,
		svc.StreamVacantGames,
		connect.WithSchema(gameServiceStreamVacantGamesMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	gameServiceJoinGameHandler := connect.NewUnaryHandler(
		GameServiceJoinGameProcedure,
		svc.JoinGame,
		connect.WithSchema(gameServiceJoinGameMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	gameServiceStreamGameEventsHandler := connect.NewServerStreamHandler(
		GameServiceStreamGameEventsProcedure,
		svc.StreamGameEvents,
		connect.WithSchema(gameServiceStreamGameEventsMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	gameServiceMakeMoveHandler := connect.NewUnaryHandler(
		GameServiceMakeMoveProcedure,
		svc.MakeMove,
		connect.WithSchema(gameServiceMakeMoveMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	gameServiceResignHandler := connect.NewUnaryHandler(
		GameServiceResignProcedure,
		svc.Resign,
		connect.WithSchema(gameServiceResignMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	gameServicePassHandler := connect.NewUnaryHandler(
		GameServicePassProcedure,
		svc.Pass,
		connect.WithSchema(gameServicePassMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/gogame.v1.GameService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case GameServiceCreateGameProcedure:
			gameServiceCreateGameHandler.ServeHTTP(w, r)
		case GameServiceStreamVacantGamesProcedure:
			gameServiceStreamVacantGamesHandler.ServeHTTP(w, r)
		case GameServiceJoinGameProcedure:
			gameServiceJoinGameHandler.ServeHTTP(w, r)
		case GameServiceStreamGameEventsProcedure:
			gameServiceStreamGameEventsHandler.ServeHTTP(w, r)
		case GameServiceMakeMoveProcedure:
			gameServiceMakeMoveHandler.ServeHTTP(w, r)
		case GameServiceResignProcedure:
			gameServiceResignHandler.ServeHTTP(w, r)
		case GameServicePassProcedure:
			gameServicePassHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedGameServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedGameServiceHandler struct{}

func (UnimplementedGameServiceHandler) CreateGame(context.Context, *connect.Request[v1.CreateGameRequest]) (*connect.Response[v1.CreateGameResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("gogame.v1.GameService.CreateGame is not implemented"))
}

func (UnimplementedGameServiceHandler) StreamVacantGames(context.Context, *connect.Request[v1.StreamVacantGamesRequest], *connect.ServerStream[v1.StreamVacantGamesResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("gogame.v1.GameService.StreamVacantGames is not implemented"))
}

func (UnimplementedGameServiceHandler) JoinGame(context.Context, *connect.Request[v1.JoinGameRequest]) (*connect.Response[v1.JoinGameResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("gogame.v1.GameService.JoinGame is not implemented"))
}

func (UnimplementedGameServiceHandler) StreamGameEvents(context.Context, *connect.Request[v1.StreamGameEventsRequest], *connect.ServerStream[v1.StreamGameEventsResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("gogame.v1.GameService.StreamGameEvents is not implemented"))
}

func (UnimplementedGameServiceHandler) MakeMove(context.Context, *connect.Request[v1.MakeMoveRequest]) (*connect.Response[v1.MakeMoveResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("gogame.v1.GameService.MakeMove is not implemented"))
}

func (UnimplementedGameServiceHandler) Resign(context.Context, *connect.Request[v1.ResignRequest]) (*connect.Response[v1.ResignResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("gogame.v1.GameService.Resign is not implemented"))
}

func (UnimplementedGameServiceHandler) Pass(context.Context, *connect.Request[v1.PassRequest]) (*connect.Response[v1.PassResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("gogame.v1.GameService.Pass is not implemented"))
}
