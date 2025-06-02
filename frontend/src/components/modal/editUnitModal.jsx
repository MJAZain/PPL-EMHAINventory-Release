import React, { useEffect, useState } from "react";
import Modal from "./modal";
import InputField from "../inputField";
import Button from "../buttonComp";
import useUnitActions from "../../hooks/useUnitAction";
import { apiClient } from "../../config/api";
import Toast from "../toast";

export default function EditUnitModal({ isOpen, close, unitId, onSuccess }) {
  const [form, setForm] = useState({ name: "", description: "" });
  const { getUnitById } = useUnitActions();
  const [loading, setLoading] = useState(false);
  const [toast, setToast] = useState(null);

  useEffect(() => {
    if (unitId && isOpen) {
      (async () => {
        try {
          const data = await getUnitById(unitId);
          setForm({ name: data.name, description: data.description || "" });
        } catch {
          setToast({
          message: "Gagal menambahkan kategori.",
          type: "error",
        });
        }
      })();
    }
  }, [unitId, isOpen, getUnitById]);

  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

  const handleSubmit = async () => {
  setLoading(true);

  try {
    const res = await apiClient.put(`/units/${unitId}`, form);
    onSuccess();
  } catch (err) {
    setToast({
      message: err.message || "Gagal mengedit satuan.",
      type: "error",
    });
  } finally {
    setLoading(false);
  }
};


  return (
    <Modal isOpen={isOpen} close={close} contentClassName="w-96">
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
      <div className="py-5 pr-2">
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
