import { useState, useCallback } from "react";
import { apiClient } from "../../config/api";

export default function useDoctorActions() {
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

  const getDoctorById = useCallback(
    (id) => {
      return handleRequest(() =>
        apiClient.get(`/doctors/${id}`).then(res => res.data.data)
      );
    },
    [handleRequest]
  );

  const deleteDoctor = useCallback(
    (id) => {
      return handleRequest(() => apiClient.delete(`/doctors/${id}`));
    },
    [handleRequest]
  );

  return {
    getDoctorById,
    deleteDoctor,
    loading,
    error,
  };
}
