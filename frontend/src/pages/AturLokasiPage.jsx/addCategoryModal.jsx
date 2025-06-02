import React, { useState } from "react";
import Modal from "./modal";
import InputField from "../inputField";
import Button from "../buttonComp";
import { apiClient } from "../../config/api";
import Toast from "../toast";

const initialFormState = {
  name: "",
  description: "",
};

export default function AddCategoryModal({ isOpen, close, onSuccess }) {
  const [form, setForm] = useState(initialFormState);
  const [loading, setLoading] = useState(false);
  const [toast, setToast] = useState(null);


  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

  const handleSubmit = async () => {
    const allFilled = Object.values(form).every((val) => val.trim() !== "");
    if (!allFilled) {
      setToast({
          message:"Semua kolom harus di isi.",
          type: "error",
        });
      return;
    }

    setLoading(true);
    try {
        await apiClient.post("/categories/", form);
        setToast({
          message: "Kategori berhasil ditambahkan!",
          type: "success",
        });
        onSuccess();
        setForm(initialFormState);
      } catch (err) {
        setToast({
          message: err.message || "Gagal menambahkan kategori.",
          type: "error",
        });
      } finally {
        setLoading(false);
      }
  };

  return (
    <Modal isOpen={isOpen} close={close} contentClassName="w-96">
      <h2 className="text-xl font-semibold mb-4 text-center py-5">Tambah Kategori</h2>
      <div className="gap-5">
        <InputField
          label="Nama"
          value={form.name}
          onChange={handleChange("name")}
          placeholder="Masukkan nama kategori"
          className="w-full h-10"
        />
        <InputField
          label="Deskripsi"
          value={form.description}
          onChange={handleChange("description")}
          placeholder="Masukkan deskripsi"
          className="w-full h-20"
        />

        <Button onClick={handleSubmit} disabled={loading} className="w-full">
          {loading ? "Menyimpan..." : "Simpan"}
        </Button>

      </div>
      <div>
        {toast && (
          <Toast
            message={toast.message}
            type={toast.type}
            onClose={() => setToast(null)}
          />
        )}
      </div>
    </Modal>
  );
}
