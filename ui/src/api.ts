import { createPromiseClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { GameService } from "./gen/gogame/v1/server_connect";

const transport = createConnectTransport({
  baseUrl: import.meta.env.VITE_API_URL || "/",
});

export const client = createPromiseClient(GameService, transport);
