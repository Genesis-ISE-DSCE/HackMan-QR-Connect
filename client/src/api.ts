import axios from "axios";

const token = localStorage.getItem("token");
axios.defaults.headers.common["Authorization"] = `Bearer ${token}`;

const connector = axios.create({
  baseURL: "http://localhost:7500",
});

export default connector;
