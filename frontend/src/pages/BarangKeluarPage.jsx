import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import Sidebar from "../components/Sidebar";
import Button from "../components/buttonComp";
import InputField from "../components/inputField";
import Toast from "../components/toast";


const initialFormState = {
  date: "",
  customer: "",
  no_faktur: "",
};

const labelMap = {
  date: "Tanggal",
  customer: "Nama Customer",
  no_faktur: "NSFP",
};

export default function BarangKeluarPage() {
  const [toast, setToast] = useState(null);

  const navigate = useNavigate();
  const [form, setForm] = useState(initialFormState);

  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

  const handleNext = () => {
    const allFilled = Object.values(form).every((val) => val.trim() !== "");
    if (!allFilled) {
      setToast({
        message: "Mohon isi semua field.",
        type: "error",
      });
      return;
    }
    localStorage.setItem("barangKeluarForm", JSON.stringify(form));
    navigate("/barang-keluar-detail");
  };


  return (
    <div className="flex min-h-screen bg-gray-100">

      <div className="bg-white min-h-screen">
        <Sidebar />
      </div>

      <div className="bg-white max-w-xl mx-auto w-full p-6 mt-10 border rounded-md border-gray-300 shadow-md">
        <h1 className="text-2xl font-bold text-center mb-6">Barang Terjual</h1>

        {/* Dynamically render input fields */}
        <div className="space-y-4">
          {Object.keys(form).map((key) => (
            <InputField
              key={key}
              label={labelMap[key] || key}
              type={key === "date" ? "date" : "text"}
              value={form[key]}
              onChange={handleChange(key)}
              placeholder={key === "date" ? undefined : `Masukkan ${labelMap[key] || key}`}
              className="w-full h-10"
            />
          ))}

        </div>
          {toast && (
            <Toast
              message={toast.message}
              type={toast.type}
              onClose={() => setToast(null)}
            />
          )}
        <div className="mt-6">
          <Button className="w-full" onClick={handleNext}>
            Selanjutnya
          </Button>
        </div>
      </div>
    </div>
  );
}
