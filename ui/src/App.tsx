import { useEffect, useRef, useState } from "react";
import clsx from "clsx";
import { client } from "./api";
import { Game } from "./gen/gogame/v1/server_pb";

export function App() {
  const [turn, setTurn] = useState(1);
  const [size] = useState(19);
  const [board, setBoard] = useState<number[]>([]);
  const [currentGame, setCurrentGame] = useState<Game | null>(null);
  const [games, setGames] = useState<Game[]>([]);
  const nameField = useRef<HTMLInputElement>(null);

  function resetBoard() {
    setBoard(new Array(size * size).fill(0));
  }

  async function processGames() {
    for await (const { game, joinable } of client.streamVacantGames({})) {
      const exist = games.some((g) => g.id === game?.id);

      if (joinable && !exist) {
        setGames((prev) => [...prev, game!]);
      }

      if (!joinable && exist) {
        setGames((prev) => prev.filter((g) => g.id !== game?.id));
      }
    }
  }

  async function createGame() {
    const playerName = nameField.current?.value;

    if (!playerName) {
      alert("Name is required");
      return;
    }

    const { game } = await client.createGame({
      name: `Game-${Date.now()}`,
      player1: { id: playerName, name: playerName },
    });

    alert(`Game created: ${game?.name}`);
  }

  async function joinGame(selectedGame: Game) {
    const playerName = nameField.current?.value;

    if (!playerName) {
      alert("Name is required");
      return;
    }

    if (playerName === selectedGame.player1?.name) {
      alert("Names must be different");
      return;
    }

    const { game } = await client.joinGame({
      gameId: selectedGame.id,
      player2: { id: playerName, name: playerName },
    });

    if (!game) {
      alert("Game is not joinable");
      return;
    }

    alert(`Game joined: ${game?.name}`);

    setCurrentGame(game!);
    resetBoard();
  }

  useEffect(() => {
    processGames();
  }, []);

  useEffect(() => {
    if (!currentGame) return;

    const newBoard = new Array(size * size).fill(0);

    for (const move of currentGame.previousMoves) {
      const index = +move.y * 19 + Number(move.x);
      newBoard[index] = move.color;
    }

    setBoard(newBoard);
  }, [currentGame]);

  const putStone = (i: number) => {
    if (board[i] !== 0) return;

    const newBoard = [...board];
    newBoard[i] = turn;

    setTurn(turn === 1 ? 2 : 1);
    setBoard(newBoard);
  };

  return (
    <div className="App">
      Your name: <input ref={nameField} />
      <br />
      <button onClick={createGame}>Create Game</button>
      <ul>
        {games.map((game) => (
          <li key={game.id}>
            {game.name} ({game.player1?.name})
            <button onClick={() => joinGame(game)}>Join</button>
          </li>
        ))}
      </ul>

      <h1>Go game</h1>

      {currentGame && (
        <div className="field">
          {board.map((stone, i) => (
            <button
              onClick={() => putStone(i)}
              className={clsx("cell", {
                black: stone === 1,
                white: stone === 2,
              })}
              key={i}
            ></button>
          ))}
        </div>
      )}
    </div>
  );
}
