import { useEffect, useState } from "react";
import clsx from "clsx";
import { client } from "./api";

export function App() {
  const [turn, setTurn] = useState(1);
  const [size] = useState(19);
  const [board, setBoard] = useState<number[]>(Array(size ** 2).fill(0));

  async function processGames() {
    for await (const res of client.streamVacantGames({})) {
      console.log(res);
    }
  }

  async function createGame() {
    const res = await client.createGame({});
    console.log(res);
  }

  useEffect(() => {
    processGames();
  }, []);

  const putStone = (i: number) => {
    if (board[i] !== 0) return;

    const newBoard = [...board];
    newBoard[i] = turn;

    setTurn(turn === 1 ? 2 : 1);
    setBoard(newBoard);
  };

  return (
    <div className="App">
      <button onClick={createGame}>Create Game</button>
      <h1>Go game</h1>

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
    </div>
  );
}
