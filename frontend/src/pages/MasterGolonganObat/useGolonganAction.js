import { useState, useCallback } from "react";
import { apiClient } from "../../config/api";

export default function useGolonganActions() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleRequest = useCallback(async (requestFn) => {
    setLoading(true);
    setError(null);
    try {
      const result = await requestFn();
      return result;
    } catch (err) {
      setError(err.response?.data?.message || err.message || "Unknown error");
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const getGolonganById = useCallback(
    (id) => {
      return handleRequest(() =>
        apiClient.get(`/drug-categories/${id}`).then((res) => res.data.data)
      );
    },
    [handleRequest]
  );

  const deleteGolongan = useCallback(
    (id) => {
      return handleRequest(() =>
        apiClient.delete(`/drug-categories/${id}`)
      );
    },
    [handleRequest]
  );

  return {
    getGolonganById,
    deleteGolongan,
    loading,
    error,
  };
}
