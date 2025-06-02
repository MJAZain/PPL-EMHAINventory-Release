import { useEffect, useState } from "react";
import { apiClient } from "../config/api"; // adjust the path if needed

export default function useBarangData() {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchProducts = async () => {
    setLoading(true);
    try {
      const response = await apiClient.get("/products/");
      if (response.status !== 200 || !Array.isArray(response.data.data)) {
        throw new Error(response.data?.message || "Failed to fetch products.");
      }
      setData(response.data.data);
    } catch (err) {
      console.error("Fetch error:", err);
      setError(err.message || "An unexpected error occurred.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProducts();
  }, []);

  return { data, loading, error, refetch: fetchProducts };
}
