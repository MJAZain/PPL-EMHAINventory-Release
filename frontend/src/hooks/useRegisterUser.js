import { useState } from "react";
import axios from "axios";

export default function useRegisterUser() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const registerUser = async (form) => {
    try {
      setLoading(true);
      const res = await axios.post(
        "/api/users/register",
        form,
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
        }
      );
      return res.data;
    } catch (err) {
      setError(err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return {
    registerUser,
    loading,
    error,
  };
}
