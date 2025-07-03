import { useState, useCallback } from "react";
import { apiClient } from "../../config/api";

export default function useSupplierActions() {
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

  const getSupplierById = useCallback(
    (id) => {
      return handleRequest(() =>
        apiClient.get(`/suppliers/${id}`).then(res => res.data.data)
      );
    },
    [handleRequest]
  );

  const deleteSupplier = useCallback(
    (id) => {
      return handleRequest(() => apiClient.delete(`/suppliers/${id}`));
    },
    [handleRequest]
  );

  return {
    getSupplierById,
    deleteSupplier,
    loading,
    error,
  };
}
