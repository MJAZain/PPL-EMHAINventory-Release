import React, { useState, useEffect } from "react";
import Modal from "../../components/modal/modal";
import InputField from "../../components/inputField";
import Button from "../../components/buttonComp";
import { apiClient } from "../../config/api";
import Toast from "../../components/toast";
import { getFriendlyErrorMessage } from "../../utils/errorHandler";
import TextArea from "../../components/textareacomp";

const formFields = [
  { label: "Nama Satuan", key: "name", placeholder: "Nama satuan" },
  { label: "Deskripsi", key: "description", placeholder: "Deskripsi" },
];

export default function SatuanModal({
  isOpen,
  close,
  onSuccess,
  mode = "add",
  satuan = null,
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
    if (isOpen) {
      if (mode === "edit" && satuan) {
      const base = generateInitialFormState();
      const populated = { ...base, ...satuan };

      formFields.forEach(({ key }) => {
        if (typeof populated[key] === "object" && populated[key] !== null) {
          populated[key] = populated[key].id;
        }
      });

      setForm(populated);
    } else {
      setForm(generateInitialFormState());
    }
    }
    
  }, [isOpen, mode, satuan]);

  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

  const handleSubmit = async () => {
    const allFilled = Object.values(form).every(
      (val) => val?.toString().trim() !== ""
    );

    const payload = { ...form };
    setLoading(true);

    try {
      if (mode === "edit" && satuan?.id) {
        await apiClient.put(`/units/${satuan.id}`, payload);
        setToast({ message: "Satuan berhasil diperbarui!", type: "success" });
      } else {
        await apiClient.post("/units/", payload);
        setToast({ message: "Satuan berhasil ditambahkan!", type: "success" });
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
        {mode === "edit" ? "Edit Satuan" : "Tambah Satuan"}
      </h2>

      <div className="max-h-[60vh] overflow-y-auto pr-2 px-5">
        <div className="gap-4 w-full">
          {formFields.map(({ label, key, placeholder, type }) => {
            const isTextArea = key.includes("description");
            if (isTextArea) {
              return (
                <div key={key} className="flex flex-col col-span-full">
                  <label className="text-sm font-medium mb-1">{label}</label>
                  <TextArea
                    value={form[key]}
                    onChange={handleChange(key)}
                    placeholder={placeholder}
                    rows={4}
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
