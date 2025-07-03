import { apiClient } from "../config/api";
import { useCallback } from "react";

const useStorageLocationActions = () => {
  const getStorageLocationById = useCallback(async (id) => {
    try {
      const response = await apiClient.get(`/storage-locations/${id}`);
      return response.data.data;
    } catch (error) {
      console.error("Failed to fetch storage location by ID", error);
      throw error;
    }
  }, []);

  const addStorageLocation = useCallback(async (data) => {
    try {
      const response = await apiClient.post("/storage-locations", data);
      return response.data;
    } catch (error) {
      console.error("Failed to add storage location", error);
      throw error;
    }
  }, []);

  const updateStorageLocation = useCallback(async (id, data) => {
    try {
      const response = await apiClient.put(`/storage-locations/${id}`, data);
      return response.data;
    } catch (error) {
      console.error("Failed to update storage location", error);
      throw error;
    }
  }, []);

  const deleteStorageLocation = useCallback(async (id) => {
    try {
      const response = await apiClient.delete(`/storage-locations/${id}`);
      return response.data;
    } catch (error) {
      console.error("Failed to delete storage location", error);
      throw error;
    }
  }, []);

  return {
    getStorageLocationById,
    addStorageLocation,
    updateStorageLocation,
    deleteStorageLocation,
  };
};

export default useStorageLocationActions;
