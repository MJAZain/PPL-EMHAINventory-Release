import { useState } from "react";
import axios from "axios";
import { apiClient } from "../config/api";

export default function useProductActions() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const getProductById = async (id) => {
    setLoading(true);
    setError(null);
    try {
      const res = await apiClient.get(`/products/${id}`);
      return res.data;
    } catch (err) {
      setError(err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const deleteProduct = async (id) => {
    setLoading(true);
    setError(null);
    try {
      await apiClient.delete(`/products/${id}`);
      return true;
    } catch (err) {
      setError(err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return {
    getProductById,
    deleteProduct,
    loading,
    error,
  };
}

