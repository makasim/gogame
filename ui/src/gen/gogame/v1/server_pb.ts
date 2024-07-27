// @generated by protoc-gen-es v1.10.0 with parameter "target=ts"
// @generated from file gogame/v1/server.proto (package gogame.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from enum gogame.v1.Color
 */
export enum Color {
  /**
   * @generated from enum value: COLOR_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * @generated from enum value: COLOR_BLACK = 1;
   */
  BLACK = 1,

  /**
   * @generated from enum value: COLOR_WHITE = 2;
   */
  WHITE = 2,
}
// Retrieve enum metadata with: proto3.getEnumType(Color)
proto3.util.setEnumType(Color, "gogame.v1.Color", [
  { no: 0, name: "COLOR_UNSPECIFIED" },
  { no: 1, name: "COLOR_BLACK" },
  { no: 2, name: "COLOR_WHITE" },
]);

/**
 * @generated from message gogame.v1.Player
 */
export class Player extends Message<Player> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: string name = 2;
   */
  name = "";

  /**
   * @generated from field: int32 captured_stones = 3;
   */
  capturedStones = 0;

  constructor(data?: PartialMessage<Player>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "gogame.v1.Player";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "captured_stones", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Player {
    return new Player().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Player {
    return new Player().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Player {
    return new Player().fromJsonString(jsonString, options);
  }

  static equals(a: Player | PlainMessage<Player> | undefined, b: Player | PlainMessage<Player> | undefined): boolean {
    return proto3.util.equals(Player, a, b);
  }
}

/**
 * @generated from message gogame.v1.Move
 */
export class Move extends Message<Move> {
  /**
   * @generated from field: string player_id = 1;
   */
  playerId = "";

  /**
   * @generated from field: gogame.v1.Color color = 2;
   */
  color = Color.UNSPECIFIED;

  /**
   * @generated from field: int32 x = 3;
   */
  x = 0;

  /**
   * @generated from field: int32 y = 4;
   */
  y = 0;

  constructor(data?: PartialMessage<Move>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "gogame.v1.Move";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "player_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "color", kind: "enum", T: proto3.getEnumType(Color) },
    { no: 3, name: "x", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
    { no: 4, name: "y", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Move {
    return new Move().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Move {
    return new Move().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Move {
    return new Move().fromJsonString(jsonString, options);
  }

  static equals(a: Move | PlainMessage<Move> | undefined, b: Move | PlainMessage<Move> | undefined): boolean {
    return proto3.util.equals(Move, a, b);
  }
}

/**
 * @generated from message gogame.v1.Game
 */
export class Game extends Message<Game> {
  /**
   * created
   *
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: int32 rev = 2;
   */
  rev = 0;

  /**
   * @generated from field: string name = 3;
   */
  name = "";

  /**
   * @generated from field: gogame.v1.Player player1 = 4;
   */
  player1?: Player;

  /**
   * @generated from field: gogame.v1.Player player2 = 5;
   */
  player2?: Player;

  /**
   * @generated from field: string state = 6;
   */
  state = "";

  /**
   * started
   *
   * @generated from field: gogame.v1.Move current_move = 7;
   */
  currentMove?: Move;

  /**
   * @generated from field: repeated gogame.v1.Move previous_moves = 8;
   */
  previousMoves: Move[] = [];

  /**
   * @generated from field: gogame.v1.Board board = 11;
   */
  board?: Board;

  /**
   * ended
   *
   * @generated from field: gogame.v1.Player winner = 9;
   */
  winner?: Player;

  /**
   * @generated from field: string won_by = 10;
   */
  wonBy = "";

  constructor(data?: PartialMessage<Game>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "gogame.v1.Game";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "rev", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
    { no: 3, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "player1", kind: "message", T: Player },
    { no: 5, name: "player2", kind: "message", T: Player },
    { no: 6, name: "state", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 7, name: "current_move", kind: "message", T: Move },
    { no: 8, name: "previous_moves", kind: "message", T: Move, repeated: true },
    { no: 11, name: "board", kind: "message", T: Board },
    { no: 9, name: "winner", kind: "message", T: Player },
    { no: 10, name: "won_by", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Game {
    return new Game().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Game {
    return new Game().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Game {
    return new Game().fromJsonString(jsonString, options);
  }

  static equals(a: Game | PlainMessage<Game> | undefined, b: Game | PlainMessage<Game> | undefined): boolean {
    return proto3.util.equals(Game, a, b);
  }
}

/**
 * @generated from message gogame.v1.Row
 */
export class Row extends Message<Row> {
  /**
   * @generated from field: repeated gogame.v1.Color colors = 1;
   */
  colors: Color[] = [];

  constructor(data?: PartialMessage<Row>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "gogame.v1.Row";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "colors", kind: "enum", T: proto3.getEnumType(Color), repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Row {
    return new Row().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Row {
    return new Row().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Row {
    return new Row().fromJsonString(jsonString, options);
  }

  static equals(a: Row | PlainMessage<Row> | undefined, b: Row | PlainMessage<Row> | undefined): boolean {
    return proto3.util.equals(Row, a, b);
  }
}

/**
 * @generated from message gogame.v1.Board
 */
export class Board extends Message<Board> {
  /**
   * @generated from field: int32 size = 1;
   */
  size = 0;

  /**
   * @generated from field: repeated gogame.v1.Row rows = 2;
   */
  rows: Row[] = [];

  constructor(data?: PartialMessage<Board>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "gogame.v1.Board";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "size", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
    { no: 2, name: "rows", kind: "message", T: Row, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Board {
    return new Board().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Board {
    return new Board().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Board {
    return new Board().fromJsonString(jsonString, options);
  }

  static equals(a: Board | PlainMessage<Board> | undefined, b: Board | PlainMessage<Board> | undefined): boolean {
    return proto3.util.equals(Board, a, b);
  }
}

/**
 * @generated from message gogame.v1.CreateGameRequest
 */
export class CreateGameRequest extends Message<CreateGameRequest> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: gogame.v1.Player player1 = 2;
   */
  player1?: Player;

  constructor(data?: PartialMessage<CreateGameRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "gogame.v1.CreateGameRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "player1", kind: "message", T: Player },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CreateGameRequest {
    return new CreateGameRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CreateGameRequest {
    return new CreateGameRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CreateGameRequest {
    return new CreateGameRequest().fromJsonString(jsonString, options);
  }

  static equals(a: CreateGameRequest | PlainMessage<CreateGameRequest> | undefined, b: CreateGameRequest | PlainMessage<CreateGameRequest> | undefined): boolean {
    return proto3.util.equals(CreateGameRequest, a, b);
  }
}

/**
 * @generated from message gogame.v1.CreateGameResponse
 */
export class CreateGameResponse extends Message<CreateGameResponse> {
  /**
   * @generated from field: gogame.v1.Game game = 1;
   */
  game?: Game;

  constructor(data?: PartialMessage<CreateGameResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "gogame.v1.CreateGameResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "game", kind: "message", T: Game },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CreateGameResponse {
    return new CreateGameResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CreateGameResponse {
    return new CreateGameResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CreateGameResponse {
    return new CreateGameResponse().fromJsonString(jsonString, options);
  }

  static equals(a: CreateGameResponse | PlainMessage<CreateGameResponse> | undefined, b: CreateGameResponse | PlainMessage<CreateGameResponse> | undefined): boolean {
    return proto3.util.equals(CreateGameResponse, a, b);
  }
}

/**
 * @generated from message gogame.v1.JoinGameRequest
 */
export class JoinGameRequest extends Message<JoinGameRequest> {
  /**
   * @generated from field: string game_id = 1;
   */
  gameId = "";

  /**
   * @generated from field: gogame.v1.Player player2 = 2;
   */
  player2?: Player;

  constructor(data?: PartialMessage<JoinGameRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "gogame.v1.JoinGameRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "game_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "player2", kind: "message", T: Player },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): JoinGameRequest {
    return new JoinGameRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): JoinGameRequest {
    return new JoinGameRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): JoinGameRequest {
    return new JoinGameRequest().fromJsonString(jsonString, options);
  }

  static equals(a: JoinGameRequest | PlainMessage<JoinGameRequest> | undefined, b: JoinGameRequest | PlainMessage<JoinGameRequest> | undefined): boolean {
    return proto3.util.equals(JoinGameRequest, a, b);
  }
}

/**
 * @generated from message gogame.v1.JoinGameResponse
 */
export class JoinGameResponse extends Message<JoinGameResponse> {
  /**
   * @generated from field: gogame.v1.Game game = 1;
   */
  game?: Game;

  constructor(data?: PartialMessage<JoinGameResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "gogame.v1.JoinGameResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "game", kind: "message", T: Game },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): JoinGameResponse {
    return new JoinGameResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): JoinGameResponse {
    return new JoinGameResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): JoinGameResponse {
    return new JoinGameResponse().fromJsonString(jsonString, options);
  }

  static equals(a: JoinGameResponse | PlainMessage<JoinGameResponse> | undefined, b: JoinGameResponse | PlainMessage<JoinGameResponse> | undefined): boolean {
    return proto3.util.equals(JoinGameResponse, a, b);
  }
}

/**
 * @generated from message gogame.v1.StreamVacantGamesRequest
 */
export class StreamVacantGamesRequest extends Message<StreamVacantGamesRequest> {
  constructor(data?: PartialMessage<StreamVacantGamesRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "gogame.v1.StreamVacantGamesRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): StreamVacantGamesRequest {
    return new StreamVacantGamesRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): StreamVacantGamesRequest {
    return new StreamVacantGamesRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): StreamVacantGamesRequest {
    return new StreamVacantGamesRequest().fromJsonString(jsonString, options);
  }

  static equals(a: StreamVacantGamesRequest | PlainMessage<StreamVacantGamesRequest> | undefined, b: StreamVacantGamesRequest | PlainMessage<StreamVacantGamesRequest> | undefined): boolean {
    return proto3.util.equals(StreamVacantGamesRequest, a, b);
  }
}

/**
 * @generated from message gogame.v1.StreamVacantGamesResponse
 */
export class StreamVacantGamesResponse extends Message<StreamVacantGamesResponse> {
  /**
   * @generated from field: gogame.v1.Game game = 1;
   */
  game?: Game;

  constructor(data?: PartialMessage<StreamVacantGamesResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "gogame.v1.StreamVacantGamesResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "game", kind: "message", T: Game },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): StreamVacantGamesResponse {
    return new StreamVacantGamesResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): StreamVacantGamesResponse {
    return new StreamVacantGamesResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): StreamVacantGamesResponse {
    return new StreamVacantGamesResponse().fromJsonString(jsonString, options);
  }

  static equals(a: StreamVacantGamesResponse | PlainMessage<StreamVacantGamesResponse> | undefined, b: StreamVacantGamesResponse | PlainMessage<StreamVacantGamesResponse> | undefined): boolean {
    return proto3.util.equals(StreamVacantGamesResponse, a, b);
  }
}

/**
 * @generated from message gogame.v1.StreamGameEventsRequest
 */
export class StreamGameEventsRequest extends Message<StreamGameEventsRequest> {
  /**
   * @generated from field: string game_id = 1;
   */
  gameId = "";

  constructor(data?: PartialMessage<StreamGameEventsRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "gogame.v1.StreamGameEventsRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "game_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): StreamGameEventsRequest {
    return new StreamGameEventsRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): StreamGameEventsRequest {
    return new StreamGameEventsRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): StreamGameEventsRequest {
    return new StreamGameEventsRequest().fromJsonString(jsonString, options);
  }

  static equals(a: StreamGameEventsRequest | PlainMessage<StreamGameEventsRequest> | undefined, b: StreamGameEventsRequest | PlainMessage<StreamGameEventsRequest> | undefined): boolean {
    return proto3.util.equals(StreamGameEventsRequest, a, b);
  }
}

/**
 * @generated from message gogame.v1.StreamGameEventsResponse
 */
export class StreamGameEventsResponse extends Message<StreamGameEventsResponse> {
  /**
   * @generated from field: gogame.v1.Game game = 1;
   */
  game?: Game;

  constructor(data?: PartialMessage<StreamGameEventsResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "gogame.v1.StreamGameEventsResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "game", kind: "message", T: Game },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): StreamGameEventsResponse {
    return new StreamGameEventsResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): StreamGameEventsResponse {
    return new StreamGameEventsResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): StreamGameEventsResponse {
    return new StreamGameEventsResponse().fromJsonString(jsonString, options);
  }

  static equals(a: StreamGameEventsResponse | PlainMessage<StreamGameEventsResponse> | undefined, b: StreamGameEventsResponse | PlainMessage<StreamGameEventsResponse> | undefined): boolean {
    return proto3.util.equals(StreamGameEventsResponse, a, b);
  }
}

/**
 * @generated from message gogame.v1.MakeMoveRequest
 */
export class MakeMoveRequest extends Message<MakeMoveRequest> {
  /**
   * @generated from field: string game_id = 1;
   */
  gameId = "";

  /**
   * @generated from field: int32 game_rev = 2;
   */
  gameRev = 0;

  /**
   * @generated from field: gogame.v1.Move move = 3;
   */
  move?: Move;

  constructor(data?: PartialMessage<MakeMoveRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "gogame.v1.MakeMoveRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "game_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "game_rev", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
    { no: 3, name: "move", kind: "message", T: Move },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MakeMoveRequest {
    return new MakeMoveRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MakeMoveRequest {
    return new MakeMoveRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MakeMoveRequest {
    return new MakeMoveRequest().fromJsonString(jsonString, options);
  }

  static equals(a: MakeMoveRequest | PlainMessage<MakeMoveRequest> | undefined, b: MakeMoveRequest | PlainMessage<MakeMoveRequest> | undefined): boolean {
    return proto3.util.equals(MakeMoveRequest, a, b);
  }
}

/**
 * @generated from message gogame.v1.MakeMoveResponse
 */
export class MakeMoveResponse extends Message<MakeMoveResponse> {
  /**
   * @generated from field: gogame.v1.Game game = 1;
   */
  game?: Game;

  constructor(data?: PartialMessage<MakeMoveResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "gogame.v1.MakeMoveResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "game", kind: "message", T: Game },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MakeMoveResponse {
    return new MakeMoveResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MakeMoveResponse {
    return new MakeMoveResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MakeMoveResponse {
    return new MakeMoveResponse().fromJsonString(jsonString, options);
  }

  static equals(a: MakeMoveResponse | PlainMessage<MakeMoveResponse> | undefined, b: MakeMoveResponse | PlainMessage<MakeMoveResponse> | undefined): boolean {
    return proto3.util.equals(MakeMoveResponse, a, b);
  }
}

/**
 * @generated from message gogame.v1.ResignRequest
 */
export class ResignRequest extends Message<ResignRequest> {
  /**
   * @generated from field: string game_id = 1;
   */
  gameId = "";

  /**
   * @generated from field: string player_id = 3;
   */
  playerId = "";

  constructor(data?: PartialMessage<ResignRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "gogame.v1.ResignRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "game_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "player_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ResignRequest {
    return new ResignRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ResignRequest {
    return new ResignRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ResignRequest {
    return new ResignRequest().fromJsonString(jsonString, options);
  }

  static equals(a: ResignRequest | PlainMessage<ResignRequest> | undefined, b: ResignRequest | PlainMessage<ResignRequest> | undefined): boolean {
    return proto3.util.equals(ResignRequest, a, b);
  }
}

/**
 * @generated from message gogame.v1.ResignResponse
 */
export class ResignResponse extends Message<ResignResponse> {
  /**
   * @generated from field: gogame.v1.Game game = 1;
   */
  game?: Game;

  constructor(data?: PartialMessage<ResignResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "gogame.v1.ResignResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "game", kind: "message", T: Game },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ResignResponse {
    return new ResignResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ResignResponse {
    return new ResignResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ResignResponse {
    return new ResignResponse().fromJsonString(jsonString, options);
  }

  static equals(a: ResignResponse | PlainMessage<ResignResponse> | undefined, b: ResignResponse | PlainMessage<ResignResponse> | undefined): boolean {
    return proto3.util.equals(ResignResponse, a, b);
  }
}

