import { Navigate } from "react-router-dom";

export default function PrivateRoute({ children }) {
  let auth = useAuth();
  return auth ? children : <Navigate to="/login" replace={true} />;
}

export const useAuth = () => {
  const auth = localStorage.getItem("token");
  return auth === "true" ? true : false;
};
