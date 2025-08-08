import React, { useState, useEffect } from "react";
import {
  VictoryChart,
  VictoryLine,
  VictoryAxis,
  VictoryLegend,
  VictoryTheme,
  VictoryTooltip,
  VictoryScatter,
  VictoryBar,
  VictoryLabel
} from "victory";
import DataTable from "../../components/tableCompo";
import { apiClient } from "../../config/api";
import Sidebar from "../../components/Sidebar";

const colorPalette = [
  "#8884d8",
  "#82ca9d",
  "#ffc658",
  "#ff7300",
  "#a4de6c",
  "#d0ed57",
];

const groupByProduct = (rawData, timeRange) => {
  const grouped = {};

  rawData.forEach((item) => {
    let date;
    
    if (timeRange === "monthly") {
      const [year, month] = item.date.split('-');
      date = new Date(year, month - 1);
    } else if (timeRange === "yearly") {
      date = new Date(item.date, 0);
    } else {
      date = new Date(item.date);
    }

    if (!grouped[item.period]) {
      grouped[item.period] = [];
    }

    grouped[item.period].push({
      x: date,
      y: item.total,
    });
  });

  Object.values(grouped).forEach(data => {
    data.sort((a, b) => a.x - b.x);
  });

  return Object.entries(grouped).map(([name, data]) => ({
    name,
    data,
  }));
};

