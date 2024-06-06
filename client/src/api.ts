import axios from "axios";

const token = localStorage.getItem("token");
axios.defaults.headers.common["Authorization"] = `Bearer ${token}`;

const connector = axios.create({
  baseURL: "https://hackman-qr-connect-production.up.railway.app",
});

export default connector;
