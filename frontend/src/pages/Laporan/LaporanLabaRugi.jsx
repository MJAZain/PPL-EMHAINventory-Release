import React, { useEffect, useState } from "react";
import {
  VictoryPie,
  VictoryBar,
  VictoryChart,
  VictoryAxis,
  VictoryTheme,
} from "victory";
import { apiClient } from "../../config/api";
import { Card, CardContent } from "../../components/cardComp";
import Sidebar from "../../components/Sidebar";

const ProfitAnalysisPage = () => {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [chartWidth, setChartWidth] = useState(window.innerWidth * 0.85);
  const [selectedChart, setSelectedChart] = useState("Pendapatan");

  useEffect(() => {
    const handleResize = () => setChartWidth(window.innerWidth * 0.85);
    handleResize(); // Initial set
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  useEffect(() => {
    const fetchAnalysisData = async () => {
      try {
        const response = await apiClient.get("/analysis/");
        setData(response.data.data);
      } catch (err) {
        console.error(err);
        setError("Failed to load analysis data.");
      } finally {
        setLoading(false);
      }
    };

    fetchAnalysisData();
  }, []);

  const formatToIDR = (value) => {
    if (value >= 1_000_000_000) return `Rp ${(value / 1_000_000_000).toFixed(1)}B`;
    if (value >= 1_000_000) return `Rp ${(value / 1_000_000).toFixed(1)}M`;
    if (value >= 1_000) return `Rp ${(value / 1_000).toFixed(1)}K`;
    return `Rp ${value}`;
  };

  if (loading) return <div className="p-6 text-gray-500">Loading analysis...</div>;
  if (error) return <div className="p-6 text-red-500">{error}</div>;
  if (!data) return null;

  const {
    total_gross_revenue,
    total_expense,
    net_profit,
    profit_loss_compare,
    top_selling_products,
    pie_chart,
    bar_chart_revenue,
    bar_chart_expense,
  } = data;

  const topProduct = top_selling_products?.[0];
  const currentBarChart =
    selectedChart === "Pendapatan" ? bar_chart_revenue : bar_chart_expense;
  const barColor = selectedChart === "Pendapatan" ? "#60a5fa" : "#f87171";

  return (
    <div className="flex">
      <div className="bg-white min-h-screen">
        <Sidebar />
      </div>

      <div className="p-5 w-full py-10">
        <h1 className="text-2xl font-bold mb-4">Analisis Laba/Rugi</h1>

        {/* Summary Cards */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <Card>
            <CardContent className="p-6 text-center">
              <p className="text-sm text-muted-foreground">Total Pendapatan Kotor</p>
              <p className="text-xl font-bold text-green-600">
                Rp {total_gross_revenue.toLocaleString()}
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-6 text-center">
              <p className="text-sm text-muted-foreground">Total Pengeluaran</p>
              <p className="text-xl font-bold text-red-600">
                Rp {total_expense.toLocaleString()}
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-6 text-center">
              <p className="text-sm text-muted-foreground">Net Profit</p>
              <p className="text-xl font-bold text-blue-600">
                Rp {net_profit.toLocaleString()}
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Profit Status & Top Product */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-5">
          <Card>
            <CardContent className="p-6 text-center m-5">
              <p className="text-sm text-muted-foreground">Status</p>
              <p className="font-semibold">{profit_loss_compare.status}</p>
              <p className="text-sm text-muted-foreground mt-2">Pertumbuhan</p>
              <p className="font-bold">
                {profit_loss_compare.percentage}%
              </p>
              <p className="text-sm mt-2">{profit_loss_compare.message}</p>
            </CardContent>
          </Card>

          {topProduct && (
            <Card>
              <CardContent className="p-6 text-center m-5">
                <p className="text-sm text-muted-foreground">Top Obat</p>
                <p className="font-semibold">{topProduct.product_name}</p>
                <p className="text-sm text-muted-foreground mt-2">Total Pendapatan Obat</p>
                <p className="font-bold text-purple-600">
                  Rp {topProduct.total_revenue.toLocaleString()}
                </p>
              </CardContent>
            </Card>
          )}
        </div>

        {/* Pie Chart */}
        <div className="m-5">
          <Card>
            <CardContent className="p-4">
              <h2 className="text-lg font-semibold mb-4 text-center">
                Pendapatan vs Pengeluaran
              </h2>
{pie_chart?.labels?.length > 0 && pie_chart?.values?.length > 0 ? (
  <VictoryPie
    colorScale={["#34d399", "#f87171"]}
    data={pie_chart.labels.map((label, idx) => ({
      x: label,
      y: pie_chart.values[idx],
    }))}
    labels={({ datum }) => `${datum.x}: ${datum.y.toLocaleString()}`}
    style={{ labels: { fontSize: 15, fill: "#4b5563" } }}
    width={chartWidth}
    height={350}
  />
) : (
  <p className="text-center text-gray-500 py-10">
    Belum ada data untuk bulan ini
  </p>
)}

            </CardContent>
          </Card>
        </div>

        {/* Combined Bar Chart with Picker */}
        <div className="m-5">
          <Card>
            <CardContent className="p-4">
              <div className="flex justify-between items-center mb-4">
                <h2 className="text-lg font-semibold text-center w-full">
                  Rincian {selectedChart}
                </h2>
                <select
                  value={selectedChart}
                  onChange={(e) => setSelectedChart(e.target.value)}
                  className="ml-auto border border-gray-300 rounded px-2 py-1 text-sm"
                >
                  <option value="Pendapatan">Pendapatan</option>
                  <option value="Pengeluaran">Pengeluaran</option>
                </select>
              </div>
{currentBarChart?.labels?.length > 0 && currentBarChart?.values?.length > 0 ? (
  <VictoryChart
    domainPadding={selectedChart === "Pendapatan" ? { x: 250 } : 20}
    theme={VictoryTheme.material}
    width={chartWidth}
    height={300}
    padding={{ top: 20, bottom: 60, left: 80, right: 40 }}
  >
    <VictoryAxis style={{ tickLabels: { fontSize: 15, padding: 10 } }} />
    <VictoryAxis
      dependentAxis
      tickFormat={(tick) => formatToIDR(tick)}
    />
    <VictoryBar
      style={{ data: { fill: barColor } }}
      data={currentBarChart.labels.map((label, idx) => ({
        x: label,
        y: currentBarChart.values[idx],
      }))}
    />
  </VictoryChart>
) : (
  <p className="text-center text-gray-500 py-10">
    Belum ada data untuk bulan ini
  </p>
)}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
};

export default ProfitAnalysisPage;