const AnalyticsDashboard = () => {
  const [timeRange, setTimeRange] = useState("weekly");
  const [lineData, setLineData] = useState([]);
  const [topProducts, setTopProducts] = useState([]);
  const [leastProducts, setLeastProducts] = useState([]);
  const [isLoading, setIsLoading] = useState(false);

  const fetchAllData = async () => {
    setIsLoading(true);
    try {
      const [lineRes, topRes, leastRes] = await Promise.all([
        apiClient.post("/sales/analytics/line-chart", { time_range: timeRange }),
        apiClient.post("/sales/analytics/top-products", { time_range: timeRange }),
        apiClient.post("/sales/analytics/least-products", { time_range: timeRange }),
      ]);

      const rawLineData = lineRes.data?.data || [];
      const groupedLineData = groupByProduct(rawLineData, timeRange);

      // If only one data point exists, add a previous point for visualization
      const processedData = groupedLineData.map(series => {
        if (series.data.length === 1) {
          const singlePoint = series.data[0];
          let prevDate;
          
          if (timeRange === "monthly") {
            prevDate = new Date(singlePoint.x);
            prevDate.setMonth(prevDate.getMonth() - 1);
          } else if (timeRange === "yearly") {
            prevDate = new Date(singlePoint.x);
            prevDate.setFullYear(prevDate.getFullYear() - 1);
          } else {
            prevDate = new Date(singlePoint.x);
            prevDate.setDate(prevDate.getDate() - 7);
          }

          return {
            ...series,
            data: [
              { x: prevDate, y: 0 },
              ...series.data
            ]
          };
        }
        return series;
      });

      setLineData(processedData);
      setTopProducts(topRes.data?.data || []);
      setLeastProducts(leastRes.data?.data || []);
    } catch (err) {
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchAllData();
  }, [timeRange]);

  const renderChart = () => {
    if (isLoading) {
      return (
        <div className="flex justify-center items-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
        </div>
      );
    }

    if (lineData.length === 0) {
      return (
        <div className="flex justify-center items-center h-64">
          <p className="text-gray-500">No data available</p>
        </div>
      );
    }

    const hasEnoughData = lineData.some(series => series.data.length > 1);
    const singlePointSeries = lineData.find(series => series.data.length === 1);

    return (
      <VictoryChart
        theme={VictoryTheme.material}
        scale={{ x: "time" }}
        width={window.innerWidth * 0.9}
        height={400}
        domainPadding={{ x: 20, y: 20 }}
      >
        <VictoryAxis
          tickFormat={(x) => {
            const date = new Date(x);
            if (timeRange === "monthly") {
              return date.toLocaleDateString("id-ID", { month: "short", year: "numeric" });
            } else if (timeRange === "yearly") {
              return date.getFullYear();
            } else {
              return date.toLocaleDateString("id-ID", { day: "numeric", month: "short" });
            }
          }}
        />
        <VictoryAxis
          dependentAxis
          tickFormat={(y) => `Rp${y.toLocaleString("id-ID")}`}
        />

        {hasEnoughData ? (
          lineData.map((series, idx) => (
            <VictoryLine
              key={series.name}
              data={series.data}
              style={{
                data: { 
                  stroke: colorPalette[idx % colorPalette.length],
                  strokeWidth: 3
                }
              }}
              labels={({ datum }) => `Rp${datum.y.toLocaleString("id-ID")}`}
              labelComponent={<VictoryTooltip />}
            />
          ))
        ) : singlePointSeries && (
          <>
            <VictoryBar
              data={singlePointSeries.data}
              style={{
                data: { fill: colorPalette[0] }
              }}
              labels={({ datum }) => `Rp${datum.y.toLocaleString("id-ID")}`}
              labelComponent={<VictoryTooltip />}
            />
            <VictoryLabel
              x={window.innerWidth * 0.45}
              y={30}
              text={`Only one data point available (${new Date(singlePointSeries.data[0].x).toLocaleDateString("id-ID", { 
                month: "long", 
                year: "numeric" 
              })})`}
              textAnchor="middle"
              style={{ fontSize: 12 }}
            />
          </>
        )}

        <VictoryLegend
          x={125}
          y={10}
          orientation="horizontal"
          gutter={20}
          data={lineData.map((series, idx) => ({
            name: series.name,
            symbol: { fill: colorPalette[idx % colorPalette.length] },
          }))}
        />
      </VictoryChart>
    );
  };

  return (
    <div className="flex">
      <div className="bg-white min-h-screen">
        <Sidebar />
      </div>

      <div className="p-5 w-full py-10">
        <h1 className="text-2xl font-bold mb-4">Laporan Penjualan</h1>

        <div className="mb-4">
          <label className="block text-sm font-medium">Rentang Waktu</label>
          <select
            value={timeRange}
            onChange={(e) => setTimeRange(e.target.value)}
            className="border rounded px-2 py-1"
          >
            <option value="weekly">Mingguan</option>
            <option value="monthly">Bulanan</option>
            <option value="yearly">Tahunan</option>
          </select>
        </div>

        {/* Line Chart */}
        <div className="bg-white p-4 rounded shadow m-5">
          <h2 className="text-lg font-semibold mb-2 text-center">Trend Penjualan</h2>
          <div className="w-full">
            {renderChart()}
          </div>
        </div>

        {/* Top Products */}
        <div className="bg-white p-4 rounded shadow m-5">
          <h2 className="text-lg font-semibold mb-2 text-center">Obat Terlaris</h2>
          <DataTable
            columns={[
              { header: "Nama Obat", accessor: "product_name" },
              { header: "Kuantitas", accessor: "total_qty" },
              { header: "Total Pemasukan", accessor: "total_amount" },
              { header: "Total Penjualan", accessor: "sales_count" },
            ]}
            data={topProducts}
          />
        </div>

        {/* Least Products */}
        <div className="bg-white p-4 rounded shadow m-5">
          <h2 className="text-lg font-semibold mb-2 text-center">Obat Stok Macet</h2>
          <DataTable
            columns={[
              { header: "Nama Obat", accessor: "product_name" },
              { header: "Kuantitas", accessor: "total_qty" },
              { header: "Total Pemasukan", accessor: "total_amount" },
              { header: "Total Penjualan", accessor: "sales_count" },
            ]}
            data={leastProducts}
          />
        </div>
      </div>
    </div>
  );
};

export default AnalyticsDashboard;