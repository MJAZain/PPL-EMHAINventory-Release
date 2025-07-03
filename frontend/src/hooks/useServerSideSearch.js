import { useState, useEffect, useCallback } from "react";
import { apiClient } from "../config/api";

const ITEMS_PER_PAGE_DEFAULT = 10;
const SEARCH_DEBOUNCE_DELAY_DEFAULT = 500;

function useServerSideSearch(
  apiPath,
  initialItemsPerPage = ITEMS_PER_PAGE_DEFAULT,
  debounceDelay = SEARCH_DEBOUNCE_DELAY_DEFAULT
) {
  const [items, setItems] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const [totalData, setTotalData] = useState(0);
  const [searchTerm, setSearchTerm] = useState("");
  const [searchInput, setSearchInput] = useState("");

  useEffect(() => {
    const timerId = setTimeout(() => {
      setSearchTerm(searchInput);
      setCurrentPage(1);
    }, debounceDelay);
    return () => clearTimeout(timerId);
  }, [searchInput, debounceDelay]);

  const fetchData = useCallback(
    async (pageToFetch, currentSearchTerm) => {
      setIsLoading(true);
      try {
        const response = await apiClient.get(
          `${apiPath}?page=${pageToFetch}&limit=${initialItemsPerPage}&search=${encodeURIComponent(
            currentSearchTerm
          )}`
        );
        if (
          response.data &&
          response.data.data &&
          Array.isArray(response.data.data.data)
        ) {
          setItems(response.data.data.data);
          setTotalData(response.data.data.total_data || 0);
          setTotalPages(response.data.data.total_pages || 0);
          setCurrentPage(response.data.data.current_page || 1);
        } else {
          console.error(
            `Unexpected data format from ${apiPath}:`,
            response.data
          );
          setItems([]);
          setTotalData(0);
          setTotalPages(0);
          setCurrentPage(1);
        }
      } catch (error) {
        console.error(`Failed to fetch data from ${apiPath}`, error);
        setItems([]);
        setTotalData(0);
        setTotalPages(0);
        setCurrentPage(1);
      } finally {
        setIsLoading(false);
      }
    },
    [apiPath, initialItemsPerPage]
  );

  useEffect(() => {
    fetchData(currentPage, searchTerm);
  }, [fetchData, currentPage, searchTerm]);

  const handleSearchInputChange = (newInput) => {
    setSearchInput(newInput);
  };

  const handlePageChange = (newPage) => {
    if (
      newPage >= 1 &&
      (newPage <= totalPages || totalPages === 0) &&
      newPage !== currentPage
    ) {
      setCurrentPage(newPage);
    }
  };

  const refetch = useCallback(() => {
    fetchData(currentPage, searchTerm);
  }, [fetchData, currentPage, searchTerm]);

  return {
    items,
    isLoading,
    currentPage,
    totalPages,
    totalData,
    searchInput,
    searchTerm,
    handleSearchInputChange,
    handlePageChange,
    setCurrentPage,
    ITEMS_PER_PAGE: initialItemsPerPage,
    refetch,
  };
}

export default useServerSideSearch;
