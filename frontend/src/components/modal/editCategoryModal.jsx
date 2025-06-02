import React, { useEffect, useState } from "react";
import Modal from "./modal";
import InputField from "../inputField";
import Button from "../buttonComp";
import useCategoryActions from "../../hooks/useCategoryAction";
import { apiClient } from "../../config/api";
import Toast from "../toast";

export default function EditCategoryModal({ isOpen, close, categoryId, onSuccess }) {
  const [form, setForm] = useState({ name: "", description: "" });
  const { getCategoryById } = useCategoryActions();
  const [loading, setLoading] = useState(false);
  const [toast, setToast] = useState(null);

  useEffect(() => {
    if (categoryId && isOpen) {
      (async () => {
        try {
          const data = await getCategoryById(categoryId);
          setForm({ name: data.name, description: data.description || ""});
        } catch {
          setToast({
          message: "Gagal mengambil data kategori.",
          type: "error",
        });
        }
      })();
    }
  }, [categoryId, isOpen, getCategoryById]);

  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

  const handleSubmit = async () => {
    setLoading(true);
    try {
        await apiClient.put(`/categories/${categoryId}`, form);
        setToast({
          message: "Kategori berhasil ditambahkan!",
          type: "success",
        });
        onSuccess();
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
    <Modal isOpen={isOpen} close={close}>
      <h2 className="text-xl font-semibold mb-4 text-center">Edit Satuan</h2>
      <div className="gap-5">
        <InputField
          label="Nama"
          value={form.name}
          onChange={handleChange("name")}
          placeholder="Masukkan nama satuan"
          className="w-full h-10"
        />
        <InputField
          label="Deskripsi"
          value={form.description}
          onChange={handleChange("description")}
          placeholder="Masukkan deskripsi"
          className="w-full h-20"
        />
      </div>
      
        <Button onClick={handleSubmit} disabled={loading} className="w-full">
          {loading ? "Menyimpan..." : "Simpan"}
        </Button>
      
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
