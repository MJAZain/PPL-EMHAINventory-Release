import React, { useState, useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import { ChevronDown, ChevronUp, Menu, X } from "lucide-react";
import { apiClient } from "../config/api";

export default function Sidebar() {
  const location = useLocation();
  const [isOpen, setIsOpen] = useState(false);
  const [showPelacakan, setShowPelacakan] = useState(false);
  const [showSetting, setShowSetting] = useState(false);
  const [showRiwayat, setShowRiwayat] = useState(false);
  const [showApotek, setShowApotek] = useState(false);

  const isActive = (path) =>
    location.pathname === path || location.pathname.startsWith(path + "/");

  const isAnyChildActiveMasterBarang =
    isActive("/satuan") ||
    isActive("/kategori") ||
    isActive("/storage-locations");
  isActive("/brands");
  useEffect(() => {
    if (isAnyChildActiveMasterBarang) {
      setShowSetting(true);
    }
  }, [isAnyChildActiveMasterBarang]);

  const toggleSidebar = () => setIsOpen((prev) => !prev);
  const togglePelacakanDropdown = () => setShowPelacakan((prev) => !prev);
  const toggleSettingDropdown = () => setShowSetting((prev) => !prev);
  const toggleRiwayatDropdown = () => setShowRiwayat((prev) => !prev);
  const toggleApotekDropdown = () => setShowApotek((prev) => !prev);

  return (
    <div className="relative">
      {/* Hamburger Button */}
      {!isOpen && (
        <button
          onClick={toggleSidebar}
          className="p-3 m-2 text-gray-700 hover:text-black"
        >
          <Menu size={24} />
        </button>
      )}

      {/* Sidebar */}
      {isOpen && (
        <div className="fixed top-0 left-0 h-full w-64 bg-white shadow-lg z-50 p-4">
          {/* Close Button */}
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-lg font-bold translate-x-20">EIMS</h2>
            <button
              onClick={toggleSidebar}
              className="text-gray-500 hover:text-black"
            >
              <X size={24} />
            </button>
          </div>

          {/* Dashboard */}
          <Link
            to="/dashboard"
            className={`block px-4 py-2 rounded hover:bg-gray-100 ${
              isActive("/dashboard") ? "bg-gray-200" : ""
            }`}
          >
            Tabel Barang
          </Link>
          {/*
          <Link
            to="/shift-resep"
            className={`block px-4 py-2 rounded hover:bg-gray-100 ${
              isActive("/shift-resep") ? "bg-gray-200" : ""
            }`}
          >
            Buka Shift Kasir dengan Resep
          </Link>
          <Link
            to="/shift-tanpa-resep"
            className={`block px-4 py-2 rounded hover:bg-gray-100 ${
              isActive("/shift-tanpa-resep") ? "bg-gray-200" : ""
            }`}
          >
            Buka Shift Kasir
          </Link>
          */}

          {/* Master Setting Dropdown */}
          <button
            onClick={toggleSettingDropdown}
            className={`flex items-center justify-between w-full px-4 py-2 mt-2 rounded hover:bg-gray-100 ${
              isActive("/satuan") ||
              isActive("/kategori") ||
              isActive("/storage-locations")
                ? "bg-gray-200"
                : ""
            }`}
          >
            <span>Setting Master Barang</span>
            {showSetting ||
            isActive("/satuan") ||
            isActive("/kategori") ||
            isActive("/storage-locations") ? (
              <ChevronUp size={18} />
            ) : (
              <ChevronDown size={18} />
            )}
          </button>

          {showSetting && (
            <div className="ml-4">
              <Link
                to="/golongan"
                className={`block px-4 py-1 hover:bg-gray-100 rounded ${
                  isActive("/golongan") ? "bg-gray-200" : ""
                }`}
              >
                Master Golongan
              </Link>
              <Link
                to="/satuan"
                className={`block px-4 py-1 hover:bg-gray-100 rounded ${
                  isActive("/satuan") ? "bg-gray-200" : ""
                }`}
              >
                Master Satuan
              </Link>
              <Link
                to="/kategori"
                className={`block px-4 py-1 hover:bg-gray-100 rounded ${
                  isActive("/kategori") ? "bg-gray-200" : ""
                }`}
              >
                Master Kategori
              </Link>
              <Link
                to="/brands"
                className={`block px-4 py-1 hover:bg-gray-100 rounded ${
                  isActive("/brands") ? "bg-gray-200" : ""
                }`}
              >
                Master Brands
              </Link>
              <Link
                to="/storage-locations"
                className={`block px-4 py-1 hover:bg-gray-100 rounded ${
                  isActive("/storage-locations") ? "bg-gray-200" : ""
                }`}
              >
                Master Lokasi Penyimpanan
              </Link>
            </div>
          )}

          {/* Pelacakan Barang Dropdown */}
          <button
            onClick={togglePelacakanDropdown}
            className="flex items-center justify-between w-full px-4 py-2 mt-2 rounded hover:bg-gray-100"
          >
            <span>Pelacakan Barang Masuk</span>
            {showPelacakan ? (
              <ChevronUp size={18} />
            ) : (
              <ChevronDown size={18} />
            )}
          </button>

          {showPelacakan && (
            <div className="ml-4">
              <Link
                to="/pbf-detail"
                className={`block px-4 py-1 hover:bg-gray-100 rounded ${
                  isActive("/pbf-detail") ? "bg-gray-200" : ""
                }`}
              >
                Barang Masuk PBF
              </Link>
              <Link
                to="/non-pbf-detail"
                className={`block px-4 py-1 hover:bg-gray-100 rounded ${
                  isActive("/non-pbf-detail") ? "bg-gray-200" : ""
                }`}
              >
                Barang Masuk Non-PBF
              </Link>
            </div>
          )}

          {/* Riwayat Barang Dropdown */}
          <button
            onClick={toggleRiwayatDropdown}
            className="flex items-center justify-between w-full px-4 py-2 mt-2 rounded hover:bg-gray-100"
          >
            <span>Riwayat Transaksi</span>
            {showRiwayat ? <ChevronUp size={18} /> : <ChevronDown size={18} />}
          </button>

          {showRiwayat && (
            <div className="ml-4">
              <Link
                to="/riwayat-pbf"
                className={`block px-4 py-1 hover:bg-gray-100 rounded ${
                  isActive("/riwayat-pbf") ? "bg-gray-200" : ""
                }`}
              >
                Riwayat Pemesanan PBF
              </Link>
              <Link
                to="/riwayat-non-pbf"
                className={`block px-4 py-1 hover:bg-gray-100 rounded ${
                  isActive("/riwayat-non-pbf") ? "bg-gray-200" : ""
                }`}
              >
                Riwayat Pemesanan Non-PBF
              </Link>
            {/*
            <Link
                to="/laporan-terjual"
                className={`block px-4 py-1 rounded hover:bg-gray-100 ${
                  isActive("/laporan-terjual") ? "bg-gray-200" : ""
                }`}
              >
                Riwayat Laporan Kasir
              </Link>
              */}
            </div>
          )}

          {/* Pelacakan Barang Dropdown */}
          <button
            onClick={toggleApotekDropdown}
            className="flex items-center justify-between w-full px-4 py-2 mt-2 rounded hover:bg-gray-100"
          >
            <span>Atur Data Apotek</span>
            {showApotek ? (
              <ChevronUp size={18} />
            ) : (
              <ChevronDown size={18} />
            )}
          </button>

          {showApotek && (
            <div className="ml-4">
              <Link
                to="/user"
                className={`block px-4 py-2 mt-2 rounded hover:bg-gray-100 ${
                  isActive("/user") ? "bg-gray-200" : ""
                }`}
              >
                Atur Karyawan
              </Link>

              <Link
                to="/doctor"
                className={`block px-4 py-2 mt-2 rounded hover:bg-gray-100 ${
                  isActive("/doctor") ? "bg-gray-200" : ""
                }`}
              >
                Atur Dokter
              </Link>

              <Link
                to="/patients"
                className={`block px-4 py-2 mt-2 rounded hover:bg-gray-100 ${
                  isActive("/patients") ? "bg-gray-200" : ""
                }`}
              >
                Atur Pasien
              </Link>

              <Link
                to="/supplier"
                className={`block px-4 py-2 mt-2 rounded hover:bg-gray-100 ${
                  isActive("/Supplier") ? "bg-gray-200" : ""
                }`}
              >
                Atur Supplier
              </Link>
            </div>
          )}
{/* <Link
                to="/stock-opname"
                className={`block px-4 py-1 hover:bg-gray-100 rounded ${
                  isActive("/stock-opname") ? "bg-gray-200" : ""
                }`}
              >
                Stock Opname
              </Link>
                           <Link
                to="/koreksi"
                className={`block px-4 py-1 hover:bg-gray-100 rounded ${
                  isActive("/koreksi") ? "bg-gray-200" : ""
                }`}
              >
                Koreksi Stok
              </Link> */}
              

          {/* Logout */}
          <button
            className="block px-4 py-2 mt-6 text-red-600 hover:text-red-800"
            onClick={async () => {
              try {
                await apiClient.post("/users/logout");
              } catch (err) {
                console.error("Logout API error:", err);
              } finally {
                localStorage.removeItem("token");
                localStorage.removeItem("user");
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
