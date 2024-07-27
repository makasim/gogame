import { Game } from "./gen/gogame/v1/server_pb";

type Props = {
  games: Game[];
  onJoin: (gameId: string) => void;
};

export const Games = ({ games, onJoin }: Props) => {
  if (games.length === 0) {
    return <h2>No games available</h2>;
  }

  return (
    <div>
      <h2>Games</h2>
      <ul>
        {games.map((game) => (
          <li key={game.id}>
            <span>{game.id} ({game.player1?.name})</span>
            <button onClick={() => onJoin(game.id)}>Join</button>
          </li>
        ))}
      </ul>
    </div>
  );
};
