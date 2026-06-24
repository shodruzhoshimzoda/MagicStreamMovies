import { useContext } from "react";
import AuthContext from "../context/AuthProvider";

const useAuth = () => {
    // ВАЖНО: вызываем useContext, а не useAuth!
    return useContext(AuthContext);
}

export default useAuth;
