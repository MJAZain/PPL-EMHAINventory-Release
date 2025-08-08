import React, { useState, useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import { ChevronDown, ChevronUp, Menu, X } from "lucide-react";
import { apiClient } from "../config/api";

export default function Sidebar() {
  const location = useLocation();
  const [isOpen, setIsOpen] = useState(false);
  const [showSetting, setShowSetting] = useState(false);
  const [showPelacakan, setShowPelacakan] = useState(false);
  const [showRiwayat, setShowRiwayat] = useState(false);
  const [showApotek, setShowApotek] = useState(false);
  const [showPengeluaran, setShowPengeluaran] = useState(false);
  const [showStok, setShowStok] = useState(false);

  const isActive = (path) =>
    location.pathname === path || location.pathname.startsWith(path + "/");

  const isAnyChildActiveSetting =
    isActive("/satuan") || isActive("/kategori") || isActive("/storage-locations") || isActive("/brands");

  useEffect(() => {
    if (isAnyChildActiveSetting) setShowSetting(true);
  }, [isAnyChildActiveSetting]);

  const toggleSidebar = () => setIsOpen((prev) => !prev);

  return (
    <div className="relative">
      {!isOpen && (
        <button
          onClick={toggleSidebar}
          className="p-3 m-2 text-gray-700 hover:text-black"
        >
          <Menu size={24} />
        </button>
      )}

      {isOpen && (
        <div className="fixed top-0 left-0 h-full w-64 bg-white shadow-lg z-50 p-4 overflow-y-auto">
          {/* Header */}
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-lg font-bold translate-x-20">EIMS</h2>
            <button
              onClick={toggleSidebar}
              className="text-gray-500 hover:text-black"
            >
              <X size={24} />
            </button>
          </div>

          <Link
            to="/dashboard"
            className={`block px-4 py-2 rounded hover:bg-gray-100 ${
              isActive("/dashboard") ? "bg-gray-200" : ""
            }`}
          >
            Dashboard
          </Link>
          <Link
            to="/sales-analytics"
            className={`block px-4 py-2 rounded hover:bg-gray-100 ${
              isActive("/sales-analytics") ? "bg-gray-200" : ""
            }`}
          >
            Laporan Penjualan
          </Link>
          <Link
            to="/profit"
            className={`block px-4 py-2 rounded hover:bg-gray-100 ${
              isActive("/profit") ? "bg-gray-200" : ""
            }`}
          >
            Laporan Laba/Rugi
          </Link>
          {/* Master Data */}
          <Link
            to="/master-obat"
            className={`block px-4 py-2 rounded hover:bg-gray-100 ${
              isActive("/master-obat") ? "bg-gray-200" : ""
            }`}
          >
            Tabel Master Barang
          </Link>

          {/* Master Setting */}
          <Dropdown
            label="Setting Master Barang"
            isOpen={showSetting}
            onToggle={() => setShowSetting((prev) => !prev)}
            activePaths={["/golongan", "/satuan", "/kategori", "/brands", "/storage-locations"]}
            links={[
              { to: "/golongan", label: "Master Golongan" },
              { to: "/satuan", label: "Master Satuan" },
              { to: "/kategori", label: "Master Kategori" },
              { to: "/brands", label: "Master Brands" },
              { to: "/storage-locations", label: "Master Lokasi Penyimpanan" },
            ]}
          />

          {/* Shift Management */}
          <Link
            to="/shift-resep"
            className={`block px-4 py-2 rounded hover:bg-gray-100 ${
              isActive("/shift-resep") ? "bg-gray-200" : ""
            }`}
          >
            Buka Shift Kasir Resep Dokter
          </Link>
          <Link
            to="/shift-tanpa-resep"
            className={`block px-4 py-2 rounded hover:bg-gray-100 ${
              isActive("/shift-tanpa-resep") ? "bg-gray-200" : ""
            }`}
          >
            Buka Shift Kasir Bebas
          </Link>

          {/* Riwayat */}
          <Dropdown
            label="Riwayat Transaksi"
            isOpen={showRiwayat}
            onToggle={() => setShowRiwayat((prev) => !prev)}
            activePaths={[
              "/shift-riwayat",
              "/regular-riwayat",
              "/pres-riwayat",
              "/riwayat-pbf",
              "/riwayat-non-pbf",
            ]}
            links={[
              { to: "/shift-riwayat", label: "Riwayat Shift" },
              { to: "/regular-riwayat", label: "Penjualan Bebas" },
              { to: "/pres-riwayat", label: "Penjualan Resep" },
              { to: "/riwayat-pbf", label: "Pemesanan PBF" },
              { to: "/riwayat-non-pbf", label: "Pemesanan Non-PBF" },
            ]}
          />

          {/* Pelacakan */}
          <Dropdown
            label="Pelacakan Barang Masuk"
            isOpen={showPelacakan}
            onToggle={() => setShowPelacakan((prev) => !prev)}
            activePaths={["/pbf-detail", "/non-pbf-detail"]}
            links={[
              { to: "/pbf-detail", label: "Barang Masuk PBF" },
              { to: "/non-pbf-detail", label: "Barang Masuk Non-PBF" },
            ]}
          />

          {/* Pengeluaran */}
          <Dropdown
            label="Atur Pengeluaran"
            isOpen={showPengeluaran}
            onToggle={() => setShowPengeluaran((prev) => !prev)}
            activePaths={["/atur-jenis", "/atur-pengeluaran"]}
            links={[
              { to: "/atur-jenis", label: "Jenis Pengeluaran" },
              { to: "/atur-pengeluaran", label: "Pengeluaran" },
            ]}
          />

          {/* Stok */}
          <Dropdown
            label="Manajemen Stok"
            isOpen={showStok}
            onToggle={() => setShowStok((prev) => !prev)}
            activePaths={["/draft", "/koreksi"]}
            links={[
              { to: "/draft", label: "Draft Stock Opname" },
              { to: "/koreksi", label: "Koreksi Stok" },
            ]}
          />

          {/* Apotek */}
          <Dropdown
            label="Atur Data Apotek"
            isOpen={showApotek}
            onToggle={() => setShowApotek((prev) => !prev)}
            activePaths={["/user", "/doctor", "/patients", "/supplier"]}
            links={[
              { to: "/user", label: "Atur Karyawan" },
              { to: "/doctor", label: "Atur Dokter" },
              { to: "/patients", label: "Atur Pasien" },
              { to: "/supplier", label: "Atur Supplier" },
            ]}
          />

          {/* Logout */}
          <button
            className="block px-4 py-2 mt-6 text-red-600 hover:text-red-800"
            onClick={async () => {
              try {
                await apiClient.post("/users/logout");
              } catch (err) {
                console.error("Logout API error:", err);
              } finally {
                localStorage.clear();
                window.location.href = "/";
              }
            }}
          >
            Logout
          </button>
        </div>
      )}
    </div>
  );
}

/**
 * Dropdown helper component
 */
function Dropdown({ label, isOpen, onToggle, activePaths, links }) {
  const location = useLocation();
  const isActive = (path) =>
    location.pathname === path || location.pathname.startsWith(path + "/");
  const isAnyActive = activePaths.some(isActive);

  return (
    <>
      <button
        onClick={onToggle}
        className={`flex items-center justify-between w-full px-4 py-2 mt-2 rounded hover:bg-gray-100 ${
          isAnyActive ? "bg-gray-200" : ""
        }`}
      >
        <span>{label}</span>
        {isOpen || isAnyActive ? <ChevronUp size={18} /> : <ChevronDown size={18} />}
      </button>
      {isOpen && (
        <div className="ml-4">
          {links.map(({ to, label }) => (
            <Link
              key={to}
              to={to}
              className={`block px-4 py-1 mt-1 rounded hover:bg-gray-100 ${
                isActive(to) ? "bg-gray-200" : ""
              }`}
            >
              {label}
            </Link>
          ))}
        </div>
      )}
    </>
  );
}
