import { useState } from "react";
import { Link, useLocation } from "react-router-dom";
import { ChevronDown, ChevronUp, Menu, X } from "lucide-react";
import { apiClient } from "../config/api";

export default function Sidebar() {
  const location = useLocation();
  const [isOpen, setIsOpen] = useState(false);
  const [showPelacakan, setShowPelacakan] = useState(false);
  const [showSetting, setShowSetting] = useState(false);
  const [showLaporan, setShowLaporan] = useState(false);

  const toggleSidebar = () => setIsOpen((prev) => !prev);
  const togglePelacakanDropdown = () => setShowPelacakan((prev) => !prev);
  const toggleSettingDropdown = () => setShowSetting((prev) => !prev);
  const toggleLaporanDropdown = () => setShowLaporan((prev) => !prev);
  const isActive = (path) => location.pathname === path;

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

          {/* Master Setting Dropdown */}
          <button
            onClick={toggleSettingDropdown}
            className="flex items-center justify-between w-full px-4 py-2 mt-2 rounded hover:bg-gray-100"
          >
            <span>Setting Master Barang</span>
            {showSetting ? <ChevronUp size={18} /> : <ChevronDown size={18} />}
          </button>

          {showSetting && (
            <div className="ml-4">
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
            </div>
          )}

          {/* Pelacakan Barang Dropdown */}
          <button
            onClick={togglePelacakanDropdown}
            className="flex items-center justify-between w-full px-4 py-2 mt-2 rounded hover:bg-gray-100"
          >
            <span>Pelacakan Barang</span>
            {showPelacakan ? <ChevronUp size={18} /> : <ChevronDown size={18} />}
          </button>

          {showPelacakan && (
            <div className="ml-4">
              <Link
                to="/barang-masuk"
                className={`block px-4 py-1 hover:bg-gray-100 rounded ${
                  isActive("/barang-masuk") ? "bg-gray-200" : ""
                }`}
              >
                Barang Masuk
              </Link>
              <Link
                to="/barang-keluar"
                className={`block px-4 py-1 hover:bg-gray-100 rounded ${
                  isActive("/barang-keluar") ? "bg-gray-200" : ""
                }`}
              >
                Barang Keluar
              </Link>
            </div>
          )}

          {/* Laporan Barang Dropdown */}
          <button
            onClick={toggleLaporanDropdown}
            className="flex items-center justify-between w-full px-4 py-2 mt-2 rounded hover:bg-gray-100"
          >
            <span>Laporan Barang</span>
            {showLaporan ? <ChevronUp size={18} /> : <ChevronDown size={18} />}
          </button>

          {showLaporan && (
            <div className="ml-4">
              <Link
                to="/laporan-masuk"
                className={`block px-4 py-1 hover:bg-gray-100 rounded ${
                  isActive("/laporan-masuk") ? "bg-gray-200" : ""
                }`}
              >
                Laporan Barang Masuk
              </Link>
              <Link
                to="/laporan-terjual"
                className={`block px-4 py-1 rounded hover:bg-gray-100 ${
                  isActive("/laporan-terjual") ? "bg-gray-200" : ""
                }`}
              >
                Laporan Barang Terjual
              </Link>
            </div>
          )}

          {/* Laporan Keuangan */}
          

          {/*  */}
          <Link
            to="/user"
            className={`block px-4 py-2 mt-2 rounded hover:bg-gray-100 ${
              isActive("/user") ? "bg-gray-200" : ""
            }`}
          >
            User Management
          </Link>

          {/* Logout */}
          <button
            className="block px-4 py-2 mt-6 text-red-600 hover:text-red-800"
            onClick={async () => {
              try {
                await apiClient.post("/users/logout"); // 🔁 Call the backend logout API
              } catch (err) {
                console.error("Logout API error:", err); // Optional logging
              } finally {
                localStorage.removeItem("token"); // 🚪 Clear token
                window.location.href = "/";       // ⏩ Redirect to login/home
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
