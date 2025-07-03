import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import Sidebar from "../../components/Sidebar";
import Button from "../../components/buttonComp";
import InputField from "../../components/inputField";
import Toast from "../../components/toast";
import { apiClient } from "../../config/api";

const formFields = [
  { label: "Tanggal Shift", key: "shift_date", placeholder: "Tanggal Shift", type: "date" },
  {
    label: "Petugas Pembuka",
    key: "opening_officer_id",
    placeholder: "Pilih Petugas Pembuka",
    type: "select",
    dynamicOptionsKey: "users",
  },
  { label: "Waktu Mulai Shift", key: "opening_time", placeholder: "Jam Mulai Shift", type: "time" },
  { label: "Saldo Kas Awal", key: "opening_balance", placeholder: "Saldo Kas Awal", type: "number" },
  { label: "Catatan", key: "notes", placeholder: "Catatan tambahan", type: "textarea" },
  {
    label: "Status",
    key: "status",
    placeholder: "Pilih Status",
    type: "select",
    options: [
      { id: "open", name: "Open" },
      { id: "closed", name: "Closed" },
    ],
  },
];

export default function TanpaResepShiftUmumPage() {
  const [toast, setToast] = useState(null);
  const [form, setForm] = useState({
    shift_date: "",
    opening_officer_id: "",
    opening_time: "",
    opening_balance: "",
    notes: "",
    status: "open",
  });
  const [users, setUsers] = useState([]);

  const navigate = useNavigate();

  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

const handleNext = async () => {
  const requiredFields = [
    "shift_date",
    "opening_officer_id",
    "opening_time",
    "opening_balance",
    "status",
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
      ...form,
      opening_balance: parseFloat(form.opening_balance),
    };

    const response = await apiClient.post("/shifts/open", payload);

    if (response.status === 200 || response.status === 201) {
      setToast({
        message: "Shift berhasil dibuka.",
        type: "success",
      });

      setTimeout(() => {
        navigate("/tanpa-resep-shift-detail");
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


  const fetchUsers = async () => {
    try {
      const response = await apiClient.get("/users/");
      return response.data;
    } catch (error) {
      console.error("Failed to fetch users", error);
      throw error;
    }
  };

  useEffect(() => {
    const loadUsers = async () => {
      try {
        const data = await fetchUsers();

        if (Array.isArray(data)) {
          setUsers(data);
        } else if (Array.isArray(data?.data)) {
          setUsers(data.data);
        } else {
          console.error("Unexpected data format:", data);
          setUsers([]);
        }

      } catch (error) {
        console.error("Error loading users", error);
        setUsers([]);
      }
    };

    loadUsers();
  }, []);

  return (
    <div className="flex min-h-screen bg-gray-100">
      <div className="bg-white min-h-screen">
        <Sidebar />
      </div>

      <div className="bg-white max-w-xl mx-auto w-full p-6 mt-10 border rounded-md border-gray-300 shadow-md max-h-[90vh] overflow-y-auto">
        <h1 className="text-2xl font-bold text-center mb-6">Buka Shift Kasir khusus Resep</h1>

        <div className="pr-2">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 w-full">
            {formFields.map(({ label, key, placeholder, type, options, dynamicOptionsKey }) => {
              if (type === "select") {
                let selectOptions = options || [];

                if (dynamicOptionsKey === "users") {
                  selectOptions = users.map((user) => ({
                    id: user.id,
                    name: user.name || user.full_name || `User ${user.id}`,
                  }));
                }

                return (
                  <div key={key} className="flex flex-col mb-4">
                    <label className="text-sm font-medium mb-1">{label}</label>
                    <select
                      disabled
                      value={form[key]}
                      onChange={handleChange(key)}
                      className="border border-gray-300 rounded-md px-3 py-2"
                    >
                      <option value="">{placeholder}</option>
                      {selectOptions.map((opt) => (
                        <option key={opt.id} value={opt.id}>
                          {opt.name}
                        </option>
                      ))}
                    </select>
                  </div>
                );
              }

              if (type === "textarea") {
                return (
                  <div key={key} className="flex flex-col mb-4">
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
