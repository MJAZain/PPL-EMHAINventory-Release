import React, { useEffect, useState } from "react";
import SearchBar from "../../components/SearchBar";
import DataTable from "../../components/tableCompo";
import useSearch from "../../hooks/useSearch";
import Sidebar from "../../components/Sidebar";
import { apiClient } from "../../config/api";

function LaporanBarangKeluarPage() {
  const [incomingProducts, setIncomingProducts] = useState([]);

  // Fetch data from API
  const fetchIncomingProducts = async () => {
    try {
      const response = await apiClient.get("/outgoing-products/");
      return response.data.data; // Extract the actual data array
    } catch (error) {
      console.error("Failed to fetch incoming products:", error);
      throw error;
    }
  };

  // Load data on mount
  useEffect(() => {
    const loadData = async () => {
      try {
        const data = await fetchIncomingProducts();
        setIncomingProducts(Array.isArray(data) ? data : []);
      } catch {
        setIncomingProducts([]);
      }
    };

    loadData();
  }, []);

  // Search logic
  const { searchTerm, setSearchTerm, filteredData } = useSearch(
    incomingProducts,
    ["customer", "no_faktur", "payment_status", "date"]
  );

  // Define columns for table
  const columns = [
    { header: "Tanggal", accessor: "date" },
    { header: "Nama Pelanggan", accessor: "customer" },
    {
      header: "Total Pembelian Barang",
      accessor: "total_amount",
      render: (row) => {
        const value = row.total_amount;

        // Format the number as IDR
        const formatted = new Intl.NumberFormat("id-ID", {
          style: "currency",
          currency: "IDR",
          minimumFractionDigits: 0,
        }).format(value);

        // Replace default "Rp" with custom "Rp."
        return formatted.replace("Rp", "Rp.");
      },
    },
    { header: "NSFP", accessor: "no_faktur" },
    { header: "Status", accessor: "payment_status" },
  ];

  return (
    <div className="flex">
      <div className="bg-white min-h-screen">
            <Sidebar />
            </div>

      <div className="p-5 w-full py-10">
        <h1 className="text-2xl font-bold mb-6">Laporan Barang Terjual</h1>

        <SearchBar
          value={searchTerm}
          onChange={setSearchTerm}
          placeholder="Cari berdasarkan customer, faktur, status..."
        />

        <div className="border border-gray-400 rounded-xl bg-white p-5">
          <DataTable
            columns={columns}
            data={filteredData}
            showIndex={true}
          />
        </div>
      </div>
    </div>
  );
}

export default LaporanBarangKeluarPage;
