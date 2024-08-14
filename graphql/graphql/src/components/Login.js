import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      const response = await fetch("https://01.gritlab.ax/api/auth/signin", {
        method: "POST",
        headers: {
          Authorization: `Basic ${btoa(`${username}:${password}`)}`,
        },
      });

      if (response.ok) {
        const data = await response.json();
        const cleanedToken = data.replace(/^"(.*)"$/, "$1");
        localStorage.setItem("token", cleanedToken);
        console.log(localStorage.getItem("token"));
        navigate("/profile");
      } else {
        setError("Invalid credentials");
      }
    } catch (err) {
      setError("Error while logging in");
    }
  };

  return (
    <div>
      <h1>Login</h1>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="Username or Email"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <button type="submit">Login</button>
      </form>
      {error && <p id="error">{error}</p>}
    </div>
  );
}

export default Login;
