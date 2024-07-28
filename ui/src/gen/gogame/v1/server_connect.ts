// @generated by protoc-gen-connect-es v1.4.0 with parameter "target=ts"
// @generated from file gogame/v1/server.proto (package gogame.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { CreateGameRequest, CreateGameResponse, JoinGameRequest, JoinGameResponse, MakeMoveRequest, MakeMoveResponse, PassRequest, PassResponse, ResignRequest, ResignResponse, StreamGameEventsRequest, StreamGameEventsResponse, StreamVacantGamesRequest, StreamVacantGamesResponse } from "./server_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service gogame.v1.GameService
 */
export const GameService = {
  typeName: "gogame.v1.GameService",
  methods: {
    /**
     * @generated from rpc gogame.v1.GameService.CreateGame
     */
    createGame: {
      name: "CreateGame",
      I: CreateGameRequest,
      O: CreateGameResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc gogame.v1.GameService.StreamVacantGames
     */
    streamVacantGames: {
      name: "StreamVacantGames",
      I: StreamVacantGamesRequest,
      O: StreamVacantGamesResponse,
      kind: MethodKind.ServerStreaming,
    },
    /**
     * @generated from rpc gogame.v1.GameService.JoinGame
     */
    joinGame: {
      name: "JoinGame",
      I: JoinGameRequest,
      O: JoinGameResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc gogame.v1.GameService.StreamGameEvents
     */
    streamGameEvents: {
      name: "StreamGameEvents",
      I: StreamGameEventsRequest,
      O: StreamGameEventsResponse,
      kind: MethodKind.ServerStreaming,
    },
    /**
     * @generated from rpc gogame.v1.GameService.MakeMove
     */
    makeMove: {
      name: "MakeMove",
      I: MakeMoveRequest,
      O: MakeMoveResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc gogame.v1.GameService.Resign
     */
    resign: {
      name: "Resign",
      I: ResignRequest,
      O: ResignResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc gogame.v1.GameService.Pass
     */
    pass: {
      name: "Pass",
      I: PassRequest,
      O: PassResponse,
      kind: MethodKind.Unary,
    },
  }
} as const;

