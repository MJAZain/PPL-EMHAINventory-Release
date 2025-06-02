import React, { useState, useEffect } from "react";
import Modal from "./modal";
import InputField from "../inputField";
import Button from "../buttonComp";
import { apiClient } from "../../config/api";
import Toast from "../toast";
import useProductActions from "../../hooks/useProductsAction";

const initialFormState = {
  name: "",
  code: "",
  barcode: "",
  category_id: "",
  unit_id: "",
  package_content: "",
  purchase_price: "",
  selling_price: "",
  wholesale_price: "",
  stock_buffer: "",
  storage_location: "",
  brand: "",
};

const labelMap = {
  name: "Nama",
  code: "Kode Obat",
  barcode: "Barcode",
  category_id: "Kategori Obat",
  unit_id: "Satuan Obat",
  package_content: "Banyaknya Isi",
  purchase_price: "Harga Beli",
  selling_price: "Harga Jual",
  wholesale_price: "Harga Jual ",
  stock_buffer: "Stok Buffer",
  storage_location: "Lokasi Obat",
  brand: "Nama Brand",
};

function EditBarangModal({ isOpen, close, onSubmit, productId }) {
  const [form, setForm] = useState(initialFormState);
  const [toast, setToast] = useState(null);
  const [categories, setCategories] = useState([]);
  const [units, setUnits] = useState([]);
  const [loadingCategories, setLoadingCategories] = useState(false);
  const [loadingUnits, setLoadingUnits] = useState(false);
  const { getProductById } = useProductActions();

  // Load categories and units
  useEffect(() => {
    const fetchCategories = async () => {
      setLoadingCategories(true);
      try {
        const res = await apiClient.get("/categories/");
        setCategories(res.data.data);
      } catch (err) {
        console.error("Failed to fetch categories", err);
      } finally {
        setLoadingCategories(false);
      }
    };

    const fetchUnits = async () => {
      setLoadingUnits(true);
      try {
        const res = await apiClient.get("/units/");
        setUnits(res.data.data);
      } catch (err) {
        console.error("Failed to fetch units", err);
      } finally {
        setLoadingUnits(false);
      }
    };

    if (isOpen) {
      fetchCategories();
      fetchUnits();
    }
  }, [isOpen]);

  useEffect(() => {
  const loadProduct = async () => {
    if (!productId || !isOpen) return;
    try {
      const data = await getProductById(productId);
      const product = data.data;

      setForm({
        name: product.name || "",
        code: product.code || "",
        barcode: product.barcode || "",
        category_id: product.category_id?.toString() || "",
        unit_id: product.unit_id?.toString() || "",
        package_content: product.package_content ?? "",
        purchase_price: product.purchase_price ?? "",
        selling_price: product.selling_price ?? "",
        wholesale_price: product.wholesale_price ?? "",
        stock_buffer: product.stock_buffer ?? "",
        storage_location: product.storage_location || "",
        brand: product.brand || "",
      });
    } catch (err) {
      setToast({ message: "Gagal memuat data barang.", type: "error" });
    }
  };

  loadProduct();
}, [productId, isOpen]);


  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

  const handleSubmit = async () => {
    const allFilled = Object.values(form).every(
      (val) => val.toString().trim() !== ""
    );
    if (!allFilled) {
      setToast({ message: "Tolong isi setiap kolom.", type: "error" });
      return;
    }

    const payload = {
      name: form.name.trim(),
      code: form.code.trim(),
      barcode: form.barcode.trim(),
      category_id: Number(form.category_id),
      unit_id: Number(form.unit_id),
      package_content: Number(form.package_content),
      purchase_price: Number(form.purchase_price),
      selling_price: Number(form.selling_price),
      wholesale_price: Number(form.wholesale_price),
      stock_buffer: Number(form.stock_buffer),
      storage_location: form.storage_location.trim(),
      brand: form.brand.trim(),
    };

    if (isNaN(payload.package_content)) {
      setToast({
        message: "Isi Paket harus diisi dengan angka yang valid.",
        type: "error",
      });
      return;
    }

    try {
      await apiClient.put(`/products/${productId}`, payload);
      close();
      onSubmit?.();
    } catch (error) {
      const message =
        error.response?.data?.message || "Gagal memperbarui barang.";
      setToast({ message, type: "error" });
      console.error("Update error:", error);
    }
  };

  const numericFields = [
    "stock_buffer",
    "selling_price",
    "purchase_price",
    "wholesale_price",
    "package_content",
  ];

  const selectFieldsConfig = {
    category_id: {
      label: "Kategori",
      options: categories,
      loading: loadingCategories,
      optionLabelKey: "name",
    },
    unit_id: {
      label: "Satuan",
      options: units,
      loading: loadingUnits,
      optionLabelKey: "name",
    },
  };

  return (
    <Modal isOpen={isOpen} close={close}>
      <h2 className="text-xl font-semibold mb-4 text-center">Tambah Barang</h2>

      <div className="grid grid-cols-2 max-h-[60vh] overflow-y-auto pr-2 gap-5">
        {Object.keys(initialFormState).map((key) => {
          if (selectFieldsConfig[key]) {
            const { label, options, loading, optionLabelKey } = selectFieldsConfig[key];
            return (
              <div key={key} className="flex flex-col">
                <label className="mb-1 font-medium">{label}</label>
                {loading ? (
                  <p>Loading {label.toLowerCase()}...</p>
                ) : (
                  <select
                    value={form[key]}
                    onChange={handleChange(key)}
                    className="w-full h-10 border rounded px-2"
                  >
                    <option value="">{`Pilih ${label}`}</option>
                    {options.map((option) => (
                      <option key={option.id} value={option.id}>
                        {option[optionLabelKey]}
                      </option>
                    ))}
                  </select>
                )}
              </div>
            );
          }

          return (
            <InputField
              key={key}
              label={labelMap[key] || key.charAt(0).toUpperCase() + key.slice(1)}
              value={form[key]}
              onChange={handleChange(key)}
              placeholder={`Masukkan ${labelMap[key] || key}`}
              type={numericFields.includes(key) ? "number" : "text"}
              min={numericFields.includes(key) ? 1 : undefined}
              className="w-full h-10"
            />
          );
        })}
      </div>

      <div className="grid grid-cols-2 pr-5 max-h-[60vh] overflow-y-auto gap-5 py-5">
        <Button
          onClick={() => setForm(initialFormState)}
          className="bg-gray-200 text-black border-none"
        >
          Reset
        </Button>
        <Button onClick={handleSubmit}>Simpan</Button>
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

export default EditBarangModal;
