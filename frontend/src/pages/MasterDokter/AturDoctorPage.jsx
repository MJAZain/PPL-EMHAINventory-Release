import React, { useState, useEffect } from "react";
import useSearch from "../../hooks/useSearch";
import useDoctorActions from "./useDoctorAction";
import { useNavigate } from "react-router-dom";

import { apiClient } from "../../config/api";
import { PlusIcon } from "@heroicons/react/24/solid";

import ActionMenu from "../../components/ActionMenu";
import SearchBar from "../../components/SearchBar";
import DataTable from "../../components/tableCompo";

import DoctorModal from "./DoctorModal";
import Sidebar from "../../components/Sidebar";
import ConfirmDialog from "../../components/ConfirmDialog";
import Toast from "../../components/toast";

function AturDoctorPage() {
  const [doctorList, setDoctorList] = useState([]);
  const [loading, setLoading] = useState(true);
  const [toast, setToast] = useState(null);

  const [modalOpen, setModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState("add");
  const [modalDoctor, setModalDoctor] = useState(null);

  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [deleteTargetId, setDeleteTargetId] = useState(null);

  const navigate = useNavigate();
  const { getDoctorById, deleteDoctor } = useDoctorActions();

  const fetchDoctors = async () => {
    try {
      const response = await apiClient.get("/doctors/");
      return response.data;
    } catch (err) {
      setToast({ message: "Data gagal diambil", type: "error" });
      console.error("Fetch error:", err);
      return [];
    }
  };

  const reloadDoctors = async () => {
    const data = await fetchDoctors();
    const items = Array.isArray(data) ? data : data?.data || [];
    setDoctorList(items);
  };

  useEffect(() => {
    reloadDoctors().finally(() => setLoading(false));
  }, []);

  const openAddModal = () => {
    setModalMode("add");
    setModalDoctor(null);
    setModalOpen(true);
  };

  const openEditModal = async (id) => {
    try {
      const doctor = await getDoctorById(id);
      setModalDoctor(doctor);
      setModalMode("edit");
      setModalOpen(true);
    } catch (err) {
      setToast({ message: "Dokter gagal diambil", type: "error" });
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
      await deleteDoctor(deleteTargetId);
      await reloadDoctors();
      setToast({ message: "Dokter berhasil dinonaktifkan", type: "success" });
    } catch (err) {
      alert("Gagal menonaktifkan dokter.");
    } finally {
      setIsConfirmOpen(false);
      setDeleteTargetId(null);
    }
  };

  const handleModalSuccess = async () => {
    await reloadDoctors();
    setToast({
      message:
        modalMode === "edit"
          ? "Data berhasil diperbarui"
          : "Data berhasil ditambahkan",
      type: "success",
    });
    setModalOpen(false);
    setModalDoctor(null);
  };

  const { searchTerm, setSearchTerm, filteredData } = useSearch(doctorList, [
    "full_name",
    "specialty",
    "email",
  ]);

  const columns = [
    { header: "Nama Dokter", accessor: "full_name" },
    { header: "Spesialisasi", accessor: "specialization" },
    { header: "Nomor STR", accessor: "str_number" },
    { header: "No. Hp", accessor: "phone_number" },
    { header: "Alamat Praktik", accessor: "practice_address" },
    { header: "Email", accessor: "email" },
    {
      header: "Status",
      accessor: "status",
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
            { label: "Non-Aktifkan", onClick: () => handleDeleteRequest(item.id) },
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
        <h1 className="text-2xl font-bold mb-6">Daftar Dokter</h1>
        <SearchBar value={searchTerm} onChange={setSearchTerm} />

        <div className="border-1 rounded-md border-gray-300 bg-white p-5">
          <div className="flex gap-4 mb-4">
            <button
              onClick={openAddModal}
              className="flex items-center text-blue-700 font-semibold space-x-1 bg-transparent border border-blue-700 py-2 px-4 rounded-md"
            >
              <PlusIcon className="w-4 h-4" />
              <span>Tambah Dokter</span>
            </button>
          </div>

          <div className="max-w-[1121px]">
            <DataTable columns={columns} data={filteredData} showIndex={true} />
          </div>
        </div>
      </div>

      <DoctorModal
        isOpen={modalOpen}
        close={() => setModalOpen(false)}
        onSuccess={handleModalSuccess}
        mode={modalMode}
        doctor={modalDoctor}
      />

      <ConfirmDialog
        isOpen={isConfirmOpen}
        title="Konfirmasi Penghapusan"
        description="Apakah Anda yakin ingin menghapus dokter ini?"
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

export default AturDoctorPage;
