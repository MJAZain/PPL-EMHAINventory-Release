import React, { useState, useEffect } from "react";
import Modal from "../../components/modal/modal";
import InputField from "../../components/inputField";
import Button from "../../components/buttonComp";
import { apiClient } from "../../config/api";
import Toast from "../../components/toast";
import { getFriendlyErrorMessage } from "../../utils/errorHandler";

const formFields = [
  { label: "Nama Kategori", key: "name", placeholder: "Nama kategori" },
  { label: "Deskripsi", key: "description", placeholder: "Deskripsi" },
];

export default function KategoriModal({
  isOpen,
  close,
  onSuccess,
  mode = "add",
  kategori = null,
}) {
  const [loading, setLoading] = useState(false);
  const [toast, setToast] = useState(null);
  const [form, setForm] = useState({});

  const generateInitialFormState = () => {
    const state = {};
    formFields.forEach(({ key }) => {
      state[key] = "";
    });
    return state;
  };

  useEffect(() => {
    if (mode === "edit" && kategori) {
      const base = generateInitialFormState();
      const populated = { ...base, ...kategori };

      formFields.forEach(({ key }) => {
        if (typeof populated[key] === "object" && populated[key] !== null) {
          populated[key] = populated[key].id;
        }
      });

      setForm(populated);
    } else {
      setForm(generateInitialFormState());
    }
  }, [mode, kategori]);

  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

  const handleSubmit = async () => {
    const allFilled = Object.values(form).every(
      (val) => val?.toString().trim() !== ""
    );
    if (!allFilled) {
      setToast({ message: "Semua kolom harus diisi.", type: "error" });
      return;
    }

    const payload = { ...form };
    setLoading(true);

    try {
      if (mode === "edit" && kategori?.id) {
        await apiClient.put(`/categories/${kategori.id}`, payload);
        setToast({ message: "Kategori berhasil diperbarui!", type: "success" });
      } else {
        await apiClient.post("/categories/", payload);
        setToast({ message: "Kategori berhasil ditambahkan!", type: "success" });
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
        {mode === "edit" ? "Edit Kategori" : "Tambah Kategori"}
      </h2>

      <div className="max-h-[60vh] overflow-y-auto pr-2 px-5">
        <div className="gap-4 w-full">
          {formFields.map(({ label, key, placeholder, type }) => {
            const isTextArea = key.includes("description");
            if (isTextArea) {
              return (
                <div key={key} className="flex flex-col col-span-full">
                  <label className="text-sm font-medium mb-1">{label}</label>
                  <textarea
                    value={form[key]}
                    onChange={handleChange(key)}
                    placeholder={placeholder}
                    rows={4}
                    className="border border-gray-300 rounded-md px-3 py-2 resize-none"
                  />
                </div>
              );
            }

            return (
              <InputField
                key={key}
                label={label}
                value={form[key]}
                onChange={handleChange(key)}
                placeholder={placeholder}
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
