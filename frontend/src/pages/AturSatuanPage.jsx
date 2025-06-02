import React, { useEffect, useState } from "react";
import SearchBar from "../components/SearchBar";
import Button from "../components/buttonComp";
import DataTable from "../components/tableCompo";
import ActionMenu from "../components/ActionMenu";
import AddUnitModal from "../components/modal/addUnitModal";
import EditUnitModal from "../components/modal/editUnitModal";
import useSearch from "../hooks/useSearch";
import useUnitActions from "../hooks/useUnitAction";
import ConfirmDialog from "../components/ConfirmDialog";
import Sidebar from "../components/Sidebar";
import { apiClient } from "../config/api";
import { PlusIcon } from '@heroicons/react/24/solid';
import Toast from "../components/toast";

function AturSatuanPage() {
  const [units, setUnits] = useState([]);
  const [isAddOpen, setIsAddOpen] = useState(false);
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [editId, setEditId] = useState(null);
  const [toast, setToast] = useState(null);

  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState(null);

  const { getUnitById, deleteUnit } = useUnitActions();

  const fetchUnits = async () => {
    try {
      const response = await apiClient.get("/units/");
      return response.data;
    } catch (error) {
      throw error;
    }
    };

    useEffect(() => {
      const loadUnits = async () => {
        try {
          const data = await fetchUnits();

          if (Array.isArray(data)) {
            setUnits(data);
          } else if (Array.isArray(data?.data)) {
            setUnits(data.data);
          } else {
            setUnits([]);
          }

        } catch (error) {
          setUnits([]);
        }
      };

      loadUnits();
    }, []);

  const handleEdit = (id) => {
    setEditId(id);
    setIsEditOpen(true);
  };

  const handleDeleteRequest = (unit) => {
    setDeleteTarget(unit);
    setIsConfirmOpen(true);
  };

  const handleConfirmDelete = async () => {
    try {
      await deleteUnit(deleteTarget.id);
      await fetchUnits();
    } catch (err) {
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

  const data = await fetchUnits();
    if (Array.isArray(data)) {
      setUnits(data);
    } else if (Array.isArray(data?.data)) {
      setUnits(data.data);
    }

    closeModal();
  };

  const handleEditSuccess = () =>
    handleSuccess({
      message: "Berhasil mengedit satuan.",
      closeModal: () => setIsEditOpen(false),
    });

  const handleAddSuccess = () =>
    handleSuccess({
      message: "Berhasil menambahkan satuan.",
      closeModal: () => setIsAddOpen(false),
    });

  const { searchTerm, setSearchTerm, filteredData } = useSearch(units, ["name"]);

  const columns = [
    { header: "Satuan", accessor: "name" },
    { header: "Deskripsi", accessor: "description" },
    {
      header: "Aksi",
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
        <h1 className="text-2xl font-bold mb-6">Daftar Satuan</h1>

        <SearchBar
          value={searchTerm}
          onChange={setSearchTerm}
          placeholder="Cari Satuan"
        />

        <div className="border border-gray-400 rounded-xl bg-white p-5">

          <div className="flex justify-start py-5">
            <button onClick={() => setIsAddOpen(true)}
              className="flex items-center text-blue-700 font-semibold space-x-1 bg-transparent border border-blue-700 py-2 px-4 rounded-md"
            >
              <PlusIcon className="w-4 h-4" />
                <span>Tambah Satuan</span>
            </button>
          </div>
        <DataTable columns={columns} data={filteredData} showIndex={true} />
        </div>

        <AddUnitModal
          isOpen={isAddOpen}
          close={() => setIsAddOpen(false)}
          onSuccess={handleAddSuccess}
        />

        <EditUnitModal
          isOpen={isEditOpen}
          close={() => setIsEditOpen(false)}
          unitId={editId}
          onSuccess={handleEditSuccess}
        />

        <ConfirmDialog
          isOpen={isConfirmOpen}
          title={`Konfirmasi penghapusan ${deleteTarget?.nama}`}
          description={`Apakah Anda yakin ingin menghapus ${deleteTarget?.nama} dari daftar satuan?`}
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
  );
}

export default AturSatuanPage;
