import React, { useEffect, useState } from "react";
import SearchBar from "../components/SearchBar";
import Button from "../components/buttonComp";
import DataTable from "../components/tableCompo";
import ActionMenu from "../components/ActionMenu";
import AddCategoryModal from "../components/modal/addCategoryModal";
import EditCategoryModal from "../components/modal/editCategoryModal";
import useSearch from "../hooks/useSearch";
import useCategoryActions from "../hooks/useCategoryAction";
import ConfirmDialog from "../components/ConfirmDialog";
import Sidebar from "../components/Sidebar";
import { apiClient } from "../config/api";
import { PlusIcon } from '@heroicons/react/24/solid';
import Toast from "../components/toast";

function AturLokasiPage() {
  const [categories, setCategories] = useState([]);
  const [isAddOpen, setIsAddOpen] = useState(false);
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [editId, setEditId] = useState(null);
  const [toast, setToast] = useState(null);

  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState(null);

  const { getCategoryById, deleteCategory } = useCategoryActions();

  const fetchCategories = async () => {
    try {
      const response = await apiClient.get("/categories/");
      return response.data;
    } catch (error) {
      console.error("Failed to fetch categories", error);
      throw error;
    }
    };

    useEffect(() => {
      const loadCategories = async () => {
        try {
          const data = await fetchCategories();

          if (Array.isArray(data)) {
            setCategories(data);
          } else if (Array.isArray(data?.data)) {
            setCategories(data.data);
          } else {
            console.error("Unexpected data format:", data);
            setCategories([]);
          }

        } catch (error) {
          console.error("Error loading categories", error);
          setCategories([]);
        }
      };

      loadCategories();
    }, []);

  const handleEdit = (id) => {
    setEditId(id);
    setIsEditOpen(true);
  };

  const handleDeleteRequest = (category) => {
    setDeleteTarget(category);
    setIsConfirmOpen(true);
  };

  const handleConfirmDelete = async () => {
    try {
      await deleteCategory(deleteTarget.id);
      await fetchCategories();
    } catch (err) {
      console.error("Delete failed:", err);
    } finally {
      setIsConfirmOpen(false);
      setDeleteTarget(null);
    }
  };

  const handleCancelDelete = () => {
    setIsConfirmOpen(false);
    setDeleteTarget(null);
  };

  const handleSuccess = async ({ message, closeModal }) => {
  setToast({ message, type: "success" });

  const data = await fetchCategories();
    if (Array.isArray(data)) {
      setCategories(data);
    } else if (Array.isArray(data?.data)) {
      setCategories(data.data);
    }

    closeModal();
  };

  const handleEditSuccess = () =>
    handleSuccess({
      message: "Berhasil mengedit kategori.",
      closeModal: () => setIsEditOpen(false),
    });

  const handleAddSuccess = () =>
    handleSuccess({
      message: "Berhasil menambahkan kategori.",
      closeModal: () => setIsAddOpen(false),
    });

  const { searchTerm, setSearchTerm, filteredData } = useSearch(categories, ["name"]);

  const columns = [
    { header: "Kategori", accessor: "name" },
    { header: "Deskripsi", accessor: "description" },
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
    <div className="flex">
      <div className="bg-white min-h-screen">
      <Sidebar />
      </div>
      <div className="p-5 w-full py-10">
        <h1 className="text-2xl font-bold mb-6">Daftar Kategori</h1>
          <SearchBar
            value={searchTerm}
            onChange={setSearchTerm}
            placeholder="Cari Kategori"
          />

        <div className="border border-gray-400 rounded-xl bg-white p-5">

          <div className="flex justify-start py-5">
            <button
              onClick={() => setIsAddOpen(true)}
              className="flex items-center text-blue-700 font-semibold space-x-1 bg-transparent border border-blue-700 py-2 px-4 rounded-md"
            >
              <PlusIcon className="w-4 h-4" />
              <span>Tambah Kategori</span>
            </button>
          </div>

          <DataTable columns={columns} data={filteredData} showIndex={true} />

          <AddCategoryModal
            isOpen={isAddOpen}
            close={() => setIsAddOpen(false)}
            onSuccess={handleAddSuccess}
          />

          <EditCategoryModal
            isOpen={isEditOpen}
            close={() => setIsEditOpen(false)}
            categoryId={editId}
            onSuccess={handleEditSuccess}
          />

          <ConfirmDialog
            isOpen={isConfirmOpen}
            title={`Konfirmasi penghapusan ${deleteTarget?.name}`}
            description={`Apakah Anda yakin ingin menghapus ${deleteTarget?.name} dari daftar kategori?`}
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
      </div>
    </div>
  );
}

export default AturKategoriPage;
