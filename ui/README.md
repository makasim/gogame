# Go Game - React Client

This is a React cient for https://github.com/makasim/gogame.

You need to have [Node.js LTS](https://nodejs.org/en) to build or run the UI.

It also uses `pnpm`. To install it globally run:

```shell
npm i -g pnpm
```

> All the next commands should be executed from the `ui` folder.

To install deps run:

```shell
pnpm i
```

To build the UI for Go server run:

```shell
pnpm build
```

Run back and dev front servers:
```shell
pnpm dev
```

```shell
CORS_ENABLED=true go run main/main.go
```

Open the page at http://localhost:5173/, back is at http://localhost:8181
