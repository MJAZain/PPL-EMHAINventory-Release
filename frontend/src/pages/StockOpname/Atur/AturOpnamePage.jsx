import React, { useState, useEffect } from "react";
import useSearch from "../../../hooks/useSearch";
import useOpnameActions from "../useOpnameAction";
import { useNavigate } from "react-router-dom";

import { apiClient } from "../../../config/api";
import { PlusIcon } from "@heroicons/react/24/solid";

import ActionMenu from "../../../components/ActionMenu";
import SearchBar from "../../../components/SearchBar";
import DataTable from "../../../components/tableCompo";

import OpnameModal from "../OpnameModal";
import Sidebar from "../../../components/Sidebar";
import ConfirmDialog from "../../../components/ConfirmDialog";
import Toast from "../../../components/toast";

import { formatDateTime } from "../../../utils/formatter";

function AturOpnamePage() {
  const [opnameList, setOpnameList] = useState([]);
  const [loading, setLoading] = useState(true);
  const [toast, setToast] = useState(null);

  const [modalOpen, setModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState("add");
  const [modalOpname, setModalOpname] = useState(null);

  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [deleteTargetId, setDeleteTargetId] = useState(null);

  const navigate = useNavigate();
  const { getOpnameById, deleteOpname } = useOpnameActions();

  const fetchOpnames = async () => {
    try {
      const response = await apiClient.get("/stock-opname");
      return response.data;
    } catch (err) {
      setToast({ message: "Data gagal diambil", type: "error" });
      console.error("Fetch error:", err);
      return [];
    }
  };

  const reloadOpnames = async () => {
    const data = await fetchOpnames();
    const items = Array.isArray(data) ? data : data?.data || [];
    const draftItems = items.filter(item => item.status === "draft");
    setOpnameList(draftItems);
  };

  useEffect(() => {
    reloadOpnames().finally(() => setLoading(false));
  }, []);

  const openAddModal = () => {
    setModalMode("add");
    setModalOpname(null);
    setModalOpen(true);
  };

  const openEditModal = async (id) => {
    try {
      const opname = await getOpnameById(id);
      setModalOpname(opname);
      setModalMode("edit");
      setModalOpen(true);
    } catch (err) {
      setToast({ message: "Opname gagal diambil", type: "error" });
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
      await deleteOpname(deleteTargetId);
      await reloadOpnames();
      setToast({ message: "Opname berhasil dihapus", type: "success" });
    } catch (err) {
      alert("Gagal menghapus opname.");
    } finally {
      setIsConfirmOpen(false);
      setDeleteTargetId(null);
    }
  };

  const handleModalSuccess = async () => {
    await reloadOpnames();
    setToast({
      message:
        modalMode === "edit"
          ? "Data berhasil diperbarui"
          : "Data berhasil ditambahkan",
      type: "success",
    });
    setModalOpen(false);
    setModalOpname(null);
  };

  const { searchTerm, setSearchTerm, filteredData } = useSearch(opnameList, [
    "item_name",
    "location",
    "category",
  ]);

  const columns = [
    {
    header: "Tanggal",
    accessor: "opname_date",
    render: (item) => formatDateTime(item.opname_date),
    },
    {
      header: "Catatan",
      accessor: "notes",
    },
     {
    header: "Mulai Opname",
    accessor: "start", // you can name it anything
    render: ({ row }) => {
      const handleStartOpname = () => {
        const opnameId = row.original.opname_id; // or row.original.id if that’s what you have
        localStorage.setItem("opnameId", opnameId);
        navigate("/opname-details");
      };

      return (
        <button
          onClick={handleStartOpname}
          className="text-blue-600 hover:underline"
        >
          Mulai Opname
        </button>
      );
    }
  },
    {
      header: "Pilih Aksi",
      accessor: "actions",
      isAction: true,
      render: (item) => (
        <ActionMenu
          actions={[
            { label: "Edit", onClick: () => openEditModal(item.opname_id) },
            { label: "Hapus", onClick: () => handleDeleteRequest(item.opname_id) },
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
        <h1 className="text-2xl font-bold mb-6">Daftar Draft Stock Opname</h1>
        <SearchBar value={searchTerm} onChange={setSearchTerm} />

        <div className="border-1 rounded-md border-gray-300 bg-white p-5">
          <div className="flex gap-4 mb-4">
            <button
              onClick={openAddModal}
              className="flex items-center text-blue-700 font-semibold space-x-1 bg-transparent border border-blue-700 py-2 px-4 rounded-md"
            >
              <PlusIcon className="w-4 h-4" />
              <span>Tambah Opname</span>
            </button>
          </div>

          <div className="max-w-[1121px]">
            <DataTable columns={columns} data={filteredData} showIndex={true} />
          </div>
        </div>
      </div>

      <OpnameModal
        isOpen={modalOpen}
        close={() => setModalOpen(false)}
        onSuccess={handleModalSuccess}
        mode={modalMode}
        opname={modalOpname}
      />

      <ConfirmDialog
        isOpen={isConfirmOpen}
        title="Konfirmasi Penghapusan"
        description="Apakah Anda yakin ingin menghapus opname ini?"
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

export default AturOpnamePage;
