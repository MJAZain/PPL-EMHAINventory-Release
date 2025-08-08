import { useState, useCallback } from "react";
import { apiClient } from "../../config/api";

export default function usePengeluaranActions() {
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

  const getPengeluaranById = useCallback(
    (id) => {
      return handleRequest(() =>
        apiClient.get(`/expenses/${id}`).then(res => res.data.data)
      );
    },
    [handleRequest]
  );

  const deletePengeluaran = useCallback(
    (id) => {
      return handleRequest(() => apiClient.delete(`/expenses/${id}`));
    },
    [handleRequest]
  );

  return {
    getPengeluaranById,
    deletePengeluaran,
    loading,
    error,
  };
}
