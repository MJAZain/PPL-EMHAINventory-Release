import React, { useState } from "react";
import Modal from "./modal";
import InputField from "../inputField";
import Button from "../buttonComp";
import { apiClient } from "../../config/api";
import Toast from "../toast";
import { getFriendlyErrorMessage } from "../../utils/errorHandler";
import Select from "../SelectComp";

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

export default function AddUserModal({ isOpen, close, onSuccess }) {
  const [form, setForm] = useState(initialFormState);
  const [loading, setLoading] = useState(false);
  const [toast, setToast] = useState(null);

  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

  const handleSubmit = async () => {
    const { password, confirmPassword, ...rest } = form;

    const allFilled = Object.values(form).every(val =>
      typeof val === "string" ? val.trim() !== "" : val !== null && val !== undefined
    );


    if (!allFilled) {
      setToast({ message: "Semua kolom harus diisi.", type: "error" });
      return;
    }

    if (password !== confirmPassword) {
      setToast({ message: "Password dan konfirmasi tidak cocok.", type: "error" });
      return;
    }

    const payload = {
      full_name: rest.full_name?.trim(),
      email: rest.email?.trim(),
      phone: rest.phone?.trim(),
      nip: rest.nip?.trim(),
      role: rest.role?.toLowerCase(),
      password: password.trim(),
      active: true,
    };

    try {
      console.log("Payload:", payload);

      await apiClient.post("/users/register", payload);
      onSuccess?.();
      setForm(initialFormState);
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
                <Select
                  value={form[key]}
                  onChange={(e) => setForm({ ...form, [key]: e.target.value })}
                  className="w-full h-10 border rounded px-2"
                >
                  <option value="">Pilih peran</option>
                  <option value="Admin">Admin</option>
                  <option value="Pegawai">Pegawai</option>
                </Select>
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
