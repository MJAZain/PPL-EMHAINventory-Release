import React, { useState, useEffect } from "react";
import useSearch from "../../hooks/useSearch";
import { useNavigate } from "react-router-dom";

import { apiClient } from "../../config/api";
import { PlusIcon } from "@heroicons/react/24/solid";

import ActionMenu from "../../components/ActionMenu";
import SearchBar from "../../components/SearchBar";
import DataTable from "../../components/tableCompo";

import SupplierModal from "./SupplierModal";
import Sidebar from "../../components/Sidebar";
import ConfirmDialog from "../../components/ConfirmDialog";
import Toast from "../../components/toast";

function AturSuppliersPage() {
  const [supplierList, setSupplierList] = useState([]);
  const [loading, setLoading] = useState(true);
  const [toast, setToast] = useState(null);

  const [modalOpen, setModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState("add");
  const [modalSupplier, setModalSupplier] = useState(null);

  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [deleteTargetId, setDeleteTargetId] = useState(null);

  const navigate = useNavigate();

  const fetchSuppliers = async () => {
    try {
      const response = await apiClient.get("/suppliers/");
      return response.data;
    } catch (err) {
      setToast({ message: "Data gagal diambil", type: "error" });
      console.error("Fetch error:", err);
      return [];
    }
  };

  const reloadSuppliers = async () => {
    const data = await fetchSuppliers();
    let items = [];

    if (Array.isArray(data)) {
      items = data;
    } else if (Array.isArray(data?.data)) {
      items = data.data;
    } else {
      setToast({ message: "Terdeteksi format tak dikenal", type: "error" });
      items = [];
    }

    setSupplierList(items);
  };

  useEffect(() => {
    reloadSuppliers().finally(() => setLoading(false));
  }, []);

  const openAddModal = () => {
    setModalMode("add");
    setModalSupplier(null);
    setModalOpen(true);
  };

  const openEditModal = async (id) => {
    try {
      const response = await apiClient.get(`/suppliers/${id}`);
      setModalSupplier(response.data.data);
      setModalMode("edit");
      setModalOpen(true);
    } catch (err) {
      setToast({ message: "Supplier gagal diambil", type: "error" });
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
      await apiClient.delete(`/suppliers/${deleteTargetId}`);
      await reloadSuppliers();
      setToast({ message: "Data berhasil dihapus", type: "success" });
    } catch (err) {
      alert("Gagal menghapus data.");
    } finally {
      setIsConfirmOpen(false);
      setDeleteTargetId(null);
    }
  };

  const handleModalSuccess = async () => {
    await reloadSuppliers();
    setToast({
      message: modalMode === "edit" ? "Data berhasil diperbarui" : "Data berhasil ditambahkan",
      type: "success",
    });
    setModalOpen(false);
    setModalSupplier(null);
  };

  const { searchTerm, setSearchTerm, filteredData } = useSearch(supplierList, ["name"]);

  const columns = [
    { header: "Nama Supplier", accessor: "name" },
    { header: "Jenis Supplier", accessor: "type" },
    { header: "Alamat Supplier", accessor: "address" },
    { header: "No Supplier", accessor: "phone" },
    { header: "Email", accessor: "email" },
    { header: "Nama Kontak Person", accessor: "contact_person" },
    { header: "No. Kontak Person", accessor: "contact_number" },
    { header: "Provinsi", accessor: "province" },
    { header: "Kota", accessor: "city" },
    { header: "Status", accessor: "status",
        render: (row) => {
        const value = row.status;
        const isAktif = value === "Aktif";
        return (
          <span
            className={`px-2 py-1 rounded-full text-sm font-medium ${
              isAktif ? "bg-green-100 text-green-800" : "bg-red-100 text-red-800"
            }`}
          >
            {value}
          </span>
        );
      },
     },
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

      <div className="p-5 w-full py-10">
        <h1 className="text-2xl font-bold mb-6">Daftar Supplier</h1>
        <SearchBar value={searchTerm} onChange={setSearchTerm} />

        <div className="border-1 rounded-md border-gray-300 bg-white p-5">
          <div className="flex gap-4 mb-4">
            <button
              onClick={openAddModal}
              className="flex items-center text-blue-700 font-semibold space-x-1 bg-transparent border border-blue-700 py-2 px-4 rounded-md"
            >
              <PlusIcon className="w-4 h-4" />
              <span>Tambah Supplier</span>
            </button>
          </div>

          <div className="max-w-[1121px]">
            <DataTable columns={columns} data={filteredData} showIndex={true} />
          </div>
        </div>
      </div>

      <SupplierModal
        isOpen={modalOpen}
        close={() => setModalOpen(false)}
        onSuccess={handleModalSuccess}
        mode={modalMode}
        supplier={modalSupplier}
      />

      <ConfirmDialog
        isOpen={isConfirmOpen}
        title="Konfirmasi Penghapusan"
        description="Apakah Anda yakin ingin menghapus supplier ini?"
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

export default AturSuppliersPage;
