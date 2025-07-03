import React, { useState, useEffect } from "react";
import SearchBar from "../../components/SearchBar";
import DataTable from "../../components/tableCompo";
import ActionMenu from "../../components/ActionMenu";
import AddStorageLocationModal from "../../components/modal/AddStorageLocationModal";
import EditStorageLocationModal from "../../components/modal/EditStorageLocationModal";
import useStorageLocationActions from "../../hooks/useStorageLocationActions";
import ConfirmDialog from "../../components/ConfirmDialog";
import Sidebar from "../../components/Sidebar";
import { PlusIcon } from "@heroicons/react/24/solid";
import Toast from "../../components/toast";
import Pagination from "../../components/Pagination";
import useServerSideSearch from "../../hooks/useServerSideSearch";

function StorageLocationPage() {
  const {
    items: storageLocations,
    isLoading,
    currentPage,
    totalPages,
    totalData,
    searchInput,
    handleSearchInputChange,
    handlePageChange,
    setCurrentPage: setCurrentPageHook,
    ITEMS_PER_PAGE,
    refetch,
  } = useServerSideSearch("/storage-locations");

  const [isAddOpen, setIsAddOpen] = useState(false);
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [editId, setEditId] = useState(null);
  const [toast, setToast] = useState(null);
  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState(null);

  const { deleteStorageLocation } = useStorageLocationActions();

  const handleEdit = (id) => {
    setEditId(id);
    setIsEditOpen(true);
  };

  const handleDeleteRequest = (location) => {
    setDeleteTarget(location);
    setIsConfirmOpen(true);
  };

  const handleConfirmDelete = async () => {
    if (!deleteTarget) return;
    try {
      await deleteStorageLocation(deleteTarget.id);
      setToast({
        message: `Lokasi "${deleteTarget.name}" berhasil dihapus.`,
        type: "success",
      });
      if (storageLocations.length === 1 && currentPage > 1) {
        setCurrentPageHook(currentPage - 1);
      } else {
        refetch();
      }
    } catch (err) {
      console.error("Delete failed:", err);
      setToast({
        message: `Gagal menghapus lokasi "${deleteTarget.name}".`,
        type: "error",
      });
    } finally {
      setIsConfirmOpen(false);
      setDeleteTarget(null);
    }
  };

  const handleCancelDelete = () => {
    setIsConfirmOpen(false);
    setDeleteTarget(null);
  };

  const handleModalSuccess = ({ message, closeModal }) => {
    setToast({ message, type: "success" });
    refetch();
    if (closeModal) closeModal();
  };

  const handleAddSuccess = (result) => {
    handleModalSuccess({
      message: result.message,
      closeModal: () => setIsAddOpen(false),
    });
  };

  const handleEditSuccess = (result) => {
    handleModalSuccess({
      message: result.message,
      closeModal: () => setIsEditOpen(false),
    });
  };

  const columns = [
    { header: "Nama Lokasi", accessor: "name" },
    { header: "Deskripsi", accessor: "description" },
    { header: "Dibuat Oleh", accessor: "created_by_name" },
    {
      header: "Pilih Aksi",
      accessor: "actions",
      isAction: true,
      render: (row) => (
        <ActionMenu
          actions={[
            { label: "Edit", onClick: () => handleEdit(row.id) },
            { label: "Hapus", onClick: () => handleDeleteRequest(row) },
          ]}
        />
      ),
    },
  ];

  return (
    <div className="flex h-screen overflow-hidden">
      <div className="bg-white min-h-screen">
        <Sidebar />
      </div>
      <div className="flex-1 flex flex-col p-5 py-10 overflow-y-auto">
        <h1 className="text-2xl font-bold mb-6 shrink-0">
          Daftar Lokasi Penyimpanan
        </h1>
        <div className="mb-4 shrink-0">
          <SearchBar
            value={searchInput}
            onChange={handleSearchInputChange}
            placeholder="Cari Lokasi Penyimpanan..."
          />
        </div>

        <div className="border border-gray-300 rounded-xl bg-white p-5 shadow flex flex-col flex-grow">
          <div className="flex justify-start py-5 shrink-0">
            <button
              onClick={() => setIsAddOpen(true)}
              className="flex items-center text-blue-700 font-semibold space-x-1 bg-transparent border border-blue-700 hover:bg-blue-50 py-2 px-4 rounded-md transition duration-150"
            >
              <PlusIcon className="w-5 h-5" />
              <span>Tambah Lokasi</span>
            </button>
          </div>

          <div className="flex-grow overflow-auto">
            {isLoading && <p className="text-center py-4">Memuat data...</p>}
            {!isLoading && storageLocations.length === 0 && totalData === 0 && (
              <p className="text-center py-4 text-gray-500">
                {searchInput
                  ? `Tidak ada hasil untuk "${searchInput}".`
                  : "Tidak ada data lokasi penyimpanan."}
              </p>
            )}
            {!isLoading && totalData > 0 && (
              <DataTable
                columns={columns}
                data={storageLocations}
                showIndex={true}
                currentPage={currentPage}
                itemsPerPage={ITEMS_PER_PAGE}
              />
            )}
          </div>

          {totalData > 0 && (
            <Pagination
              currentPage={currentPage}
              totalPages={totalPages}
              onPageChange={handlePageChange}
              itemsPerPage={ITEMS_PER_PAGE}
              totalData={totalData}
            />
          )}
        </div>

        {toast && (
          <Toast
            message={toast.message}
            type={toast.type}
            onClose={() => setToast(null)}
          />
        )}
      </div>
      {isAddOpen && (
        <AddStorageLocationModal
          isOpen={isAddOpen}
          close={() => setIsAddOpen(false)}
          onSuccess={handleAddSuccess}
        />
      )}
      {editId && isEditOpen && (
        <EditStorageLocationModal
          isOpen={isEditOpen}
          close={() => {
            setIsEditOpen(false);
            setEditId(null);
          }}
          storageLocationId={editId}
          onSuccess={handleEditSuccess}
        />
      )}
      {deleteTarget && isConfirmOpen && (
        <ConfirmDialog
          isOpen={isConfirmOpen}
          title={`Konfirmasi Penghapusan`}
          description={`Apakah Anda yakin ingin menghapus lokasi "${deleteTarget.name}"?`}
          onCancel={handleCancelDelete}
          onConfirm={handleConfirmDelete}
          confirmText="Hapus"
          cancelText="Batal"
        />
      )}
    </div>
  );
}

export default StorageLocationPage;
