import { useNavigate, useParams } from "react-router-dom";
import { Games } from "./Games";
import { client } from "./api";
import { useEffect, useState } from "react";
import { Game, State } from "./gen/gogame/v1/server_pb";

export const GamesPage = () => {
  const [duration, setDuration] = useState<number>(60);
  const navigate = useNavigate();
  const { playerId } = useParams();

  const [availableGames, setAvailableGames] = useState<Game[]>([]);

  useEffect(() => {
    const abortController = new AbortController();

    try {
      listenToAvailableGames(abortController.signal);
    } catch (error) {
      console.log("No games found", error);
    }

    return () => abortController.abort();
  }, [playerId]);

  async function listenToAvailableGames(signal: AbortSignal) {
    for await (const { game } of client.streamVacantGames({}, { signal })) {
      if (!game) {
        alert("No games found");
        continue;
      }

      setAvailableGames((games) => {
        const filteredGames = games.filter((g) => g.id !== game.id);

        return game.state === State.CREATED && game.player1?.id !== playerId
          ? [game, ...filteredGames]
          : filteredGames;
      });
    }
  }

  async function joinGame(gameId: string) {
    const { game } = await client.joinGame({
      gameId,
      player2: { id: playerId, name: playerId },
    });

    if (!game) return alert("Game is not joinable");

    navigate(`/player/${playerId}/game/${game.id}`);
  }

  async function createGame() {
    const { game } = await client.createGame({
      name: `Game-${Date.now()}`,
      player1: { id: playerId, name: playerId },
      moveDurationSec: duration,
    });

    if (!game) return alert("Can't create a game");

    navigate(`/player/${playerId}/game/${game.id}`);
  }

  if (!playerId) return void navigate("/");

  return (
    <>
      <input
        type="number"
        min="5"
        defaultValue={duration}
        onChange={(e) => setDuration(+e.target.value)}
      />
      <button onClick={createGame}>Create Game</button>
      <Games games={availableGames} onJoin={joinGame} />
    </>
  );
};
