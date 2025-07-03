import React, { useState, useEffect } from "react";
import Sidebar from "../../components/Sidebar";
import Button from "../../components/buttonComp";
import InputField from "../../components/inputField";
import Toast from "../../components/toast";
import { apiClient } from "../../config/api";
import Select from "../../components/SelectComp";

export default function KoreksiFormPage() {
  const [toast, setToast] = useState(null);
  const [products, setProducts] = useState([]);
  const [form, setForm] = useState({
    product_id: "",
    old_stock: 0,
    new_stock: "",
    reason: "",
    correction_date: new Date().toISOString().split("T")[0],
    correction_officer: "",
    notes: "",
  });

  useEffect(() => {
    fetchProducts();
  }, []);

  const fetchProducts = async () => {
    try {
      const response = await apiClient.get("/stock-current/");
      setProducts(response.data);
    } catch (error) {
      console.error("Failed to fetch current stock", error);
    }
  };

  const handleChange = (key) => (e) => {
    const value = e.target.value;
    const updatedForm = { ...form, [key]: value };

    if (key === "product_id") {
      const selectedProduct = products.find(p => String(p.id) === value);
      updatedForm.old_stock = selectedProduct?.stock || 0;
    }

    setForm(updatedForm);
  };

  const handleSubmit = async () => {
    const { product_id, new_stock, reason, correction_officer } = form;

    if (!product_id || new_stock === "" || !reason || !correction_officer) {
      setToast({ message: "Mohon isi semua field yang wajib diisi.", type: "error" });
      return;
    }

    const payload = {
      ...form,
      new_stock: parseInt(form.new_stock),
      difference: parseInt(form.new_stock) - form.old_stock,
    };

    try {
      const response = await apiClient.post("/stock-corrections/", payload);

      if (response.status === 200 || response.status === 201) {
        setToast({ message: "Koreksi stok berhasil disimpan.", type: "success" });
        setForm({
          product_id: "",
          old_stock: 0,
          new_stock: "",
          reason: "",
          correction_date: new Date().toISOString().split("T")[0],
          correction_officer: "",
          notes: "",
        });
      } else {
        throw new Error("Unexpected server response");
      }
    } catch (error) {
      console.error("Error submitting stock correction:", error);
      setToast({ message: "Gagal menyimpan koreksi stok.", type: "error" });
    }
  };

  return (
    <div className="flex min-h-screen bg-gray-100">
      <div className="bg-white min-h-screen">
        <Sidebar />
      </div>

      <div className="bg-white max-w-xl mx-auto w-full p-6 mt-10 border rounded-md border-gray-300 shadow-md max-h-[90vh] overflow-y-auto">
        <h1 className="text-2xl font-bold text-center mb-6">Koreksi Stok</h1>

        <div className="grid grid-cols-1 gap-4">
          <div className="flex flex-col mb-4">
            <label className="text-sm font-medium mb-1">Produk</label>
            <Select
              value={form.product_id}
              onChange={handleChange("product_id")}
              className="border border-gray-300 rounded-md px-3 py-2"
            >
              <option value="">Pilih Produk</option>
              {products.map((product) => (
                <option key={product.id} value={product.id}>
                  {product.name}
                </option>
              ))}
            </Select>
          </div>

          <InputField
            label="Stok Lama"
            value={form.old_stock}
            type="number"
            disabled
          />

          <InputField
            label="Stok Baru"
            value={form.new_stock}
            onChange={handleChange("new_stock")}
            type="number"
          />

          <InputField
            label="Alasan Koreksi"
            value={form.reason}
            onChange={handleChange("reason")}
          />

          <InputField
            label="Petugas Koreksi"
            value={form.correction_officer}
            onChange={handleChange("correction_officer")}
          />

          <InputField
            label="Tanggal Koreksi"
            value={form.correction_date}
            onChange={handleChange("correction_date")}
            type="date"
          />

          <div className="flex flex-col mb-4">
            <label className="text-sm font-medium mb-1">Catatan</label>
            <textarea
              value={form.notes}
              onChange={handleChange("notes")}
              className="border border-gray-300 rounded-md px-3 py-2 resize-none"
              rows={3}
              placeholder="Tambahkan catatan opsional..."
            />
          </div>

          <Button className="w-full" onClick={handleSubmit}>
            Simpan Koreksi Stok
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
