import React, { useState, useEffect } from "react";
import Modal from "../../components/modal/modal";
import InputField from "../../components/inputField";
import Button from "../../components/buttonComp";
import { apiClient } from "../../config/api";
import Toast from "../../components/toast";
import { getFriendlyErrorMessage } from "../../utils/errorHandler";
import Select from "../../components/SelectComp";

const fields = [
  { accessor: "full_name", label: "Nama Karyawan" },
  { accessor: "nip", label: "NIP", type:"phone" },
  { accessor: "phone", label: "No. Hp", type: "phone" },
  { accessor: "email", label: "Email" },
  { accessor: "role", label: "Role", type:"select", options:["Admin", "Pegawai"] },
  { accessor: "password", label: "Password", type: "password" },
  { accessor: "reset_password", label: "Reset Password", type: "password" },
];

export default function KaryawanModal({ isOpen, close, onSuccess, mode = "add", karyawan = null }) {
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
      if (mode === "edit" && karyawan) {
        const base = generateInitialFormState();
        setForm({ ...base, ...karyawan });
      } else {
        setForm(generateInitialFormState());
      }
    }
  }, [isOpen, mode, karyawan]);

  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

const handleSubmit = async () => {
  const requiredFields = ["full_name", "email", "phone", "role", "nip"];
  const allFilled = requiredFields.every((key) => form[key]?.toString().trim() !== "");
  if (!allFilled) {
    setToast({ message: "Semua kolom wajib diisi.", type: "error" });
    return;
  }

  const payload = {
    ...form,
    active: true,
  };

  setLoading(true);
  try {
    if (mode === "edit" && karyawan?.id) {
      await apiClient.put(`/users/${karyawan.id}`, payload);
      if (form.reset_password?.trim()) {
        await apiClient.put(`/users/${karyawan.id}/reset-password`, {
          new_password: form.reset_password,
        });
      }

      setToast({ message: "Karyawan berhasil diperbarui!", type: "success" });
    } else {
      await apiClient.post("/users/register", payload);
      setToast({ message: "Karyawan berhasil ditambahkan!", type: "success" });
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
        {mode === "edit" ? "Edit Karyawan" : "Tambah Karyawan"}
      </h2>

      <div className="max-h-[60vh] overflow-y-auto pr-2 px-5">
       <div className="grid grid-cols-1 md:grid-cols-2 gap-4 w-full">
        {fields.map(({ accessor, label, type, options }) => {
          if (accessor === "password" && mode === "edit") return null;
          if (accessor === "reset_password" && mode !== "edit") return null;

          if (type === "select") {
            return (
              <div key={accessor} className="flex flex-col mb-4">
                <label className="text-sm font-medium mb-1">{label}</label>
                <Select
                  value={form[accessor]}
                  onChange={handleChange(accessor)}
                >
                  <option value="">Pilih {label}</option>
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
        <Button onClick={handleSubmit} disabled={loading} className="w-full">
          {loading ? "Menyimpan..." : mode === "edit" ? "Update" : "Simpan"}
        </Button>
        <button
          onClick={() => setForm(generateInitialFormState())}
          className="w-full bg-gray-200 border border-black text-black rounded-md py-2 hover:bg-gray-300 transition"
        >
          Reset
        </button>
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
