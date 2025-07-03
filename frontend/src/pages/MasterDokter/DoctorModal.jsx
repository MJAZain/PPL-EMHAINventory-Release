import React, { useState, useEffect } from "react";
import Modal from "../../components/modal/modal";
import InputField from "../../components/inputField";
import Button from "../../components/buttonComp";
import { apiClient } from "../../config/api";
import Toast from "../../components/toast";
import { getFriendlyErrorMessage } from "../../utils/errorHandler";
import Select from '../../components/SelectComp'

const fields = [
  { accessor: "full_name", label: "Nama Dokter" },
  { accessor: "specialization", label: "Spesialisasi" },
  { accessor: "str_number", label: "Nomor STR" },
  { accessor: "phone_number", label: "No. Hp", type: "phone"},
  { accessor: "practice_address", label: "Alamat Praktik" },
  { accessor: "email", label: "Email" },
  { accessor: "status", label: "Status", type: "select", options: ["Aktif", "Nonaktif"] },
];

export default function DoctorModal({ isOpen, close, onSuccess, mode = "add", doctor = null }) {
  const [loading, setLoading] = useState(false);
  const [toast, setToast] = useState(null);
  const [form, setForm] = useState({});

  const generateInitialFormState = () => {
    const state = {};
    fields.forEach(({ accessor }) => {
      state[accessor] = "";
    });
    return state;
  };

  useEffect(() => {
    if (isOpen) {
      setToast(null);
      if (mode === "edit" && doctor) {
        const base = generateInitialFormState();
        setForm({ ...base, ...doctor });
      } else {
        setForm(generateInitialFormState());
      }
    }
  }, [isOpen, mode, doctor]);


  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

  const handleSubmit = async () => {
    const allFilled = Object.values(form).every((val) => val?.toString().trim() !== "");
    if (!allFilled) {
      setToast({ message: "Semua kolom harus diisi.", type: "error" });
      return;
    }

    setLoading(true);
    try {
      if (mode === "edit" && doctor?.id) {
        await apiClient.put(`/doctors/${doctor.id}`, form);
        setToast({ message: "Dokter berhasil diperbarui!", type: "success" });
      } else {
        await apiClient.post("/doctors/", form);
        setToast({ message: "Dokter berhasil ditambahkan!", type: "success" });
      }

      onSuccess();
      close();
    } catch (err) {
      const message = getFriendlyErrorMessage(err);
      setToast({ message, type: "error" });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal isOpen={isOpen} close={close} contentClassName="w-full max-w-2xl">
      <h2 className="text-xl font-semibold mb-4 text-center py-5">
        {mode === "edit" ? "Edit Dokter" : "Tambah Dokter"}
      </h2>

      <div className="max-h-[60vh] overflow-y-auto pr-2 px-5">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 w-full">
          {fields.map(({ accessor, label, type, options }) => {
            if (type === "select") {
              return (
                <div key={accessor} className="flex flex-col mb-4">
                  <label className="text-sm font-medium mb-1">{label}</label>
                  <Select
                    value={form[accessor]}
                    onChange={handleChange(accessor)}
                  >
                    <option value="">Pilih Status</option>
                    {options.map((opt) => (
                      <option key={opt} value={opt}>
                        {opt}
                      </option>
                    ))}
                  </Select>
                </div>
              );
            }

            return (
              <InputField
                key={accessor}
                label={label}
                value={form[accessor]}
                onChange={handleChange(accessor)}
                placeholder={label}
                type={type || "text"}
              />
            );
          })}
        </div>
      </div>

      <div className="mt-6 flex justify-between gap-4">
        <button
          onClick={() => setForm(generateInitialFormState())}
          className="w-full bg-gray-200 border border-black text-black rounded-md py-2 hover:bg-gray-300 transition"
        >
          Reset
        </button>
        <Button onClick={handleSubmit} disabled={loading} className="w-full">
          {loading ? "Menyimpan..." : mode === "edit" ? "Update" : "Simpan"}
        </Button>
      </div>

      {toast && (
        <Toast
          message={toast.message}
          type={toast.type}
          onClose={() => setToast(null)}
        />
      )}
    </Modal>
  );
}
