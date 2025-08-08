import { useState, useCallback } from "react";
import { apiClient } from "../../config/api";

export default function useJenisActions() {
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

  const getJenisById = useCallback(
    (id) => {
      return handleRequest(() =>
        apiClient.get(`/expense-types/${id}`).then(res => res.data.data)
      );
    },
    [handleRequest]
  );

  const deleteJenis = useCallback(
    (id) => {
      return handleRequest(() => apiClient.delete(`/expense-types/${id}`));
    },
    [handleRequest]
  );

  return {
    getJenisById,
    deleteJenis,
    loading,
    error,
  };
}
