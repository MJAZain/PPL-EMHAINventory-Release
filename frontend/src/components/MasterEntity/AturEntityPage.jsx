import React, { useState, useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { apiClient } from "../../config/api";

import useSearch from "../../hooks/useSearch";
import ActionMenu from "../../components/ActionMenu";
import SearchBar from "../../components/SearchBar";
import DataTable from "../../components/tableCompo";
import Sidebar from "../../components/Sidebar";
import ConfirmDialog from "../../components/ConfirmDialog";
import Toast from "../../components/toast";
import { PlusIcon } from "@heroicons/react/24/solid";

const AturEntityPage = ({
  endpoint,
  entityName = "Item",
  entityKey = "name",
  columns,
  ModalComponent,
}) => {
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [toast, setToast] = useState(null);

  const [modalOpen, setModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState("add");
  const [modalData, setModalData] = useState(null);

  const [confirmOpen, setConfirmOpen] = useState(false);
  const [deleteId, setDeleteId] = useState(null);

  const [actionLoading, setActionLoading] = useState(false);
  const [actionError, setActionError] = useState(null);

  const handleRequest = useCallback(async (requestFn) => {
    setActionLoading(true);
    setActionError(null);
    try {
      return await requestFn();
    } catch (err) {
      const message = err.response?.data?.message || err.message || "Unknown error";
      setActionError(message);
      throw err;
    } finally {
      setActionLoading(false);
    }
  }, []);

  const getItemById = useCallback(
    (id) =>
      handleRequest(() =>
        apiClient.get(`${endpoint}${id}`).then((res) => res.data.data || res.data)
      ),
    [endpoint, handleRequest]
  );

  const deleteItem = useCallback(
    (id) => handleRequest(() => apiClient.delete(`${endpoint}${id}`)),
    [endpoint, handleRequest]
  );

  const fetchData = async () => {
    try {
      const response = await apiClient.get(endpoint);
      const data = response.data;
      return Array.isArray(data) ? data : data?.data || [];
    } catch (err) {
      setToast({ message: `Gagal memuat ${entityName.toLowerCase()}`, type: "error" });
      return [];
    }
  };

  const reload = async () => {
    setLoading(true);
    const data = await fetchData();
    setItems(data);
    setLoading(false);
  };

  useEffect(() => {
    reload();
  }, []);

  const openAdd = () => {
    setModalMode("add");
    setModalData(null);
    setModalOpen(true);
  };

  const openEdit = async (id) => {
    try {
      const item = await getItemById(id);
      setModalData(item);
      setModalMode("edit");
      setModalOpen(true);
    } catch {
      setToast({ message: `Gagal mengambil ${entityName.toLowerCase()}`, type: "error" });
    }
  };

  const requestDelete = (id) => {
    setDeleteId(id);
    setConfirmOpen(true);
  };

  const cancelDelete = () => {
    setConfirmOpen(false);
    setDeleteId(null);
  };

  const confirmDelete = async () => {
    try {
      await deleteItem(deleteId);
      await reload();
      setToast({ message: `${entityName} berhasil dihapus`, type: "success" });
    } catch {
      setToast({ message: `Gagal menghapus ${entityName.toLowerCase()}`, type: "error" });
    } finally {
      setConfirmOpen(false);
      setDeleteId(null);
    }
  };

  const handleSuccess = async () => {
    await reload();
    setToast({
      message:
        modalMode === "edit"
          ? `${entityName} berhasil diperbarui`
          : `${entityName} berhasil ditambahkan`,
      type: "success",
    });
    setModalOpen(false);
    setModalData(null);
  };

  const { searchTerm, setSearchTerm, filteredData } = useSearch(items, [entityKey]);

  const enhancedColumns = [
    ...columns,
    {
      header: "Pilih Aksi",
      accessor: "actions",
      isAction: true,
      render: (item) => (
        <ActionMenu
          actions={[
            { label: "Edit", onClick: () => openEdit(item.id) },
            { label: "Delete", onClick: () => requestDelete(item.id) },
          ]}
        />
      ),
    },
  ];

  return (
    <div className="flex">
        <div className="bg-white min-h-screen">
            <Sidebar />
        </div>
      
      <div className="p-5 w-full py-10">
        <h1 className="text-2xl font-bold mb-6">Daftar {entityName}</h1>
        <SearchBar value={searchTerm} onChange={setSearchTerm} />

        <div className="border-1 rounded-md border-gray-300 bg-white p-5">
          <div className="flex gap-4 mb-4">
            <button
              onClick={openAdd}
              className="flex items-center text-blue-700 font-semibold space-x-1 bg-transparent border border-blue-700 py-2 px-4 rounded-md"
            >
              <PlusIcon className="w-4 h-4" />
              <span>Tambah {entityName}</span>
            </button>
          </div>

          <DataTable columns={enhancedColumns} data={filteredData} showIndex={true} />
        </div>
      </div>

      <ModalComponent
        isOpen={modalOpen}
        close={() => setModalOpen(false)}
        onSuccess={handleSuccess}
        mode={modalMode}
        kategori={modalData} // optional: rename prop to `data` for generality
      />

      <ConfirmDialog
        isOpen={confirmOpen}
        title="Konfirmasi Penghapusan"
        description={`Apakah Anda yakin ingin menghapus ${entityName.toLowerCase()} ini?`}
        onCancel={cancelDelete}
        onConfirm={confirmDelete}
        loading={actionLoading}
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
};

export default AturEntityPage;
