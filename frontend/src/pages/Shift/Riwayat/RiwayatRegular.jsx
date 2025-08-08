import React, { useState, useEffect } from "react";
import useSearch from "../../../hooks/useSearch";
import { apiClient } from "../../../config/api";
import SearchBar from "../../../components/SearchBar";
import DataTable from "../../../components/tableCompo";
import Toast from '../../../components/toast';
import Sidebar from "../../../components/Sidebar";

function RiwayatRegularPage() {
  const [shifts, setShifts] = useState([]);
  const [toast, setToast] = useState(null);
  const [loading, setLoading] = useState(true);

  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 10;
  const [totalCount, setTotalCount] = useState(0);

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

  const fetchShifts = async (page = 1) => {
    setLoading(true);
    try {
      const offset = (page - 1) * pageSize;
      const response = await apiClient.get(
        `/sales/regular?limit=${pageSize}&offset=${offset}`
      );

      const formattedData = (response.data.data || []).map(sales => ({
        ...sales,
        transaction_date: formatDateTime(sales.transaction_date),
      }));

      setShifts(formattedData);
      setTotalCount(response.data.total || 0);
    } catch (err) {
      setToast({ message: "Data gagal diambil", type: "error" });
      console.error("Fetch error:", err);
      setShifts([]);
      setTotalCount(0);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchShifts(currentPage);
  }, [currentPage]);

  const pageCount = Math.ceil(totalCount / pageSize);

  const columns = [
    { header: "Nama Pegawai", accessor: "cashier_name" },
    { header: "Total Pembelian", accessor: "total_pay" },
    { header: "Tanggal Pembelian", accessor: "transaction_date" },
    { header: "Cara Pembayaran", accessor: "payment_method" },
  ];

  const changePage = (page) => {
    if (page < 1 || page > pageCount) return;
    setCurrentPage(page);
  };

  const renderPageNumbers = () => {
    const pages = [];
    const maxShown = 5;

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
        <h1 className="text-2xl font-bold mb-6">Riwayat Penjualan Bebas</h1>
        <SearchBar value={""} onChange={() => {}} disabled />

        <div className="border-1 rounded-md border-gray-300 bg-white p-5">
          <div className="max-w-[1121px]">
            <DataTable columns={columns} data={shifts} showIndex={true} />
          </div>

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

export default RiwayatRegularPage;
