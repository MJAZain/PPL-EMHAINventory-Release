import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import Sidebar from "../../../components/Sidebar";
import Button from "../../../components/buttonComp";
import InputField from "../../../components/inputField";
import Toast from "../../../components/toast";
import { apiClient } from "../../../config/api";

const formFields = [
  { label: "Tanggal Shift", key: "shift_date", placeholder: "Tanggal Shift", type: "date" },
  { label: "Waktu Mulai Shift", key: "opening_time", placeholder: "Jam Mulai Shift", type: "time" },
  {
    label: "Petugas Pembuka",
    key: "opening_officer",
  },
  { label: "Saldo Kas Awal", key: "opening_balance", placeholder: "Saldo Kas Awal", type: "number" },
];

export default function ResepShiftUmumPage() {
  const [toast, setToast] = useState(null);
  const now = new Date();
  const hh = String(now.getHours()).padStart(2, '0');
  const mm = String(now.getMinutes()).padStart(2, '0');
  const storedUser = JSON.parse(localStorage.getItem("user") || "{}");
  const [form, setForm] = useState({
    shift_date: now.toISOString().split("T")[0],
    opening_time: `${hh}:${mm}`, 
    opening_officer: storedUser.full_name || "",
    opening_balance: "",
  });

  const navigate = useNavigate();

  useEffect(() => {
  
  }, []);

  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

const handleNext = async () => {
  const requiredFields = [
    "opening_balance",
  ];

  const allFilled = requiredFields.every(
    (field) => String(form[field] ?? "").trim() !== ""
  );

  if (!allFilled) {
    setToast({
      message: "Mohon isi semua field yang wajib diisi.",
      type: "error",
    });
    return;
  }

  try {
    const payload = {
      opening_balance: parseFloat(form.opening_balance),
    };

    const response = await apiClient.post("/shifts/open", payload);

    if (response.status === 200 || response.status === 201) {

       const presId = response.data?.data?.id;

      if (presId) {
        localStorage.setItem("presId", presId);
      } else {
        console.warn("Shift ID not found in response!");
      }
      
      console.log("presId:", presId);
      localStorage.setItem('closing_balance_pres', form.opening_balance);

      setToast({
        message: "Shift berhasil dibuka.",
        type: "success",
      });

      setTimeout(() => {
        navigate("/resep-detail");
      }, 1500); 
    } else {
      throw new Error("Server tidak merespons seperti yang diharapkan.");
    }
  } catch (error) {
    console.error("Error posting shift:", error);
    setToast({
      message: "Gagal membuka shift. Silakan coba lagi.",
      type: "error",
    });
  }
};

  return (
    <div className="flex min-h-screen bg-gray-100">
      <div className="bg-white min-h-screen">
        <Sidebar />
      </div>

      <div className="bg-white max-w-xl mx-auto w-full p-6 mt-10 border rounded-md border-gray-300 shadow-md max-h-[90vh] overflow-y-auto">
        <h1 className="text-2xl font-bold text-center mb-6">Buka Shift Kasir Resep</h1>

        <div className="pr-2">
          <div className="grid grid-cols-1 md:grid-cols-1 gap-4 w-full">
            {formFields.map(({ label, key, placeholder, type }) => {
              const isDisabled = key !== "opening_balance";

              return (
                <InputField
                  key={key}
                  label={label}
                  value={form[key]}
                  onChange={handleChange(key)}
                  placeholder={placeholder}
                  type={type || "text"}
                  disabled={isDisabled}
                />
              );
            })}
          </div>
        </div>

        <div className="mt-6">
          <Button className="w-full" onClick={handleNext}>
            Selanjutnya
          </Button>
        </div>

        {toast && (
          <Toast
            message={toast.message}
            type={toast.type}
            onClose={() => setToast(null)}
          />
        )}
      </div>
    </div>
  );
}
