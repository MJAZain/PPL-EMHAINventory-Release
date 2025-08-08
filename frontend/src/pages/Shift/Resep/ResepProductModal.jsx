import React, { useEffect, useState } from "react";
import Modal from "../../../components/modal/modal";
import Toast from "../../../components/toast";
import InputField from "../../../components/inputField";
import { apiClient } from "../../../config/api";
import Button from "../../../components/buttonComp";

export default function ResepProductModal({ isOpen, onClose, onSave, initialData = null }) {
  const [products, setProducts] = useState([]);
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedProduct, setSelectedProduct] = useState(null);
  const [form, setForm] = useState({
    quantity: "",
    unit_price: "",
  });

  const [toast, setToast] = useState(null);

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
        unit_price: foundProduct.selling_price || 0,
      });
      setSearchTerm(foundProduct.name || "");
    } else {
      setSelectedProduct(null);
      setForm({ quantity: "", unit_price: "" });
      setSearchTerm("");
    }
  }, [initialData, isOpen, products]);

  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

  const handleSelectProduct = (product) => {
    const existingItems = JSON.parse(localStorage.getItem("presList") || "[]");
    
    // Check if product already exists (only when not editing)
    if (!initialData && existingItems.some(item => item.product_id === product.id)) {
      setToast({ message: "Produk ini sudah ada dalam daftar.", type: "error" });
      return;
    }
    
    setSelectedProduct(product);
    setSearchTerm(product.name);
    setForm((prev) => ({
      ...prev,
      unit_price: product.selling_price || 0,
    }));
  };

  const handleSave = () => {
    if (!selectedProduct || !form.quantity || !form.unit_price) {
      setToast({ message: "Semua field wajib harus diisi.", type: "error" });
      return;
    }

    const newItem = {
      product_id: selectedProduct.id,
      code: selectedProduct.code,
      name: selectedProduct.name,
      quantity: Number(form.quantity),
      unit: selectedProduct.unit?.name ?? "",
      price: Number(form.quantity) * Number(form.unit_price),
    };

    const existing = JSON.parse(localStorage.getItem("presList") || "[]");

    if (!initialData && existing.some(item => item.product_id === selectedProduct.id)) {
      setToast({ message: "Produk ini sudah ada dalam daftar.", type: "error" });
      return;
    }

    const updated = initialData
      ? existing.map((item) =>
          item.product_id === initialData.product.id ? newItem : item
        )
      : [...existing, newItem];

    localStorage.setItem("presList", JSON.stringify(updated));
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
            .map((product) => {
              const existingItems = JSON.parse(localStorage.getItem("presList") || "[]");
              const isProductInList = existingItems.some(item => item.product_id === product.id);
              
              return (
                <li
                  key={product.id}
                  className={`p-2 hover:bg-gray-100 cursor-pointer border-b ${
                    isProductInList && !initialData?.product?.id === product.id 
                      ? "opacity-50 cursor-not-allowed" 
                      : ""
                  }`}
                  onClick={() => {
                    if (!isProductInList || initialData?.product?.id === product.id) {
                      handleSelectProduct(product);
                    }
                  }}
                >
                  <div className="font-medium">{product.name}</div>
                  <div className="text-sm text-gray-600">
                    {product.code} – {product.brand?.name || ""}
                    {isProductInList && !initialData?.product?.id === product.id && (
                      <span className="text-red-500 ml-2">(Sudah dipilih)</span>
                    )}
                  </div>
                </li>
              );
            })}
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
              label="Jumlah Pembelian"
              value={form.quantity}
              onChange={handleChange("quantity")}
              placeholder="Jumlah beli obat"
              type="number"
            />

            <InputField
              label="Harga Jual"
              value={form.unit_price || ""}
              disabled
            />
            <InputField
              label="Total Harga"
              value={
                form.quantity && form.unit_price
                  ? new Intl.NumberFormat("id-ID", {
                      style: "currency",
                      currency: "IDR",
                      minimumFractionDigits: 0,
                    })
                      .format(Number(form.quantity) * Number(form.unit_price))
                      .replace("Rp", "Rp.")
                  : ""
              }
              disabled
            />
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
              unit_price: initialData?.unit_price?.toString() || "",
            });
            setSearchTerm("");
          }}
        >
          Reset
        </button>
        <Button onClick={handleSave} className="w-full">
          Konfirmasi
        </Button>
      </div>
    </Modal>
  );
}