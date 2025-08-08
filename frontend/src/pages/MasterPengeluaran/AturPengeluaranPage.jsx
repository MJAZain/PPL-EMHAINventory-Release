import React, { useState, useEffect } from "react";
import useSearch from "../../hooks/useSearch";
import usePengeluaranActions from "./usePengeluaranAction";
import { useNavigate } from "react-router-dom";

import { apiClient } from "../../config/api";
import { PlusIcon } from "@heroicons/react/24/solid";

import ActionMenu from "../../components/ActionMenu";
import SearchBar from "../../components/SearchBar";
import DataTable from "../../components/tableCompo";

import PengeluaranModal from "./PengeluaranModal";
import Sidebar from "../../components/Sidebar";
import ConfirmDialog from "../../components/ConfirmDialog";
import Toast from "../../components/toast";

function AturPengeluaranPage() {
  const [pengeluaranList, setPengeluaranList] = useState([]);
  const [loading, setLoading] = useState(true);
  const [toast, setToast] = useState(null);

  const [modalOpen, setModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState("add");
  const [modalPengeluaran, setModalPengeluaran] = useState(null);

  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [deleteTargetId, setDeleteTargetId] = useState(null);

  const navigate = useNavigate();
  const { getPengeluaranById, deletePengeluaran } = usePengeluaranActions();

    const formatDateTime = (isoString) => {
    if (!isoString) return "-";
    const date = new Date(isoString);
    return date.toLocaleString("id-ID", {
      day: "2-digit",
      month: "short",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const fetchPengeluaran = async () => {
    try {
      const response = await apiClient.get("/expenses/");
      return response.data?.data || [];
    } catch (err) {
      setToast({ message: "Data gagal diambil", type: "error" });
      console.error("Fetch error:", err);
      return [];
    }
  };

  const reloadPengeluaran = async () => {
    const items = await fetchPengeluaran();
    setPengeluaranList(items);
  };

  useEffect(() => {
    reloadPengeluaran().finally(() => setLoading(false));
  }, []);

  const openAddModal = () => {
    setModalMode("add");
    setModalPengeluaran(null);
    setModalOpen(true);
  };

  const openEditModal = async (id) => {
    try {
      const pengeluaran = await getPengeluaranById(id);
      setModalPengeluaran(pengeluaran);
      setModalMode("edit");
      setModalOpen(true);
    } catch (err) {
      setToast({ message: "Pengeluaran gagal diambil", type: "error" });
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
      await deletePengeluaran(deleteTargetId);
      await reloadPengeluaran();
      setToast({ message: "Pengeluaran berhasil dihapus", type: "success" });
    } catch (err) {
      setToast({ message: "Gagal menghapus pengeluaran.", type: "error" });
    } finally {
      setIsConfirmOpen(false);
      setDeleteTargetId(null);
    }
  };

  const handleModalSuccess = async () => {
    await reloadPengeluaran();
    setToast({
      message:
        modalMode === "edit"
          ? "Data berhasil diperbarui"
          : "Data berhasil ditambahkan",
      type: "success",
    });
    setModalOpen(false);
    setModalPengeluaran(null);
  };

  const { searchTerm, setSearchTerm, filteredData } = useSearch(
    pengeluaranList,
    ["description", "amount", "date", "expense_type.name"]
  );

  const columns = [
    {
    header: "Tanggal",
    accessor: "date",
    render: (item) => formatDateTime(item.date),
  },
    { header: "Jenis", accessor: "expense_type.name" },
    { header: "Jumlah", accessor: "amount" },
    { header: "Deskripsi", accessor: "description" },
    {
      header: "Aksi",
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
        <h1 className="text-2xl font-bold mb-6">Daftar Pengeluaran</h1>
        <SearchBar value={searchTerm} onChange={setSearchTerm} />

        <div className="border-1 rounded-md border-gray-300 bg-white p-5">
          <div className="flex gap-4 mb-4">
            <button
              onClick={openAddModal}
              className="flex items-center text-blue-700 font-semibold space-x-1 bg-transparent border border-blue-700 py-2 px-4 rounded-md"
            >
              <PlusIcon className="w-4 h-4" />
              <span>Tambah Pengeluaran</span>
            </button>
          </div>

          <div className="max-w-[1121px]">
            <DataTable columns={columns} data={filteredData} showIndex={true} />
          </div>
        </div>
      </div>

      <PengeluaranModal
        isOpen={modalOpen}
        close={() => setModalOpen(false)}
        onSuccess={handleModalSuccess}
        mode={modalMode}
        pengeluaran={modalPengeluaran}
      />

      <ConfirmDialog
        isOpen={isConfirmOpen}
        title="Konfirmasi Penghapusan"
        description="Apakah Anda yakin ingin menghapus pengeluaran ini?"
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

export default AturPengeluaranPage;
