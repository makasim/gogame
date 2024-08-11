import { useState } from "react";
import { useNavigate } from "react-router-dom";

export const UserForm = () => {
  const [playerId, setPlayerId] = useState("");
  const navigate = useNavigate();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!playerId) return alert("Please enter your name");
    navigate(`/player/${playerId}`);
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        placeholder="Enter your name"
        onChange={(e) => setPlayerId(e.target.value)}
      />
      <button>Submit</button>
    </form>
  );
};
