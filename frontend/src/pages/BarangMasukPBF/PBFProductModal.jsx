import React, { useEffect, useState } from "react";
import Modal from "../../components/modal/modal";
import Toast from "../../components/toast";
import InputField from "../../components/inputField";
import { apiClient } from "../../config/api";
import Button from "../../components/buttonComp";

export default function PBFProductModal({ isOpen, onClose, onSave, initialData = null }) {
  const [products, setProducts] = useState([]);
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedProduct, setSelectedProduct] = useState(null);
  const [form, setForm] = useState({
    quantity: "",
    purchase_price: "",
    expiry_date: "",
  });

  const [toast, setToast] = useState(null);

  const order_number = JSON.parse(localStorage.getItem("pbfForm") || "{}")?.order_number || "";

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const res = await apiClient.get("/products/");
        setProducts(res.data?.data || []);
      } catch (err) {
        console.error("Failed to fetch products", err);
      }
    };
    fetchProducts();
  }, []);

  useEffect(() => {
    if (initialData && initialData.product) {
      const foundProduct =
        products.find((p) => p.id === initialData.product.id) || initialData.product;

      setSelectedProduct(foundProduct);
      setForm({
        quantity: initialData.quantity?.toString() || "",
        purchase_price: initialData.purchase_price?.toString() || "",
        expiry_date: initialData.expiry || "",
      });
      setSearchTerm(foundProduct.name || "");
    } else {
      setSelectedProduct(null);
      setForm({ quantity: "", purchase_price: "", expiry_date: "" });
      setSearchTerm("");
    }
  }, [initialData, isOpen, products]);

  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

  const handleSelectProduct = (product) => {
    setSelectedProduct(product);
    setSearchTerm(product.name);
  };

  const handleSave = () => {
    if (!selectedProduct || !form.quantity || !form.purchase_price || !form.expiry_date) {
      setToast({ message: "Semua field wajib harus diisi.", type: "error" });
      return;
    }

    const newItem = {
      product: selectedProduct,
      quantity: Number(form.quantity),
      purchase_price: Number(form.purchase_price),
      expiry: form.expiry_date,
      product_batch: order_number,
    };

    const existing = JSON.parse(localStorage.getItem("barangList") || "[]");

    if (
      !initialData &&
      existing.some((item) => item.product?.id === selectedProduct.id)
    ) {
      setToast({ message: "Produk ini sudah ditambahkan.", type: "error" });
      return;
    }

    const updated = initialData
      ? existing.map((item) =>
          item.product?.id === initialData.product.id ? newItem : item
        )
      : [...existing, newItem];

    localStorage.setItem("barangList", JSON.stringify(updated));
    onSave(updated);
    onClose();
  };

  return (
    <Modal isOpen={isOpen} close={onClose}>
      <h2 className="text-xl text-center font-semibold mb-4">
        {initialData ? "Edit Barang" : "Tambah Barang"}
      </h2>

      <div className="max-h-[60vh] overflow-y-auto pr-2">
        <InputField
          label="Cari Produk"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          placeholder="Nama produk"
          className="mb-2"
        />

        <ul className="border p-2 rounded mb-4 max-h-40 overflow-y-auto">
          {products
            .filter((p) => p.name.toLowerCase().includes(searchTerm.toLowerCase()))
            .map((product) => (
              <li
                key={product.id}
                className="p-2 hover:bg-gray-100 cursor-pointer border-b"
                onClick={() => handleSelectProduct(product)}
              >
                <div className="font-medium">{product.name}</div>
                <div className="text-sm text-gray-600">
                  {product.code} – {product.brand?.name || ""}
                </div>
              </li>
            ))}
        </ul>

        {selectedProduct && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <InputField label="Kode Barang" value={selectedProduct.code} disabled />
            <InputField label="Nama Barang" value={selectedProduct.name || ""} disabled />
            <InputField
              label="Satuan"
              value={typeof selectedProduct.unit?.name === "string" ? selectedProduct.unit.name : ""}
              disabled
            />
            <InputField
              label="Jumlah Masuk"
              value={form.quantity}
              onChange={handleChange("quantity")}
              placeholder="Jumlah barang"
              type="number"
            />
            <InputField
              label="Harga Beli"
              value={form.purchase_price}
              onChange={handleChange("purchase_price")}
              placeholder="Harga beli"
              type="number"
            />
            <InputField
              label="Total Harga"
              value={
                form.quantity && form.purchase_price
                  ? new Intl.NumberFormat("id-ID", {
                      style: "currency",
                      currency: "IDR",
                      minimumFractionDigits: 0,
                    })
                      .format(Number(form.quantity) * Number(form.purchase_price))
                      .replace("Rp", "Rp.")
                  : ""
              }
              disabled
            />
            <InputField
              label="Tanggal Kedaluwarsa"
              value={form.expiry_date}
              onChange={handleChange("expiry_date")}
              type="date"
            />
            <InputField label="No. Batch" value={order_number} disabled />
          </div>
        )}
      </div>

      {toast && (
        <Toast
          message={toast.message}
          type={toast.type}
          onClose={() => setToast(null)}
        />
      )}

      <div className="mt-6 flex justify-between gap-4">
        <button
          className="text-black w-full bg-gray-200 border border-black hover:bg-gray-300 rounded-md"
          onClick={() => {
            setSelectedProduct(null);
            setForm({
              quantity: initialData?.quantity?.toString() || "",
              purchase_price: initialData?.purchase_price?.toString() || "",
              expiry_date: initialData?.expiry || "",
            });
            setSearchTerm("");
          }}
        >
          Reset
        </button>
        <Button onClick={handleSave} className="w-full">
          Simpan
        </Button>
      </div>
    </Modal>
  );
}
