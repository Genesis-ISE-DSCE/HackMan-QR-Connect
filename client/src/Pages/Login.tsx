import React, { useState } from "react";
import connector from "../api";
import { useNavigate } from "react-router-dom";

function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const navigate = useNavigate();

  const handleSubmit = async (e: React.MouseEvent) => {
    e.preventDefault();
    try {
      const response = await connector.post("/auth", {
        username,
        password,
      });

      localStorage.setItem("token", response.data.token);
      const userId = localStorage.getItem("userId");
      navigate(`/dashboard/${userId}`);
    } catch (error) {
      console.error(error);
    }
  };
  return (
    <div className="w-screen h-screen flex flex-col  items-center justify-center bg-slate-800">
      <p className="font-extralight text-wrap text-center tracking-wider text-4xl md:text-6xl mb-4 bg-gradient-to-r from-purple-600 via-pink-600 to-red-600 bg-clip-text text-transparent p-4 animate-pulse">
        Login
      </p>
      <div className="flex flex-col max-w-xl">
        <label className="text-gray-300">Username</label>
        <input
          className="mt-1 bg-gray-200 rounded-md p-1 outline-none pl-2 shadow-md"
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
      </div>
      <div className="flex flex-col mt-3">
        <label className="text-gray-300">Password</label>
        <input
          className="mt-1 bg-gray-200 rounded-md p-1 outline-none pl-2 shadow-md"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <button
          onClick={handleSubmit}
          className="bg-gradient-to-r from-green-400 to-green-600 hover:from-green-500 hover:to-green-700 p-2 rounded-lg text-white mt-6 tracking-wide shadow-md transform transition-transform duration-300 ease-in-out hover:scale-105"
        >
          Submit
        </button>
      </div>
    </div>
  );
}

export default Login;
