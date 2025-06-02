import { useState, useCallback } from "react";
import { apiClient } from "../config/api";

export default function useUnitActions() {
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

  // Memoize getUnitById to keep function identity stable
  const getUnitById = useCallback(
    (id) => {
      return handleRequest(() =>
        apiClient.get(`/units/${id}`).then((res) => res.data.data)
      );
    },
    [handleRequest]
  );

  const deleteUnit = useCallback(
    (id) => {
      return handleRequest(() => apiClient.delete(`/units/${id}`));
    },
    [handleRequest]
  );

  return {
    getUnitById,
    deleteUnit,
    loading,
    error,
  };
}
