import React, { useEffect, useState } from "react";
import Modal from "../../components/modal/modal";
import Toast from "../../components/toast";
import InputField from "../../components/inputField";
import { apiClient } from "../../config/api";
import Button from "../../components/buttonComp";

export default function OpnameProductModal({ isOpen, onClose, onSave, initialData = null }) {
  const [products, setProducts] = useState([]);
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedProduct, setSelectedProduct] = useState(null);
  const [toast, setToast] = useState(null);

  // Fetch products on mount
  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const res = await apiClient.get("/stocks/current");
        const mapped = (res.data?.data || []).map((item) => ({
          id: item.product_id,
          name: item.product_name,
          quantity: item.quantity,
        }));
        setProducts(mapped);
      } catch (err) {
        console.error("Failed to fetch products", err);
      }
    };
    fetchProducts();
  }, []);

  // Reset modal state when opened
  useEffect(() => {
    if (initialData?.product) {
      const foundProduct =
        products.find((p) => p.id === initialData.product.id) || initialData.product;

      setSelectedProduct(foundProduct);
      setSearchTerm(foundProduct.name || "");
    } else {
      setSelectedProduct(null);
      setSearchTerm("");
    }
  }, [initialData, isOpen, products]);

  const handleSelectProduct = (product) => {
    const existingItems = JSON.parse(localStorage.getItem("opList") || "[]");

    if (!initialData && existingItems.some((item) => item.product_id === product.id)) {
      setToast({ message: "Produk ini sudah ada dalam daftar.", type: "error" });
      return;
    }

    setSelectedProduct(product);
    setSearchTerm(product.name);
  };

const handleSave = async () => {
  if (!selectedProduct) {
    setToast({ message: "Pilih produk terlebih dahulu.", type: "error" });
    return;
  }

  const newItem = {
    product_id: selectedProduct.id,
    product_name: selectedProduct.name,
    quantity: selectedProduct.quantity,
  };

  const existing = JSON.parse(localStorage.getItem("opList") || "[]");

  const updated = initialData
    ? existing.map((item) =>
        item.product_id === initialData.product.id ? newItem : item
      )
    : [...existing, newItem];

  localStorage.setItem("opList", JSON.stringify(updated));

  const opnameId = localStorage.getItem("opnameId");
  if (!opnameId) {
    setToast({ message: "opnameId tidak ditemukan di localStorage.", type: "error" });
    return;
  }

  try {
const res = await apiClient.post(`/stock-opname/draft/${opnameId}/products`, {
  product_id: String(selectedProduct.id),
});

    const detailId = res?.data?.data?.detail_id;
    if (!detailId) throw new Error("detail_id tidak ditemukan.");

    const latestOpList = JSON.parse(localStorage.getItem("opList") || "[]");
    const updatedWithDetailId = latestOpList.map((item) =>
      item.product_id === selectedProduct.id
        ? { ...item, detail_id: detailId }
        : item
    );

    localStorage.setItem("opList", JSON.stringify(updatedWithDetailId));

    setToast({ message: "Produk berhasil ditambahkan ke draft.", type: "success" });
    onSave(updatedWithDetailId);
    onClose();
  } catch (error) {
    console.error("Gagal mengirim produk ke server:", error);
    setToast({
      message: "Gagal mengirim produk ke server. Coba lagi.",
      type: "error",
    });
  }
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
              const existingItems = JSON.parse(localStorage.getItem("opList") || "[]");
              const isProductInList = existingItems.some(
                (item) => item.product_id === product.id
              );
              const isEditingCurrent = initialData?.product?.id === product.id;

              return (
                <li
                  key={product.id}
                  className={`p-2 hover:bg-gray-100 cursor-pointer border-b ${
                    isProductInList && !isEditingCurrent
                      ? "opacity-50 cursor-not-allowed"
                      : ""
                  }`}
                  onClick={() => {
                    if (!isProductInList || isEditingCurrent) {
                      handleSelectProduct(product);
                    }
                  }}
                >
                  <div className="font-medium">{product.name}</div>
                  <div className="text-sm text-gray-600">
                    {product.code}
                    {isProductInList && !isEditingCurrent && (
                      <span className="text-red-500 ml-2">(Sudah dipilih)</span>
                    )}
                  </div>
                </li>
              );
            })}
        </ul>

        {selectedProduct && (
          <div className="grid grid-cols-1 md:grid-cols-1 gap-4 mb-4">
            <InputField label="Nama Barang" value={selectedProduct.name} disabled />
            <InputField label="Kuantitas Sistem" value={selectedProduct.quantity} disabled />
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
