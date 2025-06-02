import React, { useEffect, useState } from "react";
import SearchBar from "../components/SearchBar";
import DataTable from "../components/tableCompo";
import ActionMenu from "../components/ActionMenu";
import AddUserModal from "../components/modal/addUserModal";
import EditUserModal from "../components/modal/editUserModal";
import useSearch from "../hooks/useSearch";
import useUserActions from "../hooks/useUserAction";
import ConfirmDialog from "../components/ConfirmDialog";
import Sidebar from "../components/Sidebar";
import { apiClient } from "../config/api";
import { PlusIcon } from '@heroicons/react/24/solid';
import Toast from "../components/toast";

function AturUsersPage() {
  const [users, setUsers] = useState([]);
  const [isAddOpen, setIsAddOpen] = useState(false);
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [editId, setEditId] = useState(null);
  const [toast, setToast] = useState(null);

  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState(null);

  const { getUserById, deleteUser, reactivateUser, deactivateUser } = useUserActions();

  const fetchUsers = async () => {
    try {
      const response = await apiClient.get("/users/");
      return response.data;
    } catch (error) {
      console.error("Failed to fetch users", error);
      throw error;
    }
    };

    useEffect(() => {
      const loadUsers = async () => {
        try {
          const data = await fetchUsers();

          if (Array.isArray(data)) {
            setUsers(data);
          } else if (Array.isArray(data?.data)) {
            setUsers(data.data);
          } else {
            console.error("Unexpected data format:", data);
            setUsers([]);
          }

        } catch (error) {
          console.error("Error loading users", error);
          setUsers([]);
        }
      };

      loadUsers();
    }, []);

  const handleEdit = (id) => {
    setEditId(id);
    setIsEditOpen(true);
  };

  const handleDeleteRequest = (user) => {
    setDeleteTarget(user);
    setIsConfirmOpen(true);
  };

  const handleConfirmDelete = async () => {
    try {
      await deleteUser(deleteTarget.id);
      const data = await fetchUsers();
        if (Array.isArray(data)) {
        setUsers(data);
        } else if (Array.isArray(data?.data)) {
        setUsers(data.data);
        }
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

  const handleActivateRequest = async (id) => {
    try {
        await reactivateUser(id);
        setToast({ message: "User berhasil diaktifkan.", type: "success" });

        const data = await fetchUsers();
        if (Array.isArray(data)) {
        setUsers(data);
        } else if (Array.isArray(data?.data)) {
        setUsers(data.data);
        }
    } catch (err) {
        setToast({
        message: err.response?.data?.message || err.message || "Gagal mengaktifkan user.",
        type: "error",
        });
    }
    };

const handleDeactivateRequest = async (id) => {
  console.log("Deactivating user ID:", id);

  try {
    await deactivateUser(id);
    setToast({ message: "User berhasil dinonaktifkan.", type: "success" });

    const data = await fetchUsers();
    if (Array.isArray(data)) {
      setUsers(data);
    } else if (Array.isArray(data?.data)) {
      setUsers(data.data);
    }
  } catch (err) {
    console.error("Deactivation error:", err);
    setToast({
      message: err.response?.data?.message || err.message || "Gagal menonaktifkan user.",
      type: "error",
    });
  }
};


  const handleSuccess = async ({ message, closeModal }) => {
  setToast({ message, type: "success" });

  const data = await fetchUsers();
    if (Array.isArray(data)) {
      setUsers(data);
    } else if (Array.isArray(data?.data)) {
      setUsers(data.data);
    }

    closeModal();
  };

  const handleEditSuccess = () =>
    handleSuccess({
      message: "Berhasil mengedit user.",
      closeModal: () => setIsEditOpen(false),
    });

  const handleAddSuccess = () =>
    handleSuccess({
      message: "Berhasil menambahkan user.",
      closeModal: () => setIsAddOpen(false),
    });

  const { searchTerm, setSearchTerm, filteredData } = useSearch(users, ["full_name"]);

  const columns = [
    { header: "Nama", accessor: "full_name" },
    { header: "Peran", accessor: "role" },
    { header: "Nomor HP", accessor: "phone" },
    { header: "EMail", accessor: "email" },
    {
        header: "Status",
        accessor: "active",
        render: (row) => (
            <span
            className={`px-2 py-1 rounded text-xs font-semibold
                ${row.active ? "bg-green-100 text-green-800" : "bg-red-100 text-red-800"}`}
            >
            {row.active ? "Aktif" : "Nonaktif"}
            </span>
        )
    },
    {
        header: "Pilih Aksi",
        accessor: "actions",
        isAction: true,
        render: (row) => {
            const baseActions = [
            { label: "Edit", onClick: () => handleEdit(row.id)},
            { label: "Hapus", onClick: () => handleDeleteRequest(row) },
            ];

            const statusAction = row.active
            ? { label: "Non-Aktifkan", onClick: () => handleDeactivateRequest(row.id) }
            : { label: "Aktifkan", onClick: () => handleActivateRequest(row.id) };

            return <ActionMenu actions={[...baseActions, statusAction]} />;
        },
    }
  ];

  return (
    <div className="flex">
      <div className="bg-white min-h-screen">
      <Sidebar />
      </div>
      <div className="p-5 w-full py-10">
        <h1 className="text-2xl font-bold mb-6">Daftar User</h1>
          <SearchBar
            value={searchTerm}
            onChange={setSearchTerm}
            placeholder="Cari User"
          />

        <div className="border border-gray-400 rounded-xl bg-white p-5">

          <div className="flex justify-start py-5">
            <button
              onClick={() => setIsAddOpen(true)}
              className="flex items-center text-blue-700 font-semibold space-x-1 bg-transparent border border-blue-700 py-2 px-4 rounded-md"
            >
              <PlusIcon className="w-4 h-4" />
              <span>Tambah User</span>
            </button>
          </div>

          <DataTable columns={columns} data={filteredData} showIndex={true} />

          <AddUserModal
            isOpen={isAddOpen}
            close={() => setIsAddOpen(false)}
            onSuccess={handleAddSuccess}
          />

          <EditUserModal
            isOpen={isEditOpen}
            close={() => setIsEditOpen(false)}
            userId={editId}
            onSuccess={handleEditSuccess}
          />

          <ConfirmDialog
            isOpen={isConfirmOpen}
            title={`Konfirmasi penghapusan ${deleteTarget?.name}`}
            description={`Apakah Anda yakin ingin menghapus ${deleteTarget?.name} dari daftar user?`}
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

export default AturUsersPage;
