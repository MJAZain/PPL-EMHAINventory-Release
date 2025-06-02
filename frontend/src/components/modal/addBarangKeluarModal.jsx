import React, { useEffect, useState } from "react";
import Modal from "./modal";
import InputField from "../inputField";
import { apiClient } from "../../config/api";
import Button from "../buttonComp";

export default function TambahBarangKeluarModal({
  isOpen,
  onClose,
  onSave,
  initialData = null,
}) {
  const [selectedProduct, setSelectedProduct] = useState(null);
  const [quantity, setQuantity] = useState("");
  const [price, setPrice] = useState("");

  const [productOptions, setProductOptions] = useState([]);
  const [searchTerm, setSearchTerm] = useState("");
  const [toast, setToast] = useState(null);

  // Reset modal state on open or when initialData changes
  useEffect(() => {
    if (isOpen) {
      setSelectedProduct(initialData?.product || null);
      setQuantity(initialData?.quantity?.toString() || "");
      setPrice(initialData?.price?.toString() || "");
      setSearchTerm(initialData?.product?.name || "");
      setToast(null);
    }
  }, [isOpen, initialData]);

  // Fetch products when modal opens
  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const res = await apiClient.get("/products/");
        setProductOptions(res.data?.data || []);
      } catch (err) {
        console.error("Failed to fetch products:", err);
        setProductOptions([]);
      }
    };

    if (isOpen) {
      fetchProducts();
    }
  }, [isOpen]);

  // Auto-calculate price based on selected product's selling price × quantity
  useEffect(() => {
    if (selectedProduct && quantity) {
      const total = selectedProduct.selling_price * Number(quantity);
      setPrice(total.toFixed(2));
    } else {
      setPrice("");
    }
  }, [selectedProduct, quantity]);

  const filteredProducts = productOptions.filter((p) =>
    p.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleSave = () => {
    if (!selectedProduct || !quantity) {
      setToast({ message: "Mohon pilih produk dan isi kuantitas.", type: "error" });
      return;
    }

    const newItem = {
      product: selectedProduct,
      quantity: Number(quantity),
      price: Number(price),
    };

    const existing = JSON.parse(localStorage.getItem("barangList") || "[]");

    const isDuplicate =
      !initialData && existing.some((item) => item.product.id === selectedProduct.id);

    if (isDuplicate) {
      setToast({ message: "Produk ini sudah ditambahkan.", type: "error" });
      return;
    }

    const updated = initialData
      ? existing.map((item) =>
          item.product.id === initialData.product.id ? newItem : item
        )
      : [...existing, newItem];

    localStorage.setItem("barangList", JSON.stringify(updated));
    onSave(updated);
    onClose();
  };

  return (
    <Modal isOpen={isOpen} close={onClose}>
      <h2 className="text-xl font-semibold mb-4">
        {initialData ? "Edit Barang" : "Tambah Barang"}
      </h2>

      {toast && (
        <div className="text-red-600 font-medium mb-4">{toast.message}</div>
      )}

      <InputField
        label="Cari Produk"
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        placeholder="Nama produk"
        className="w-full mb-2 h-10"
      />

      <div className="border p-2 rounded mb-4 max-h-40 overflow-y-auto">
        {filteredProducts.map((product) => (
          <li
            key={product.id}
            className="p-2 hover:bg-gray-100 cursor-pointer border-b"
            onClick={() => {
              setSelectedProduct(product);
              setSearchTerm(product.name);
            }}
          >
            <div className="font-medium">{product.name}</div>
            <div className="text-sm text-gray-600">
              {product.code} – {product.brand}
            </div>
          </li>
        ))}
      </div>

      <InputField
        label="Kuantitas"
        type="number"
        value={quantity}
        onChange={(e) => setQuantity(e.target.value)}
        placeholder="Jumlah yang dijual"
        className="w-full h-10"
      />

      <InputField
        label="Harga Total"
        type="number"
        value={price}
        readOnly
        className="w-full h-10 bg-gray-100 cursor-not-allowed"
      />

      <Button onClick={handleSave}>Simpan</Button>
    </Modal>
  );
}
