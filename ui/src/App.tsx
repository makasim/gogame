import { useEffect, useState } from "react";
import clsx from "clsx";

import { client } from "./api";
import { Color, Game, State } from "./gen/gogame/v1/server_pb";
import { UserForm } from "./UserForm";
import { Games } from "./Games";

export function App() {
  const [playerId, setPlayerId] = useState(localStorage.getItem("currentUser"));
  const [availableGames, setAvailableGames] = useState<Game[]>([]);
  const [currentGame, setCurrentGame] = useState<Game | null>(() => {
    const game = localStorage.getItem("currentGame");
    return game ? JSON.parse(game) : null;
  });

  useEffect(() => {
    listenToAwailableGames();
  }, []);

  useEffect(() => {
    if (!playerId) return;
    localStorage.setItem("currentUser", playerId);
  }, [playerId]);

  useEffect(() => {
    if (!currentGame) return;
    localStorage.setItem("currentGame", JSON.stringify(currentGame));

    try {
      listenToGame(currentGame.id);
    } catch (error) {
      console.error("Game not found", error);
      resetGame();
    }
  }, [currentGame?.id]);

  if (!playerId) {
    return <UserForm onSave={setPlayerId} />;
  }

  if (!currentGame) {
    return (
      <>
        <button onClick={createGame}>Create Game</button>
        <Games games={availableGames} onJoin={joinGame} />
      </>
    );
  }

  if (currentGame.state === State.CREATED) {
    return <h2>Waiting for another player</h2>;
  }

  function resetGame() {
    setCurrentGame(null);
    localStorage.removeItem("currentGame");
  }

  async function createGame() {
    if (!playerId) return;

    const { game } = await client.createGame({
      name: `Game-${Date.now()}`,
      player1: { id: playerId, name: playerId },
    });

    if (!game) return alert("Game not created");

    setCurrentGame(game);
  }

  async function joinGame(gameId: string) {
    if (!playerId) return;

    const { game } = await client.joinGame({
      gameId,
      player2: { id: playerId, name: playerId },
    });

    if (!game) return alert("Game is not joinable");

    setCurrentGame(game);
  }

  async function listenToAwailableGames() {
    for await (const { game } of client.streamVacantGames({})) {
      if (!game) return alert("No games found");

      setAvailableGames((games) => {
        const filteredGames = games.filter((g) => g.id !== game.id);

        return game.state === State.CREATED && game.player1?.id !== playerId
          ? [game, ...filteredGames]
          : filteredGames;
      });
    }
  }

  async function listenToGame(gameId: string) {
    for await (const { game } of client.streamGameEvents({ gameId })) {
      if (!game) return alert("Game not found");
      setCurrentGame(game);
    }
  }

  async function putStone(i: number) {
    if (!playerId || !currentGame) return;
    if (currentGame.state === State.ENDED) return;
    if (currentGame.currentMove?.playerId !== playerId) return;

    const { game } = await client.makeMove({
      gameId: currentGame?.id,
      gameRev: currentGame?.rev,
      move: {
        ...currentGame.currentMove,
        x: i % 19,
        y: Math.floor(i / 19),
      },
    });

    if (!game) return alert("Move not made");

    setCurrentGame(game);
  }

  async function resign() {
    if (!playerId || !currentGame) return;
    const { game } = await client.resign({ gameId: currentGame.id, playerId });
    if (!game) return alert("Game not resigned");
    setCurrentGame(game);
  }

  const colors = currentGame.board?.rows.map((row) => row.colors).flat() || [];

  return (
    <div className="App">
      {currentGame.state === State.ENDED ? (
        <h2>
          Game ended.{" "}
          {currentGame.winner?.id === playerId ? "You won!" : "You lost!"}
          {currentGame.wonBy}
          <button onClick={resetGame}>Reset</button>
        </h2>
      ) : (
        <h2>
          <button onClick={resign}>Resign</button>

          {currentGame.currentMove?.playerId === playerId &&
            `You turn ${currentGame.currentMove.color === Color.BLACK ? "Black" : "White"}`}
        </h2>
      )}

      <div className="field">
        {colors.map((color, i) => (
          <button
            key={i}
            onClick={() => putStone(i)}
            className={clsx("cell", {
              black: color === Color.BLACK,
              white: color === Color.WHITE,
            })}
          ></button>
        ))}
      </div>
    </div>
  );
}
