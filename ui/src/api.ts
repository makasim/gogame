import { createPromiseClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { GameService } from "./gen/gogame/v1/server_connect";

const transport = createConnectTransport({
  baseUrl: "http://localhost:8181",
});

export const client = createPromiseClient(GameService, transport);
