import React, { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import Sidebar from "../../components/Sidebar";
import Button from "../../components/buttonComp";
import InputField from "../../components/inputField";
import Toast from "../../components/toast";
import { apiClient } from "../../config/api";
import Select from '../../components/SelectComp'
import TextArea from '../../components/textareacomp'

const formFields = [
  { label: "Tanggal Pesanan", key: "order_date", placeholder: "Tanggal Pemesanan", type: "date" },
  { label: "Tanggal Barang Masuk", key: "receipt_date", placeholder: "Tanggal Barang Masuk", type: "date" },
  { label: "Nama Supplier / Non-PBF", key: "supplier_id", type: "select", optionsKey: "supplier" },
  { label: "Nomor Faktur", key: "invoice_number", placeholder: "Nomor Faktur", type: "phone"},
  {
    label: "Jenis Transaksi",
    key: "transaction_type",
    type: "select",
    options: ["Cash", "Kredit", "Konsinyasi"],
  },
  { label: "Tanggal Jatuh Tempo Pembayaran", key: "payment_due_date", placeholder: "Tanggal Jatuh Tempo", type: "date" },
  { label: "Catatan", key: "additional_notes", placeholder: "Catatan Tambahan" },
];

export default function NonPBFDetailPage() {
  const [toast, setToast] = useState(null);
  const [supplier, setSupplier] = useState([]);
  const [form, setForm] = useState({});

  const { id } = useParams();

  const navigate = useNavigate();

 useEffect(() => {
  const fetchDropdownsAndForm = async () => {
    try {
      console.debug("Fetching supplier list...");
      const supplierRes = await apiClient.get("/suppliers/");
      const allSuppliers = supplierRes?.data?.data || [];
      console.debug("All suppliers fetched:", allSuppliers);

      const nonpbfSuppliers = allSuppliers.filter(supplier => supplier.type === "Non-PBF");
      console.debug("Filtered Non-PBF suppliers:", nonpbfSuppliers);

      setSupplier(nonpbfSuppliers);

      if (id) {
        console.debug(`Fetching incoming-nonpbf detail for id: ${id}`);
        const res = await apiClient.get(`/incoming-nonpbf/${id}`);
        const data = res?.data?.data;
        console.debug("Fetched Non-PBF detail:", data);

        if (data) {
          const matchedSupplier = nonpbfSuppliers.find(s => s.name === data.supplier_name);
          console.debug("Matched supplier from name:", matchedSupplier);

          setForm({
            order_date: data.order_date?.split("T")[0] || "",
            receipt_date: data.incoming_date?.split("T")[0] || "",
            supplier_id: matchedSupplier?.id || "",
            supplier_name: data.supplier_name || "",
            invoice_number: data.invoice_number || "",
            transaction_type: data.transaction_type || "",
            payment_due_date: data.payment_due_date?.split("T")[0] || "",
            additional_notes: data.additional_notes || "",
            order_number: data.order_number || "",
          });

          console.debug("Form state set with data:", {
            order_date: data.order_date?.split("T")[0] || "",
            receipt_date: data.incoming_date?.split("T")[0] || "",
            supplier_id: matchedSupplier?.id || "",
            supplier_name: data.supplier_name || "",
            invoice_number: data.invoice_number || "",
            transaction_type: data.transaction_type || "",
            payment_due_date: data.payment_due_date?.split("T")[0] || "",
            additional_notes: data.additional_notes || "",
            order_number: data.order_number || "",
          });
        }
      } else {
        console.debug("No ID provided, skipping detail fetch.");
      }
    } catch (error) {
      console.error("Failed to fetch dropdowns or form data:", error);
    }
  };

  fetchDropdownsAndForm();
}, [id]);

const today = new Date().toISOString().split("T")[0];

  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

  const generateKodeTransaksi = () => {
    const timestamp = Date.now();
    return `Non-PBF-${timestamp}`;
  };

  const handleNext = () => {
    const requiredKeys = formFields
      .filter((field) => field.key !== "additinal_notes")
      .map((field) => field.key);

    const allFilled = requiredKeys.every((key) => form[key] && form[key].toString().trim() !== "");

    if (!allFilled) {
      setToast({
        message: "Mohon isi semua field wajib.",
        type: "error",
      });
      return;
    }

    const selectedSupplier = supplier.find((s) => s.id === Number(form.supplier_id));
    const supplierName = selectedSupplier ? selectedSupplier.name : "";

    const fullForm = {
      ...form,
      supplier_name: supplierName,
      order_number: id ? form.order_number : generateKodeTransaksi(),
    };

    localStorage.setItem("nonpbfForm", JSON.stringify(fullForm));
    console.log("nonpbfForm:", JSON.stringify(fullForm, null, 2));

    if (id) {
      navigate(`/non-pbf-list/${id}`);
    } else {
      navigate("/non-pbf-list");
    }
  };

  return (
    <div className="flex min-h-screen bg-gray-100">
      <div className="bg-white min-h-screen">
        {!id && <Sidebar />}
      </div>
      <div className="bg-white max-w-xl mx-auto w-full p-6 mt-10 border rounded-md border-gray-300 shadow-md">
        <h1 className="text-2xl font-bold text-center mb-6">
          {id ? "Edit Barang Masuk Non-PBF" : "Form Barang Masuk Non-PBF"}
        </h1>

        <div className="max-h-[60vh] overflow-y-auto pr-2">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 w-full">
            {formFields.map(({ label, key, placeholder, type, optionsKey, options }) => {
              const isTextArea = key === "additional_notes";

              if (type === "select") {
                let selectOptions = [];

                if (optionsKey === "supplier") {
                  selectOptions = supplier.map((sup) => ({
                    id: sup.id,
                    name: sup.name,
                  }));
                } else if (Array.isArray(options)) {
                  selectOptions = options.map((val) => ({
                    id: val,
                    name: val,
                  }));
                }

                return (
                  <div key={key} className="flex flex-col">
                    <label className="text-sm font-medium mb-1">{label}</label>
                    <Select
                      value={form[key] || ""}
                      onChange={handleChange(key)}
                    >
                      <option value="">Pilih {label}</option>
                      {selectOptions.map((opt) => (
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
                      value={form[key] || ""}
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
                  value={form[key] || ""}
                  onChange={handleChange(key)}
                  placeholder={placeholder}
                  type={type || "text"}
                  min={type === "date" ? today : undefined}
                />
              );
            })}
          </div>
        </div>

        {toast && (
          <Toast
            message={toast.message}
            type={toast.type}
            onClose={() => setToast(null)}
          />
        )}

        <div className="mt-6">
          <Button className="w-full" onClick={handleNext}>
            {id ? "Lanjut Edit" : "Selanjutnya"}
          </Button>
        </div>
      </div>
    </div>
  );
}
