import { useState, useCallback } from "react";
import { apiClient } from "../../config/api";

export default function useOpnameActions() {
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

  const getDrafts = useCallback(() => {
    return handleRequest(() =>
      apiClient.get(`/stock-opname/draft`).then(res =>
        res.data.data.filter(d => d.status === "draft")
      )
    );
  }, [handleRequest]);

  const getOpnameById = useCallback(
    (opnameId) => {
      return handleRequest(() =>
        apiClient.get(`/stock-opname/draft/${opnameId}`).then(res => res.data.data)
      );
    },
    [handleRequest]
  );

  const deleteOpname = useCallback(
    (opnameId) => {
      return handleRequest(() =>
        apiClient.delete(`/stock-opname/draft/${opnameId}`)
      );
    },
    [handleRequest]
  );

  return {
    getDrafts,
    getOpnameById,
    deleteOpname,
    loading,
    error,
  };
}
