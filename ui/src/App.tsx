import { useEffect, useState } from "react";
import clsx from "clsx";

import { client } from "./api";
import { Color, Game, State, Undo } from "./gen/gogame/v1/server_pb";
import { useNavigate, useParams } from "react-router-dom";

export function App() {
  const navigate = useNavigate();
  const { playerId, gameId } = useParams();
  const [currentGame, setCurrentGame] = useState<Game | null>(null);

  useEffect(() => {
    const abortController = new AbortController();

    try {
      listenToGameEvents(abortController.signal);
    } catch (error) {
      console.log("Game not found", error);
    }

    return () => abortController.abort();
  }, [gameId]);

  async function listenToGameEvents(signal: AbortSignal) {
    for await (const res of client.streamGameEvents({ gameId }, { signal })) {
      const { game, undo } = res;

      setCurrentGame((current) => {
        if (!game) return current;
        if (!current) return game;

        if (
          undo &&
          !undo.decided &&
          current?.currentMove?.playerId === playerId
        ) {
          if (confirm("Undo requested. Do you accept?")) acceptUndo(undo);
        }

        return game.rev > current.rev ? game : current;
      });
    }
  }

  function resetGame() {
    navigate(`/player/${playerId}`);
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

  async function requestUndo() {
    if (!playerId || !currentGame) return;
    if (currentGame.currentMove?.playerId === playerId) return;

    try {
      await client.undo({
        gameId: currentGame.id,
        gameRev: currentGame.rev,
        action: {
          value: { playerId },
          case: "request",
        },
      });
    } catch (error) {
      console.log("Undo not made", error);
      alert("Undo not made");
    }
  }

  async function acceptUndo(undo: Undo) {
    try {
      const { game } = await client.undo({
        gameId: undo.gameId,
        gameRev: undo.gameRev,
        action: {
          value: { playerId, accepted: true },
          case: "decision",
        },
      });

      if (!game) return alert("Undo not accepted");

      setCurrentGame(game);
    } catch (error) {
      console.log("Undo not accepted", error);
      alert("Undo not accepted");
    }
  }

  if (!playerId) return void navigate("/");
  if (!gameId) return void navigate(`/player/${playerId}`);
  if (!currentGame) {
    return (
      <h2>
        <button onClick={resetGame}>Reset</button>
        Loading...
      </h2>
    );
  }

  const { board, currentMove, state } = currentGame;

  if (state === State.CREATED) {
    return (
      <h2>
        <button onClick={resetGame}>Reset</button>
        Waiting for another player
      </h2>
    );
  }

  const colors = board?.rows.map((row) => row.colors).flat() || [];
  const yourTurn = currentMove?.playerId === playerId;
  const color = yourTurn
    ? currentMove?.color
    : currentMove?.color === Color.BLACK
      ? Color.WHITE
      : Color.BLACK;
  const colorName = color === Color.BLACK ? "black" : "white";

  return (
    <div className="App">
      {currentGame.state === State.ENDED ? (
        <h2>
          Game ended.{" "}
          {currentGame.winner?.id === playerId ? "You won!" : "You lost!"}
          <button onClick={resetGame}>Reset</button>
        </h2>
      ) : (
        <h2>
          <button onClick={resign}>Resign</button>
          Your color is {colorName}. {yourTurn ? "Your" : "Opponent's"} turn.
          {!yourTurn && <button onClick={requestUndo}>Undo</button>}
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
