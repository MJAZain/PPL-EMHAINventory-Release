import React, { useState, useEffect } from "react";
import useSearch from "../../hooks/useSearch";
import usePatientActions from "./usePatientAction";
import { useNavigate } from "react-router-dom";

import { apiClient } from "../../config/api";
import { PlusIcon } from "@heroicons/react/24/solid";

import ActionMenu from "../../components/ActionMenu";
import SearchBar from "../../components/SearchBar";
import DataTable from "../../components/tableCompo";

import PatientModal from "./PatientModal";
import Sidebar from "../../components/Sidebar";
import ConfirmDialog from "../../components/ConfirmDialog";
import Toast from "../../components/toast";

function AturPatientPage() {
  const [patientList, setPatientList] = useState([]);
  const [loading, setLoading] = useState(true);
  const [toast, setToast] = useState(null);

  const [modalOpen, setModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState("add");
  const [modalPatient, setModalPatient] = useState(null);

  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [deleteTargetId, setDeleteTargetId] = useState(null);

  const navigate = useNavigate();
  const { getPatientById, deletePatient } = usePatientActions();

  const fetchPatients = async () => {
    try {
      const response = await apiClient.get("/patients/");
      return response.data.data;
    } catch (err) {
      setToast({ message: "Data gagal diambil", type: "error" });
      console.error("Fetch error:", err);
      return [];
    }
  };

  const reloadPatients = async () => {
    const data = await fetchPatients();
    const items = Array.isArray(data) ? data : data?.data || [];
    setPatientList(items);
  };

  useEffect(() => {
    reloadPatients().finally(() => setLoading(false));
  }, []);

  const openAddModal = () => {
    setModalMode("add");
    setModalPatient(null);
    setModalOpen(true);
  };

  const openEditModal = async (id) => {
    try {
      const patient = await getPatientById(id);
      setModalPatient(patient);
      setModalMode("edit");
      setModalOpen(true);
    } catch (err) {
      setToast({ message: "Pasien gagal diambil", type: "error" });
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
      await deletePatient(deleteTargetId);
      await reloadPatients();
      setToast({ message: "Data berhasil dihapus", type: "success" });
    } catch (err) {
      alert("Gagal menghapus data.");
    } finally {
      setIsConfirmOpen(false);
      setDeleteTargetId(null);
    }
  };

  const handleModalSuccess = async () => {
    await reloadPatients();
    setToast({
      message:
        modalMode === "edit"
          ? "Data berhasil diperbarui"
          : "Data berhasil ditambahkan",
      type: "success",
    });
    setModalOpen(false);
    setModalPatient(null);
  };

  const { searchTerm, setSearchTerm, filteredData } = useSearch(patientList, [
    "full_name",
    "gender",
    "identity_number",
    "patient_type",
    "phone_number",
  ]);

  const columns = [
    { header: "Nama Pasien", accessor: "full_name" },
    { header: "Jenis Kelamin", accessor: "gender" },
    { header: "Tempat Lahir", accessor: "place_of_birth" },
    { header: "Tanggal Lahir", accessor: "date_of_birth" },
    { header: "Alamat", accessor: "address" },
    { header: "No. Hp", accessor: "phone_number" },
    { header: "Jenis Pasien", accessor: "patient_type" },
    { header: "NIK", accessor: "identity_number" },
    { header: "Nama Penjamin", accessor: "guarantor_name" },
    { header: "Status", accessor: "status",
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
        <h1 className="text-2xl font-bold mb-6">Daftar Pasien</h1>
        <SearchBar value={searchTerm} onChange={setSearchTerm} />

        <div className="border-1 rounded-md border-gray-300 bg-white p-5">
          <div className="flex gap-4 mb-4">
            <button
              onClick={openAddModal}
              className="flex items-center text-blue-700 font-semibold space-x-1 bg-transparent border border-blue-700 py-2 px-4 rounded-md"
            >
              <PlusIcon className="w-4 h-4" />
              <span>Tambah Pasien</span>
            </button>
          </div>

          <div className="max-w-[1121px]">
            <DataTable columns={columns} data={filteredData} showIndex={true} />
          </div>
        </div>
      </div>

      <PatientModal
        isOpen={modalOpen}
        close={() => setModalOpen(false)}
        onSuccess={handleModalSuccess}
        mode={modalMode}
        patient={modalPatient}
      />

      <ConfirmDialog
        isOpen={isConfirmOpen}
        title="Konfirmasi Penghapusan"
        description="Apakah Anda yakin ingin menghapus pasien ini?"
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

export default AturPatientPage;
