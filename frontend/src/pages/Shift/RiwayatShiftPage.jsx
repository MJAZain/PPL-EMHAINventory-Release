import React, { useState, useEffect } from "react";
import useSearch from "../../hooks/useSearch";
import { apiClient } from "../../config/api";
import SearchBar from "../../components/SearchBar";
import DataTable from "../../components/tableCompo";
import Toast from '../../components/toast';
import Sidebar from "../../components/Sidebar";

function RiwayatShiftPage() {
  const [shifts, setShifts] = useState([]);
  const [toast, setToast] = useState(null);
  const [loading, setLoading] = useState(true);

  // pagination state
  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 5;

  const formatDateTime = (isoString) => {
    if (!isoString) return "-";
    const date = new Date(isoString);
    return date.toLocaleString("id-ID", {
      day: "2-digit",
      month: "short",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const fetchShift = async () => {
    try {
      const response = await apiClient.get("/shifts/");

      const formattedData = (response.data.data || []).map(shift => ({
        ...shift,
        opening_time: formatDateTime(shift.opening_time),
        closing_time: formatDateTime(shift.closing_time),
      }));

      setShifts(formattedData);
    } catch (err) {
      setToast({ message: "Data gagal diambil", type: "error" });
      console.error("Fetch error:", err);
      setShifts([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchShift();
  }, []);

  const { searchTerm, setSearchTerm, filteredData } = useSearch(shifts, [
    "closing_balance",
    "opening_time",
    "closing_time",
    "total_sales",
  ]);

  useEffect(() => {
    setCurrentPage(1);
  }, [searchTerm, shifts]);

  const pageCount = Math.ceil(filteredData.length / pageSize);

  const paginatedData = filteredData.slice(
    (currentPage - 1) * pageSize,
    currentPage * pageSize
  );

  const columns = [
    { header: "Closing Balance", accessor: "closing_balance",
      render: (row) => {
        const value = row.closing_balance;
        const formatted = new Intl.NumberFormat("id-ID", {
          style: "currency",
          currency: "IDR",
          minimumFractionDigits: 0,
        }).format(value);
        return formatted.replace("Rp", "Rp.");
      },
     },
    { header: "Opening Time", accessor: "opening_time" },
    { header: "Closing Time", accessor: "closing_time" },
    { header: "Total Sales", accessor: "total_sales" },
  ];

  const changePage = (page) => {
    if (page < 1 || page > pageCount) return;
    setCurrentPage(page);
  };

  const renderPageNumbers = () => {
    const pages = [];
    const maxShown = 5; // max page buttons to show

    if (pageCount <= maxShown) {
      for (let i = 1; i <= pageCount; i++) {
        pages.push(i);
      }
    } else {
      if (currentPage <= 3) {
        pages.push(1, 2, 3, '...', pageCount);
      } else if (currentPage >= pageCount - 2) {
        pages.push(1, '...', pageCount - 2, pageCount - 1, pageCount);
      } else {
        pages.push(1, '...', currentPage, '...', pageCount);
      }
    }

    return pages.map((page, index) =>
      page === '...' ? (
        <span key={index} className="px-2">…</span>
      ) : (
        <button
          key={index}
          onClick={() => changePage(page)}
          className={`px-3 py-1 rounded ${page === currentPage ? 'bg-blue-500 text-white' : 'bg-gray-200'}`}
        >
          {page}
        </button>
      )
    );
  };

  if (loading) return <div className="text-center mt-6">Loading...</div>;

  return (
    <div className="flex">
      <div className="bg-white min-h-screen">
        <Sidebar />
      </div>

      <div className="p-5 w-full py-10">
        <h1 className="text-2xl font-bold mb-6">Riwayat Shift Resep</h1>
        <SearchBar value={searchTerm} onChange={setSearchTerm} />

        <div className="border-1 rounded-md border-gray-300 bg-white p-5">
          <div className="max-w-[1121px]">
            <DataTable columns={columns} data={paginatedData} showIndex={false} />
          </div>

          {/* Pagination Controls */}
          <div className="flex justify-between items-center mt-4">
            <div>
              Showing page {currentPage} of {pageCount} pages
            </div>

            <div className="flex gap-1">
              <button
                onClick={() => changePage(1)}
                disabled={currentPage === 1}
                className="px-2 py-1 bg-gray-200 rounded disabled:opacity-50"
              >
                &laquo;
              </button>
              <button
                onClick={() => changePage(currentPage - 1)}
                disabled={currentPage === 1}
                className="px-2 py-1 bg-gray-200 rounded disabled:opacity-50"
              >
                &lt;
              </button>

              {renderPageNumbers()}

              <button
                onClick={() => changePage(currentPage + 1)}
                disabled={currentPage === pageCount}
                className="px-2 py-1 bg-gray-200 rounded disabled:opacity-50"
              >
                &gt;
              </button>
              <button
                onClick={() => changePage(pageCount)}
                disabled={currentPage === pageCount}
                className="px-2 py-1 bg-gray-200 rounded disabled:opacity-50"
              >
                &raquo;
              </button>
            </div>
          </div>
        </div>
      </div>

      {toast && (
        <Toast
          message={toast.message}
          type={toast.type}
          onClose={() => setToast(null)}
        />
      )}
    </div>
  );
}

export default RiwayatShiftPage;
