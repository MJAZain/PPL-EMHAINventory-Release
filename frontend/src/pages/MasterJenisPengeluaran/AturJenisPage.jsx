import React, { useState, useEffect } from "react";
import useSearch from "../../hooks/useSearch";
import useJenisActions from "./useJenisAction";
import { useNavigate } from "react-router-dom";

import { apiClient } from "../../config/api";
import { PlusIcon } from "@heroicons/react/24/solid";

import ActionMenu from "../../components/ActionMenu";
import SearchBar from "../../components/SearchBar";
import DataTable from "../../components/tableCompo";

import JenisModal from "./JenisModal";
import Sidebar from "../../components/Sidebar";
import ConfirmDialog from "../../components/ConfirmDialog";
import Toast from "../../components/toast";

function AturJenisPage() {
  const [jenisList, setJenisList] = useState([]);
  const [loading, setLoading] = useState(true);
  const [toast, setToast] = useState(null);

  const [modalOpen, setModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState("add");
  const [modalJenis, setModalJenis] = useState(null);

  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [deleteTargetId, setDeleteTargetId] = useState(null);

  const navigate = useNavigate();
  const { getJenisById, deleteJenis } = useJenisActions();

  const fetchJenis = async () => {
    try {
      const response = await apiClient.get("/expense-types/");
      return response.data;
    } catch (err) {
      setToast({ message: "Data gagal diambil", type: "error" });
      console.error("Fetch error:", err);
      return [];
    }
  };

  const reloadJenis = async () => {
    const data = await fetchJenis();
    const items = Array.isArray(data) ? data : data?.data || [];
    setJenisList(items);
  };

  useEffect(() => {
    reloadJenis().finally(() => setLoading(false));
  }, []);

  const openAddModal = () => {
    setModalMode("add");
    setModalJenis(null);
    setModalOpen(true);
  };

  const openEditModal = async (id) => {
    try {
      const jenis = await getJenisById(id);
      setModalJenis(jenis);
      setModalMode("edit");
      setModalOpen(true);
    } catch (err) {
      setToast({ message: "Jenis gagal diambil", type: "error" });
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
      await deleteJenis(deleteTargetId);
      await reloadJenis();
      setToast({ message: "Jenis berhasil dihapus", type: "success" });
    } catch (err) {
      alert("Gagal menghapus jenis.");
    } finally {
      setIsConfirmOpen(false);
      setDeleteTargetId(null);
    }
  };

  const handleModalSuccess = async () => {
    await reloadJenis();
    setToast({
      message:
        modalMode === "edit"
          ? "Data berhasil diperbarui"
          : "Data berhasil ditambahkan",
      type: "success",
    });
    setModalOpen(false);
    setModalJenis(null);
  };

  const { searchTerm, setSearchTerm, filteredData } = useSearch(jenisList, [
    "name",
  ]);

  const columns = [
    { header: "Nama Jenis Pengeluaran", accessor: "name" },
    {
      header: "Pilih Aksi",
      accessor: "actions",
      isAction: true,
      render: (item) => (
        <ActionMenu
          actions={[
            { label: "Edit", onClick: () => openEditModal(item.id) },
            { label: "Hapus", onClick: () => handleDeleteRequest(item.id) },
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
        <h1 className="text-2xl font-bold mb-6">Daftar Jenis Pengeluaran</h1>
        <SearchBar value={searchTerm} onChange={setSearchTerm} />

        <div className="border-1 rounded-md border-gray-300 bg-white p-5">
          <div className="flex gap-4 mb-4">
            <button
              onClick={openAddModal}
              className="flex items-center text-blue-700 font-semibold space-x-1 bg-transparent border border-blue-700 py-2 px-4 rounded-md"
            >
              <PlusIcon className="w-4 h-4" />
              <span>Tambah Jenis</span>
            </button>
          </div>

          <div className="max-w-[1121px]">
            <DataTable columns={columns} data={filteredData} showIndex={true} />
          </div>
        </div>
      </div>

      <JenisModal
        isOpen={modalOpen}
        close={() => setModalOpen(false)}
        onSuccess={handleModalSuccess}
        mode={modalMode}
        jenis={modalJenis}
      />

      <ConfirmDialog
        isOpen={isConfirmOpen}
        title="Konfirmasi Penghapusan"
        description="Apakah Anda yakin ingin menghapus jenis ini?"
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

export default AturJenisPage;
