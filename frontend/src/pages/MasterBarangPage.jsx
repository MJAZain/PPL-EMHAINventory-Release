import React, { useState, useEffect } from "react";
import useSearch from "../hooks/useSearch";
import useProductActions from "../hooks/useProductsAction";
import { useNavigate } from "react-router-dom";

import { apiClient } from "../config/api";
import { PlusIcon } from "@heroicons/react/24/solid";

import ActionMenu from "../components/ActionMenu";
import SearchBar from "../components/SearchBar";
import DataTable from "../components/tableCompo";
import Button from "../components/buttonComp";
import BarangModal from "../components/modal/BarangModal";
import Sidebar from "../components/Sidebar";
import ConfirmDialog from "../components/ConfirmDialog";
import Toast from "../components/toast";

function MasterBarangPage() {
  const [barangList, setBarangList] = useState([]);
  const [loading, setLoading] = useState(true);
  const [toast, setToast] = useState(null);

  const [modalOpen, setModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState("add");
  const [modalProduct, setModalProduct] = useState(null);

  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [deleteTargetId, setDeleteTargetId] = useState(null);

  const navigate = useNavigate();
  const { getProductById, deleteProduct } = useProductActions();

  const fetchBarang = async () => {
    try {
      const response = await apiClient.get("/products/");
      return response.data;
    } catch (err) {
      setToast({
        message: "Barang gagal diambil",
        type: "error",
      });
      console.error("Fetch error:", err);
      return [];
    }
  };

  const reloadBarang = async () => {
    const data = await fetchBarang();
    let items = [];

    if (Array.isArray(data)) {
      items = data;
    } else if (Array.isArray(data?.data)) {
      items = data.data;
    } else {
      setToast({ message: "Terdeteksi format tak dikenal", type: "error" });
      items = [];
    }

    setBarangList(items);
  };

  useEffect(() => {
    reloadBarang().finally(() => setLoading(false));
  }, []);

  const openAddModal = () => {
    setModalMode("add");
    setModalProduct(null);
    setModalOpen(true);
  };

  const openEditModal = async (id) => {
    try {
      const product = await getProductById(id);
      setModalProduct(product);
      setModalMode("edit");
      setModalOpen(true);
    } catch (err) {
      setToast({
        message: "Barang gagal diambil",
        type: "error",
      });
    }
  };

  const handleDeleteRequest = (id) => {
    setDeleteTargetId(id);
    setIsConfirmOpen(true);
  };

  const handleCancelDelete = () => {
    setIsConfirmOpen(false);
    setDeleteTargetId(null);
  };

  const handleConfirmDelete = async () => {
    try {
      await deleteProduct(deleteTargetId);
      await reloadBarang();
      setToast({ message: "Barang berhasil dihapus", type: "success" });
    } catch (err) {
      alert("Gagal menghapus data.");
    } finally {
      setIsConfirmOpen(false);
      setDeleteTargetId(null);
    }
  };

  const handleModalSuccess = async () => {
    await reloadBarang();
    setToast({
      message:
        modalMode === "edit"
          ? "Barang berhasil diperbarui"
          : "Barang berhasil ditambahkan",
      type: "success",
    });
    setModalOpen(false);
    setModalProduct(null);
  };

  const { searchTerm, setSearchTerm, filteredData } = useSearch(barangList, [
    "name",
    "code",
    "barcode",
    "category.name",
    "brand.name",
    "unit.name",
    "storage_location.name",
    "drug_category.name",
  ]);

  const columns = [
    { header: "Nama", accessor: "name" },
    { header: "SKU", accessor: "code" },
    { header: "Barcode", accessor: "barcode" },
    { header: "Golongan Obat", accessor: (item) => item.drug_category?.name || "Tidak ada Data" },
    { header: "Kategori Obat", accessor: (item) => item.category?.name || "Tidak ada Data" },
    { header: "Satuan", accessor: (item) => item.unit?.name || "Tidak ada Data" },
    { header: "Harga Jual", accessor: "selling_price" },
    { header: "Lokasi", accessor: (item) => item.storage_location?.name || "Tidak ada Data" },
    { header: "Merk", accessor: (item) => item.brand?.name || "Tidak ada Data" },
    { header: "Stok Minimal", accessor: "min_stock" || "Tidak ada Data"},
    { header: "Dosis", accessor: "dosage_description" || "Tidak ada Data"},
    { header: "Komposisi", accessor: "composition_description" || "Tidak ada Data"},
    {
      header: "Pilih Aksi",
      accessor: "actions",
      isAction: true,
      render: (item) => (
        <ActionMenu
          actions={[
            { label: "Edit", onClick: () => openEditModal(item.id) },
            { label: "Delete", onClick: () => handleDeleteRequest(item.id) },
          ]}
        />
      ),
    },
  ];

  if (loading) return <div className="text-center mt-6">Loading...</div>;

  return (
    <div className="flex">
      <div className="bg-white min-h-screen">
        <Sidebar />
      </div>

      <div className="p-5">
        <h1 className="text-2xl font-bold mb-6">Daftar Barang</h1>
        <SearchBar value={searchTerm} onChange={setSearchTerm} />

        <div className="border-1 rounded-md border-gray-300 bg-white p-5">
          <div className="flex gap-4 mb-4">
            <button
              onClick={openAddModal}
              className="flex items-center text-blue-700 font-semibold space-x-1 bg-transparent border border-blue-700 py-2 px-4 rounded-md"
            >
              <PlusIcon className="w-4 h-4" />
              <span>Tambah Barang</span>
            </button>
          </div>

          <div className="max-w-[1121px]">
            <DataTable
              columns={columns}
              data={filteredData}
              showIndex={true}
            />
          </div>
        </div>
      </div>

      <BarangModal
        isOpen={modalOpen}
        close={() => setModalOpen(false)}
        onSuccess={handleModalSuccess}
        mode={modalMode}
        product={modalProduct}
      />

      <ConfirmDialog
        isOpen={isConfirmOpen}
        title="Konfirmasi Penghapusan"
        description="Apakah Anda yakin ingin menghapus obat ini?"
        onCancel={handleCancelDelete}
        onConfirm={handleConfirmDelete}
      />

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

export default MasterBarangPage;
