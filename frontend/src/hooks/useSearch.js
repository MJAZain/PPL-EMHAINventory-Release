// hooks/useSearch.js
import { useState, useMemo } from "react";

export default function useSearch(data, keys = []) {
  const [searchTerm, setSearchTerm] = useState("");

  const filteredData = useMemo(() => {
    if (!searchTerm) return data;

    return data.filter((item) =>
      keys.some((key) =>
        String(item[key]).toLowerCase().includes(searchTerm.toLowerCase())
      )
    );
  }, [data, searchTerm, keys]);

  return {
    searchTerm,
    setSearchTerm,
    filteredData,
  };
}
