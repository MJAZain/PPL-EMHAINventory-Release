import React, { useEffect, useState } from "react";
import Modal from "./modal";
import InputField from "../inputField";
import Button from "../buttonComp";
import useUserActions from "../../hooks/useUserAction";
import { apiClient } from "../../config/api";
import Toast from "../toast";
import { getFriendlyErrorMessage } from "../../utils/errorHandler";

const initialFormState = {
  full_name: "",
  role: "",
  phone:"",
  email:"",
  active:true,
  nip:"",
  password:"",
  confirmPassword: "",
};

const labelMap = {
  full_name: "Nama",
  role: "Peran",
  phone: "Telepon",
  email: "Email",
  active: "Status",
  nip: "NIP",
  password: "Password",
  confirmPassword: "Konfirmasi Password",
};

export default function EditUserModal({ isOpen, close, userId, onSuccess }) {
  const [form, setForm] = useState(initialFormState);
  const { getUserById } = useUserActions();
  const [loading, setLoading] = useState(false);
  const [toast, setToast] = useState(null);

  useEffect(() => {
    console.log("Fetching user", userId, "Modal open:", isOpen);
    if (userId && isOpen) {
      (async () => {
        try {
          const data = await getUserById(userId);
          setForm({
             ...initialFormState,
            full_name: data.full_name || "",
            email: data.email || "",
            phone: data.phone || "",
            nip: data.nip || "",
            role: data.role || "",
          });
        } catch {
          setToast({
            message: "Gagal mengambil data user.",
            type: "error",
          });
        }
      })();
    }
  }, [userId, isOpen, getUserById]);

  const handleChange = (key) => (e) => {
    setForm({
      ...form,
      [key]: key === "phone" || key === "nip" ? e.target.value.replace(/\D/, "") : e.target.value,
    });
  };

  const handleSubmit = async () => {
    const { full_name, email, phone, nip, role } = form;

    if (![full_name, email, phone, nip, role].every((val) => val && val.trim() !== "")) {
      setToast({ message: "Semua kolom harus diisi.", type: "error" });
      return;
    }

    setLoading(true);
    try {
      await apiClient.put(`/users/${userId}`, form);
      onSuccess?.();
      close();
    } catch (err) {
      const message = getFriendlyErrorMessage(err);
      setToast({ message, type: "error" });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal isOpen={isOpen} close={close} contentClassName="w-96">
      <h2 className="text-xl font-semibold mb-4 text-center py-5">Tambah User</h2>
      <div className="gap-5 grid grid-cols-2">
        {Object.keys(initialFormState).map((key) => {
          if (key === "active") return null;
          if (key === "role") {
            return (
              <div key={key} className="flex flex-col">
                <label className="mb-1 font-medium">{labelMap[key]}</label>
                <select
                  value={form[key]}
                  onChange={(e) => setForm({ ...form, [key]: e.target.value })}
                  className="w-full h-10 border rounded px-2"
                >
                  <option value="">Pilih peran</option>
                  <option value="Admin">Admin</option>
                  <option value="Pegawai">Pegawai</option>
                </select>
              </div>
            );
          }

          return (
            <InputField
              key={key}
              label={labelMap[key] || key.charAt(0).toUpperCase() + key.slice(1)}
              type={key.toLowerCase().includes("password") ? "password" : "text"}
              value={form[key]}
              onChange={handleChange(key)}
              placeholder={`Masukkan ${labelMap[key] || key}`}
              className="w-full h-10 border rounded px-2"
            />
          );
        })}
       
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
