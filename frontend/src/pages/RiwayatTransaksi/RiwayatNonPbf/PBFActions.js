import { useState, useCallback } from "react";
import { apiClient } from "../../../config/api";

export default function usePBFActions() {
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

  const getPBFById = useCallback(
    (id) => {
      return handleRequest(() =>
        apiClient.get(`/incoming-nonpbf/${id}`).then((res) => res.data.data)
      );
    },
    [handleRequest]
  );

  const updatePBF = useCallback(
    (id, payload) => {
      return handleRequest(() =>
        apiClient.put(`/incoming-nonpbf/${id}`, payload).then((res) => res.data)
      );
    },
    [handleRequest]
  );

  const deletePBF = useCallback(
    (id) => {
      return handleRequest(() =>
        apiClient.delete(`/incoming-nonpbf/${id}`).then((res) => res.data)
      );
    },
    [handleRequest]
  );

  return {
    getPBFById,
    updatePBF,
    deletePBF,
    loading,
    error,
  };
}
