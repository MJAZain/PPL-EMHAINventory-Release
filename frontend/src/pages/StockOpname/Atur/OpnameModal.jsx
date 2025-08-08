import React, { useState, useEffect } from "react";
import Modal from "../../components/modal/modal";
import InputField from "../../components/inputField";
import Button from "../../components/buttonComp";
import { apiClient } from "../../config/api";
import Toast from "../../components/toast";
import { getFriendlyErrorMessage } from "../../utils/errorHandler";

const fields = [
  { accessor: "opname_date", label: "Tanggal", type: "date" },
  { accessor: "notes", label: "Catatan", type: "text" },
];

export default function OpnameModal({ isOpen, close, onSuccess, mode = "add", opname = null }) {
  const [loading, setLoading] = useState(false);
  const [toast, setToast] = useState(null);
  const [form, setForm] = useState({});
  const [initialForm, setInitialForm] = useState({});


  const generateInitialFormState = () => {
    const state = {};
    fields.forEach(({ accessor }) => {
      state[accessor] = "";
    });
    return state;
  };

 useEffect(() => {
  if (isOpen) {
    setToast(null);

    console.log("[Modal] isOpen:", isOpen, "mode:", mode, "opname:", opname);

    if (mode === "edit" && opname) {
      const currentForm = generateInitialFormState();

      // Fill in known fields safely
      fields.forEach(({ accessor }) => {
        currentForm[accessor] = opname?.[accessor] ?? "";
      });

      // Include the ID for API
      if (opname?.opname_id) {
        currentForm.id = opname.opname_id;
      }

      console.log("[Modal] Raw fetched opname:", opname);

      if (currentForm.opname_date) {
        const dateObj = new Date(currentForm.opname_date);
        const yyyy = dateObj.getFullYear();
        const mm = String(dateObj.getMonth() + 1).padStart(2, '0');
        const dd = String(dateObj.getDate()).padStart(2, '0');
        currentForm.opname_date = `${yyyy}-${mm}-${dd}`;
      }

      console.log("[Modal] currentForm after merge & date format:", currentForm);

      setForm(currentForm);
      setInitialForm(currentForm);
    } else {
      const blank = generateInitialFormState();

      console.log("[Modal] New blank form:", blank);

      setForm(blank);
      setInitialForm(blank);
    }
  }
}, [isOpen, mode, opname]);


  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

  const handleSubmit = async () => {
    if (!form.opname_date || form.opname_date.toString().trim() === "") {
      setToast({ message: "Tanggal harus diisi.", type: "error" });
      return;
    }

    setLoading(true);
    try {
      if (mode === "edit") {
        await apiClient.put(`/stock-opname/draft/${form.id}`, form);
        setToast({ message: "Draft stock opname berhasil diperbarui!", type: "success" });
      } else {
        await apiClient.post("/stock-opname/draft", form);
        setToast({ message: "Draft stock opname berhasil ditambahkan!", type: "success" });
      }

      onSuccess();
      setTimeout(() => {
        close();
      }, 1000);
    } catch (err) {
      const message = getFriendlyErrorMessage(err);
      setToast({ message, type: "error" });
    } finally {
      setLoading(false);
    }
  };

  const handleReset = () => {
    setForm(initialForm);
  };

  return (
    <Modal isOpen={isOpen} close={close} contentClassName="w-full max-w-2xl">
      <h2 className="text-xl font-semibold mb-4 text-center py-5">
        {mode === "edit" ? "Edit Draft Stock Opname" : "Tambah Draft Stock Opname"}
      </h2>

      <div className="max-h-[60vh] overflow-y-auto pr-2 px-5">
        <div className=" gap-4 w-full">
          {fields.map(({ accessor, label, type }) => (
            <InputField
              key={accessor}
              label={label}
              value={form[accessor]}
              onChange={handleChange(accessor)}
              placeholder={label}
              type={type || "text"}
            />
          ))}
        </div>
      </div>

      <div className="mt-6 flex justify-between gap-4">
        <button
          onClick={handleReset}
          className="w-full bg-gray-200 border border-black text-black rounded-md py-2 hover:bg-gray-300 transition"
        >
          Reset
        </button>
        <Button onClick={handleSubmit} disabled={loading} className="w-full">
          {loading ? "Menyimpan..." : mode === "edit" ? "Update" : "Simpan"}
        </Button>
      </div>

      {toast && (
        <Toast
          message={toast.message}
          type={toast.type}
          onClose={() => setToast(null)}
        />
      )}
    </Modal>
  );
}
