import React, { useState, useEffect } from "react";
import Modal from "../../components/modal/modal";
import InputField from "../../components/inputField";
import Button from "../../components/buttonComp";
import { apiClient } from "../../config/api";
import Toast from "../../components/toast";
import { getFriendlyErrorMessage } from "../../utils/errorHandler";
import Select from '../../components/SelectComp'

export default function PengeluaranModal({ isOpen, close, onSuccess, mode = "add", pengeluaran = null }) {
  const [loading, setLoading] = useState(false);
  const [toast, setToast] = useState(null);
  const [form, setForm] = useState({
    expense_type_id: "",
    amount: "",
    description: "",
    date: "",
  });
  const [expenseTypes, setExpenseTypes] = useState([]);

  useEffect(() => {
    if (isOpen) {
      setToast(null);
      if (mode === "edit" && pengeluaran) {
        setForm({
          expense_type_id: pengeluaran.expense_type_id,
          amount: pengeluaran.amount,
          description: pengeluaran.description || "",
          date: pengeluaran.date ? pengeluaran.date.slice(0, 10) : "",
        });
      } else {
        setForm({
          expense_type_id: "",
          amount: "",
          description: "",
          date: "",
        });
      }

      fetchExpenseTypes();
    }
  }, [isOpen, mode, pengeluaran]);

  const fetchExpenseTypes = async () => {
    try {
      const res = await apiClient.get("/expense-types/");
      setExpenseTypes(res.data?.data || []);
    } catch {
      setExpenseTypes([]);
    }
  };

  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

  const handleSubmit = async () => {
    const allFilled = form.expense_type_id && form.amount && form.date;
    if (!allFilled) {
      setToast({ message: "Kolom wajib harus diisi.", type: "error" });
      return;
    }

    setLoading(true);
    try {
      const payload = {
        ...form,
        expense_type_id: parseInt(form.expense_type_id, 10),
        amount: parseFloat(form.amount),
        date: new Date(form.date + "T00:00:00Z").toISOString(),
      };

      if (mode === "edit" && pengeluaran?.id) {
        await apiClient.put(`/expenses/${pengeluaran.id}`, payload);
        setToast({ message: "Pengeluaran berhasil diperbarui!", type: "success" });
      } else {
        await apiClient.post("/expenses/", payload);
        setToast({ message: "Pengeluaran berhasil ditambahkan!", type: "success" });
      }

      onSuccess();
      close();
    } catch (err) {
      const message = getFriendlyErrorMessage(err);
      setToast({ message, type: "error" });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal isOpen={isOpen} close={close} contentClassName="w-full max-w-2xl">
      <h2 className="text-xl font-semibold mb-4 text-center py-5">
        {mode === "edit" ? "Edit Pengeluaran" : "Tambah Pengeluaran"}
      </h2>

      <div className="max-h-[60vh] overflow-y-auto pr-2 px-5">
        <div className="grid grid-cols-1 gap-4 w-full">
          <div>
            <label>Jenis Pengeluaran</label>
            <Select
              value={form.expense_type_id}
              onChange={handleChange("expense_type_id")}
              className="w-full border rounded p-2"
            >
              <option value="">Pilih Jenis</option>
              {expenseTypes.map((type) => (
                <option key={type.id} value={type.id}>
                  {type.name}
                </option>
              ))}
            </Select>
          </div>

          <InputField
            label="Jumlah"
            value={form.amount}
            onChange={handleChange("amount")}
            placeholder="Jumlah"
            type="number"
          />

          <InputField
            label="Deskripsi"
            value={form.description}
            onChange={handleChange("description")}
            placeholder="Deskripsi"
            type="text"
          />

          <div>
            <label>Tanggal</label>
            <InputField
              type="date"
              value={form.date}
              onChange={handleChange("date")}
              className="w-full border rounded p-2"
            />
          </div>
        </div>
      </div>

      <div className="mt-6 flex justify-between gap-4">
        <button
          onClick={() =>
            setForm({ expense_type_id: "", amount: "", description: "", date: "" })
          }
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
