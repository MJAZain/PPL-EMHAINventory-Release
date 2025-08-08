import React, { useState, useEffect } from "react";
import Sidebar from "../../components/Sidebar";
import Button from "../../components/buttonComp";
import InputField from "../../components/inputField";
import Toast from "../../components/toast";
import { apiClient } from "../../config/api";
import TextArea from '../../components/textareacomp'

export default function KoreksiFormPage() {
  const [searchTerm, setSearchTerm] = useState("");
  const [toast, setToast] = useState(null);
  const [products, setProducts] = useState([]);
  const [form, setForm] = useState({
    product_id: "",
    old_stock: 0,
    new_stock: "",
    reason: "",
    correction_date: new Date().toISOString().split("T")[0],
    notes: "",
  });

  useEffect(() => {
    fetchProducts();
  }, []);

  const fetchProducts = async () => {
    try {
      const response = await apiClient.get("/stocks/current");
      setProducts(response.data?.data ?? []);
    } catch (error) {
      console.error("Failed to fetch current stock", error);
    }
  };

  const handleSelectProduct = (product) => {
    setForm((prev) => ({
      ...prev,
      product_id: product.product_id,
      old_stock: product.quantity,
    }));
    setSearchTerm(product.product_name);
  };

  const handleChange = (key) => (e) => {
    const value = e.target.value;
    setForm((prev) => ({ ...prev, [key]: value }));
  };

const handleSubmit = async () => {
  const { product_id, new_stock, reason, notes } = form;

  if (!product_id || new_stock === "" || !reason) {
    setToast({ message: "Mohon isi semua field yang wajib diisi.", type: "error" });
    return;
  }

  const payload = {
    product_id: parseInt(product_id),
    old_stock: form.old_stock,
    new_stock: parseInt(new_stock),
    difference: parseInt(new_stock) - form.old_stock,
    reason,
    correction_date: new Date(form.correction_date).toISOString(), // 🟢 ISO timestamp
    notes,
  };

  console.log("Submitting stock correction payload:", payload);

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
        notes: "",
      });
      setSearchTerm("");
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
            <label className="text-sm font-medium mb-1">Cari Produk</label>
            <InputField
              type="text"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              placeholder="Cari produk berdasarkan nama atau kode…"
            />

            {searchTerm && (
              <ul className="border p-2 rounded mt-2 mb-4 max-h-40 overflow-y-auto bg-white shadow">
                {products
                  .filter((p) =>
                    p.product_name.toLowerCase().includes(searchTerm.toLowerCase())
                  )
                  .map((product) => {
                    const existingItems = JSON.parse(localStorage.getItem("presList") || "[]");
                    const isProductInList = existingItems.some(
                      (item) => item.product_id === product.product_id
                    );

                    return (
                      <li
                        key={product.product_id}
                        className={`p-2 hover:bg-gray-100 cursor-pointer border-b last:border-b-0 ${
                          isProductInList && form.product_id !== product.product_id
                            ? "opacity-50 cursor-not-allowed"
                            : ""
                        }`}
                        onClick={() => {
                          if (!isProductInList || form.product_id === product.product_id) {
                            handleSelectProduct(product);
                          }
                        }}
                      >
                        <div className="font-medium">{product.product_name}</div>
                        <div className="text-sm text-gray-600">
                          {product.product_code}{" "}
                          {isProductInList && form.product_id !== product.product_id && (
                            <span className="text-red-500 ml-2">(Sudah dipilih)</span>
                          )}
                        </div>
                      </li>
                    );
                  })}
              </ul>
            )}
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
            label="Tanggal Koreksi"
            value={form.correction_date}
            onChange={handleChange("correction_date")}
            type="date"
          />

          <div className="flex flex-col mb-4">
            <label className="text-sm font-medium mb-1">Catatan</label>
            <TextArea
              value={form.notes}
              onChange={handleChange("notes")}
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
