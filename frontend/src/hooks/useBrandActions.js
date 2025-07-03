import { apiClient } from "../config/api";
import { useCallback } from "react";

const useBrandActions = () => {
  const getBrandById = useCallback(async (id) => {
    try {
      const response = await apiClient.get(`/brands/${id}`);
      return response.data.data;
    } catch (error) {
      console.error("Failed to fetch brand by ID", error);
      throw error;
    }
  }, []);

  const addBrand = useCallback(async (data) => {
    try {
      const response = await apiClient.post("/brands", data);
      return response.data;
    } catch (error) {
      console.error("Failed to add brand", error);
      throw error;
    }
  }, []);

  const updateBrand = useCallback(async (id, data) => {
    try {
      const response = await apiClient.put(`/brands/${id}`, data);
      return response.data;
    } catch (error) {
      console.error("Failed to update brand", error);
      throw error;
    }
  }, []);

  const deleteBrand = useCallback(async (id) => {
    try {
      const response = await apiClient.delete(`/brands/${id}`);
      return response.data;
    } catch (error) {
      console.error("Failed to delete brand", error);
      throw error;
    }
  }, []);

  return {
    getBrandById,
    addBrand,
    updateBrand,
    deleteBrand,
  };
};

export default useBrandActions;
