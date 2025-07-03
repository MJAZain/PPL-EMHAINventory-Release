import { useState, useCallback } from "react";
import { apiClient } from "../../config/api";

export default function useSatuanActions() {
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

  const getSatuanById = useCallback(
    (id) => {
      return handleRequest(() =>
        apiClient.get(`/units/${id}`).then((res) => res.data.data)
      );
    },
    [handleRequest]
  );

  const deleteSatuan = useCallback(
    (id) => {
      return handleRequest(() => apiClient.delete(`/units/${id}`));
    },
    [handleRequest]
  );

  return {
    getSatuanById,
    deleteSatuan,
    loading,
    error,
  };
}
