import React, { useState, useEffect } from "react";
import useSearch from "../../hooks/useSearch";
import useSatuanActions from "./useSatuanActions";
import { useNavigate } from "react-router-dom";

import { apiClient } from "../../config/api";
import { PlusIcon } from "@heroicons/react/24/solid";

import ActionMenu from "../../components/ActionMenu";
import SearchBar from "../../components/SearchBar";
import DataTable from "../../components/tableCompo";

import SatuanModal from "./SatuanModal";
import Sidebar from "../../components/Sidebar";
import ConfirmDialog from "../../components/ConfirmDialog";
import Toast from "../../components/toast";

function AturSatuanPage() {
  const [satuanList, setSatuanList] = useState([]);
  const [loading, setLoading] = useState(true);
  const [toast, setToast] = useState(null);

  const [modalOpen, setModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState("add");
  const [modalSatuan, setModalSatuan] = useState(null);

  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [deleteTargetId, setDeleteTargetId] = useState(null);

  const navigate = useNavigate();
  const { getSatuanById, deleteSatuan } = useSatuanActions();

  const fetchSatuan = async () => {
    try {
      const response = await apiClient.get("/units/");
      return response.data;
    } catch (err) {
      setToast({
        message: "Data gagal diambil",
        type: "error",
      });
      console.error("Fetch error:", err);
      return [];
    }
  };

  const reloadSatuan = async () => {
    const data = await fetchSatuan();
    let items = [];

    if (Array.isArray(data)) {
      items = data;
    } else if (Array.isArray(data?.data)) {
      items = data.data;
    } else {
      setToast({ message: "Terdeteksi format tak dikenal", type: "error" });
      items = [];
    }

    setSatuanList(items);
  };

  useEffect(() => {
    reloadSatuan().finally(() => setLoading(false));
  }, []);

  const openAddModal = () => {
    setModalMode("add");
    setModalSatuan(null);
    setModalOpen(true);
  };

  const openEditModal = async (id) => {
    try {
      const satuan = await getSatuanById(id);
      setModalSatuan(satuan);
      setModalMode("edit");
      setModalOpen(true);
    } catch (err) {
      setToast({
        message: "Data satuan gagal diambil",
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
      await deleteSatuan(deleteTargetId);
      await reloadSatuan();
      setToast({ message: "Data berhasil dihapus", type: "success" });
    } catch (err) {
      alert("Gagal menghapus data.");
    } finally {
      setIsConfirmOpen(false);
      setDeleteTargetId(null);
    }
  };

  const handleModalSuccess = async () => {
    await reloadSatuan();
    setToast({
      message:
        modalMode === "edit"
          ? "Data berhasil diperbarui"
          : "Data berhasil ditambahkan",
      type: "success",
    });
    setModalOpen(false);
    setModalSatuan(null);
  };

  const { searchTerm, setSearchTerm, filteredData } = useSearch(satuanList, [
    "name",
  ]);

  const columns = [
    { header: "Nama", accessor: "name" },
    { header: "Deskripsi", accessor: "description" },
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
        <h1 className="text-2xl font-bold mb-6">Daftar Satuan Obat</h1>
        <SearchBar value={searchTerm} onChange={setSearchTerm} />

        <div className="border-1 rounded-md border-gray-300 bg-white p-5">
          <div className="flex gap-4 mb-4">
            <button
              onClick={openAddModal}
              className="flex items-center text-blue-700 font-semibold space-x-1 bg-transparent border border-blue-700 py-2 px-4 rounded-md"
            >
              <PlusIcon className="w-4 h-4" />
              <span>Tambah Satuan</span>
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

      <SatuanModal
        isOpen={modalOpen}
        close={() => setModalOpen(false)}
        onSuccess={handleModalSuccess}
        mode={modalMode}
        satuan={modalSatuan}
      />

      <ConfirmDialog
        isOpen={isConfirmOpen}
        title="Konfirmasi Penghapusan"
        description="Apakah Anda yakin ingin menghapus satuan ini?"
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

export default AturSatuanPage;
