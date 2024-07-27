import { useState } from "react";

type Props = {
  onSave: (userName: string) => void;
};

export const UserForm = ({ onSave }: Props) => {
  const [user, setUser] = useState("");

  return (
    <div>
      <input
        type="text"
        placeholder="Enter your name"
        onChange={(e) => setUser(e.target.value)}
      />
      <button
        onClick={() => {
          if (!user) return;

          onSave(user);
        }}
      >
        Submit
      </button>
    </div>
  );
};
