import { useState, useCallback } from "react";
import { apiClient } from "../../config/api";

export default function useKaryawanActions() {
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

  const getUserById = useCallback(
    async (id) => {
      return handleRequest(async () => {
        const res = await apiClient.get(`/users/${id}`);
        return res.data;
      });
    },
    [handleRequest]
  );

  const reactivateUser = useCallback(
    (id) => {
      return handleRequest(() => apiClient.patch(`/users/${id}/reactivate`));
    },
    [handleRequest]
  );

  const deactivateUser = useCallback(
    async (id) => {
      console.log("Attempting to deactivate user with ID:", id);

      return handleRequest(async () => {
        try {
          const response = await apiClient.patch(`/users/${id}/deactivate`);
          console.log("Deactivation successful:", response.data);
          return response.data;
        } catch (error) {
          console.error("Error during deactivation:", error);
          throw error;
        }
      });
    },
    [handleRequest]
  );

  const deleteUser = useCallback(
    (id) => {
      return handleRequest(() => apiClient.delete(`/users/${id}`));
    },
    [handleRequest]
  );

  const exportUsersCSV = useCallback(() => {
    return handleRequest(() =>
      apiClient.get("/users/export/csv", {
        responseType: "blob", // Important if you plan to download the file
      })
    );
  }, [handleRequest]);

  return {
    getUserById,
    deleteUser,
    reactivateUser,
    deactivateUser,
    exportUsersCSV,
    loading,
    error,
  };
}
