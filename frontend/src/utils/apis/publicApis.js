import axios from "axios";

const baseURL =
  process.env.NODE_ENV === "production"
    ? "http://localhost:8098"
    : "http://localhost:8098";

export default axios.create({
  baseURL,
  headers: {
    accept: "application/json",
    "Content-Type": "application/json",
  },
});
