import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import Sidebar from "../../components/Sidebar";
import Button from "../../components/buttonComp";
import InputField from "../../components/inputField";
import Toast from "../../components/toast";
import { apiClient } from "../../config/api";
import TextArea from '../../components/textareacomp';

const formFields = [
  { label: "Tanggal Opname", key: "opname_date", placeholder: "Tanggal Opname", type: "date", required: true },
  { label: "Catatan", key: "notes", placeholder: "Catatan (opsional)", type: "text", required: false },
];

export default function CreateDraftPage() {
  const [toast, setToast] = useState(null);
  const now = new Date();
  const [form, setForm] = useState({
    opname_date: now.toISOString().split("T")[0],
    notes: "",
  });

  const navigate = useNavigate();

  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

  const handleSubmit = async () => {
    // Validate required field
    if (!form.opname_date) {
      setToast({
        message: "Tanggal opname wajib diisi",
        type: "error",
      });
      return;
    }

    try {
      const payload = {
        opname_date: form.opname_date,
        notes: form.notes || undefined,
      };

      const response = await apiClient.post("/stock-opname/draft", payload);

      if (response.status === 200 || response.status === 201) {
        const opnameId = response.data?.data?.opname_id;

        if (opnameId) {
          localStorage.setItem("opnameId", opnameId);
        } else {
          console.warn("Opname ID not found in response!");
        }

        setToast({
          message: "Draft opname berhasil dibuat",
          type: "success",
        });

        setTimeout(() => {
          navigate("/draft-detail");
        }, 1500); 
        
      } else {
        throw new Error("Server tidak merespons seperti yang diharapkan");
      }
    } catch (error) {
      console.error("Error creating draft:", error);
      setToast({
        message: error.response?.data?.error || "Gagal membuat draft. Silakan coba lagi",
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
        <h1 className="text-2xl font-bold text-center mb-6">Buat Draft Opname</h1>

        <div className="pr-2">
          <div className="grid grid-cols-1 md:grid-cols-1 gap-4 w-full">
            {formFields.map(({ label, key, placeholder, type, required }) => (
              key === "notes" ? (
                <TextArea
                  key={key}
                  label={label}
                  value={form[key]}
                  onChange={handleChange(key)}
                  placeholder={placeholder}
                  rows={3}
                />
              ) : (
                <InputField
                  key={key}
                  label={label}
                  value={form[key]}
                  onChange={handleChange(key)}
                  placeholder={placeholder}
                  type={type || "text"}
                  required={required}
                />
              )
            ))}
          </div>
        </div>

        <div className="mt-6">
          <Button className="w-full" onClick={handleSubmit}>
            Buat Draft
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