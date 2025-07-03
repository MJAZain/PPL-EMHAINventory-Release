import React, { useState, useEffect } from "react";
import Modal from "./modal";
import InputField from "../inputField";
import Button from "../buttonComp";
import { apiClient } from "../../config/api";
import Toast from "../toast";
import { getFriendlyErrorMessage } from "../../utils/errorHandler";
import TextArea from "../textareacomp";
import Select from "../SelectComp";

const formFields = [
  { label: "Nama", key: "name", placeholder: "Nama produk" },
  { label: "SKU", key: "code", placeholder: "Kode unik produk" },
  { label: "Barcode", key: "barcode", placeholder: "Barcode produk", type:"phone" },
  { label: "Stok Minimal", key: "min_stock", placeholder: "Stok Minimal", type: "number" },
  { label: "Harga Jual", key: "selling_price", placeholder: "Harga jual", type: "number" },
  { label: "Golongan Obat", key: "drug_category_id", type: "select", optionsKey: "drugCategories" },
  { label: "Kategori Obat", key: "category_id", type: "select", optionsKey: "categories" },
  { label: "Satuan", key: "unit_id", type: "select", optionsKey: "units" },
  { label: "Lokasi", key: "storage_location_id", type: "select", optionsKey: "storageLocations" },
  { label: "Merk", key: "brand_id", type: "select", optionsKey: "brands" },
  { label: "Dosis", key: "dosage_description", placeholder: "Deskripsi dosis" },
  { label: "Komposisi", key: "composition_description", placeholder: "Deskripsi komposisi" },
];

export default function BarangModal({ isOpen, close, onSuccess, mode = "add", product = null }) {
  const [loading, setLoading] = useState(false);
  const [toast, setToast] = useState(null);
  const [form, setForm] = useState({});

  const [brands, setBrands] = useState([]);
  const [drugCategories, setDrugCategories] = useState([]);
  const [units, setUnits] = useState([]);
  const [categories, setCategories] = useState([]);
  const [storageLocations, setStorageLocations] = useState([]);

  const generateInitialFormState = () => {
    const state = {};
    formFields.forEach(({ key }) => {
      state[key] = "";
    });
    return state;
  };

  useEffect(() => {
    if (isOpen) {
      setToast(null);
      if (mode === "edit" && product) {
        const base = generateInitialFormState();
        const populated = { ...base, ...product };

        formFields.forEach(({ key }) => {
          if (typeof populated[key] === "object" && populated[key] !== null) {
            populated[key] = populated[key].id;
          }
        });

        setForm(populated);
      } else {
        setForm(generateInitialFormState());
      }
    }

  }, [isOpen, mode, product]);

  useEffect(() => {
    const fetchDropdowns = async () => {
      try {
        const [brandsRes, drugCatRes, unitsRes, catRes, locationsRes] = await Promise.all([
          apiClient.get("/brands?page=1&limit=10000"),
          apiClient.get("/drug-categories/"),
          apiClient.get("/units/"),
          apiClient.get("/categories/"),
          apiClient.get("/storage-locations?page=1&limit=10000"),
        ]);

        setBrands(brandsRes?.data?.data?.data || []);
        setDrugCategories(drugCatRes?.data?.data || []);
        setUnits(unitsRes?.data?.data || []);
        setCategories(catRes?.data?.data || []);
        setStorageLocations(locationsRes?.data?.data?.data || []);
      } catch (error) {
        console.error("Failed to fetch dropdown data:", error);
      }
    };

    fetchDropdowns();
  }, []);

  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

  const handleSubmit = async () => {
    const allFilled = Object.values(form).every((val) => val?.toString().trim() !== "");
    if (!allFilled) {
      setToast({ message: "Semua kolom harus diisi.", type: "error" });
      return;
    }

    const payload = {
      ...form,
      drug_category_id: Number(form.drug_category_id),
      category_id: Number(form.category_id),
      unit_id: Number(form.unit_id),
      selling_price: Number(form.selling_price),
      storage_location_id: Number(form.storage_location_id),
      brand_id: Number(form.brand_id),
      min_stock: Number(form.min_stock || 0),
    };

    setLoading(true);
    try {
      if (mode === "edit" && product?.id) {
        await apiClient.put(`/products/${product.id}`, payload);
        setToast({ message: "Produk berhasil diperbarui!", type: "success" });
      } else {
        await apiClient.post("/products/", payload);
        setToast({ message: "Produk berhasil ditambahkan!", type: "success" });
      }

      onSuccess();
      close();
    } catch (err) {
      const message = getFriendlyErrorMessage(err);
      setToast({ message, type: "error" });
    } finally {
      setForm(generateInitialFormState())
      setToast(false)
      setLoading(false);
    }
  };

  return (
    <Modal isOpen={isOpen} close={close} contentClassName="w-full max-w-2xl">
      <h2 className="text-xl font-semibold mb-4 text-center py-5">
        {mode === "edit" ? "Edit Obat" : "Tambah Obat"}
      </h2>
      <div className="max-h-[60vh] overflow-y-auto pr-2">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 w-full">
          {formFields.map(({ label, key, placeholder, type, optionsKey }) => {
            const isTextArea = key.includes("description");

            if (type === "select") {
              const optionsRaw = {
                brands,
                drugCategories,
                units,
                categories,
                storageLocations,
              }[optionsKey];
              const options = Array.isArray(optionsRaw) ? optionsRaw : [];

              return (
                <div key={key} className="flex flex-col">
                  <label className="text-sm font-medium mb-1">{label}</label>
                  <Select
                    value={form[key]}
                    onChange={handleChange(key)}
                  >
                    <option value="">Pilih {label}</option>
                    {options.map((opt) => (
                      <option key={opt.id} value={opt.id}>
                        {opt.name}
                      </option>
                    ))}
                  </Select>
                </div>
              );
            }

            if (isTextArea) {
              return (
                <div key={key} className="flex flex-col col-span-full">
                  <label className="text-sm font-medium mb-1">{label}</label>
                  <TextArea
                    value={form[key]}
                    onChange={handleChange(key)}
                    placeholder={placeholder}
                    rows={4}
                  />
                </div>
              );
            }

            return (
              <InputField
                key={key}
                label={label}
                value={form[key]}
                onChange={handleChange(key)}
                placeholder={placeholder}
                type={type || "text"}
              />
            );
          })}
        </div>
      </div>

      <div className="mt-6 flex justify-between gap-4">
        <Button onClick={handleSubmit} disabled={loading} className="w-full">
          {loading ? "Menyimpan..." : mode === "edit" ? "Update" : "Simpan"}
        </Button>
        <button
          onClick={() => setForm(generateInitialFormState())}
          className="w-full bg-gray-200 border border-black text-black rounded-md py-2 hover:bg-gray-300 transition"
        >
          Reset
        </button>
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
