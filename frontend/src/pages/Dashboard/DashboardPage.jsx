import { useEffect, useState } from "react";
import { Card, CardContent } from "../../components/cardComp";
import { apiClient } from "../../config/api";
import DataTable from "../../components/tableCompo";
import Sidebar from "../../components/Sidebar";
import { formatToIDR, formatDateTime } from "../../utils/formatter";

export default function Dashboard() {
  const [summary, setSummary] = useState(null);
  const [stockSummary, setStockSummary] = useState(null);
  const [lowStockData, setLowStockData] = useState([]);
  const [expiringSoonData, setExpiringSoonData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  useEffect(() => {
    const fetchData = async () => {
      try {
        const [
          dashboardRes,
          stocksRes,
          lowStockRes,
          expiringSoonRes,
        ] = await Promise.all([
          apiClient.get("/dashboard/summary"),
          apiClient.get("/stocks/summary"),
          apiClient.get("/stocks/low"),
          apiClient.get("/stocks/expiring-soon"),
        ]);

        if (
          dashboardRes.data.status &&
          stocksRes.data.status === 200
        ) {
          setSummary(dashboardRes.data.data);
          setStockSummary(stocksRes.data.data);
          setLowStockData(lowStockRes.data.data || []);
          setExpiringSoonData(expiringSoonRes.data.data || []);
        } else {
          setError("Failed to fetch one or more dashboard sections.");
        }
      } catch (err) {
        console.error(err);
        setError("An error occurred while fetching data.");
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) {
    return <div className="p-4 text-center">Loading...</div>;
  }

  if (error) {
    return <div className="p-4 text-center text-red-500">{error}</div>;
  }

  return (
<div className="flex">
  <div className="bg-white min-h-screen">
    <Sidebar />
  </div>

  <div className="p-5 w-full py-10">
    <h1 className="text-2xl font-bold mb-6">Dashboard</h1>

{/* Row 1: Sales */}
<div className="grid grid-cols-1 sm:grid-cols-3 gap-4 m-5">
  <Card>
    <CardContent className="p-6 text-center">
      <h2 className="text-lg font-semibold">Penjualan Reguler</h2>
      <p className="text-2xl font-bold text-blue-600">
        {formatToIDR(summary.total_sales_regular)}
      </p>
    </CardContent>
  </Card>

  <Card>
    <CardContent className="p-6 text-center">
      <h2 className="text-lg font-semibold">Penjualan Resep</h2>
      <p className="text-2xl font-bold text-green-600">
        {formatToIDR(summary.total_sales_prescription)}
      </p>
    </CardContent>
  </Card>

  <Card>
    <CardContent className="p-6 text-center">
      <h2 className="text-lg font-semibold">Total Omzet</h2>
      <p className="text-2xl font-bold text-purple-600">
        {formatToIDR(summary.total_revenue)}
      </p>
    </CardContent>
  </Card>
</div>

    {/* Low Stock Table */}
    <div className="border-1 rounded-md border-gray-300 bg-white p-5 m-5">
      <h2 className="text-xl font-semibold text-center">Stok Kritis</h2>
      <DataTable
        columns={[
          { header: "Nama", accessor: "product_name" },
          { header: "Total Stock", accessor: "total_stock" },
        ]}
        data={lowStockData}
      />
    </div>

    {/* Expiring Soon Table */}
    <div className="border-1 rounded-md border-gray-300 bg-white p-5 m-5">
      <h2 className="text-xl font-semibold text-center">Stok Kedaluarsa</h2>
      <DataTable
        columns={[
          { header: "Nama", accessor: "product_name" },
          { header: "Batch", accessor: "batch_number" },
          { header: "Tanggal Kedaluarsa", accessor: "expiry_date" },
        ]}
        data={expiringSoonData
          .slice(0, 5)
          .map((item) => ({
            ...item,
            expiry_date: formatDateTime(item.expiry_date),
          }))}
      />
    </div>
  </div>
</div>

  );
}
