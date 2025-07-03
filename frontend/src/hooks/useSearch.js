import { useState, useMemo } from "react";

function getNestedValue(obj, path) {
  return path.split(".").reduce((o, key) => (o ? o[key] : ""), obj);
}

export default function useSearch(data, keys = []) {
  const [searchTerm, setSearchTerm] = useState("");

  const filteredData = useMemo(() => {
    if (!searchTerm) return data;

    return data.filter((item) =>
      keys.some((key) =>
        String(getNestedValue(item, key)).toLowerCase().includes(searchTerm.toLowerCase())
      )
    );
  }, [data, searchTerm, keys]);

  return {
    searchTerm,
    setSearchTerm,
    filteredData,
  };
}
