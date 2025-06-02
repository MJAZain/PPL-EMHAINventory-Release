import React, { useState } from "react";
import InputField from "../components/inputField";
import Button from "../components/buttonComp";
import useRegisterUser from "../hooks/useRegisterUser";

const initialFormState = {
  email: "",
  phone: "",
  full_name: "",
  role: "pegawai",
  nip: "",
};

export default function RegisterUserPage() {
  const [form, setForm] = useState(initialFormState);
  const { registerUser, loading } = useRegisterUser();

  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

  const handleSubmit = async () => {
    const allFilled = Object.values(form).every((val) => val.trim() !== "");
    if (!allFilled) {
      alert("Mohon isi semua field.");
      return;
    }

    try {
      await registerUser(form);
      alert("User berhasil didaftarkan!");
      setForm(initialFormState);
    } catch {
      alert("Gagal mendaftarkan user.");
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center">
      <div className="bg-white p-8 rounded-xl shadow-md w-full max-w-md">
        <h2 className="text-2xl font-bold text-center mb-6">Daftar User Baru</h2>

        <div className="space-y-4">
          <InputField
            label="Email"
            value={form.email}
            onChange={handleChange("email")}
            placeholder="Masukkan email"
            type="email"
          />
          <InputField
            label="No. Telepon"
            value={form.phone}
            onChange={handleChange("phone")}
            placeholder="Masukkan no. telepon"
          />
          <InputField
            label="Nama Lengkap"
            value={form.full_name}
            onChange={handleChange("full_name")}
            placeholder="Masukkan nama lengkap"
          />
          <div>
            <label className="block mb-1 font-medium">Role</label>
            <select
              value={form.role}
              onChange={handleChange("role")}
              className="w-full border border-gray-300 rounded-lg px-3 py-2"
            >
              <option value="admin">Admin</option>
              <option value="pegawai">Pegawai</option>
            </select>
          </div>
          <InputField
            label="NIP"
            value={form.nip}
            onChange={handleChange("nip")}
            placeholder="Masukkan NIP"
          />
        </div>

        <div className="mt-6">
          <Button onClick={handleSubmit} className="w-full" disabled={loading}>
            {loading ? "Mendaftarkan..." : "Daftarkan User"}
          </Button>
        </div>
      </div>
    </div>
  );
}
