import React, { useState, useEffect } from "react";
import Modal from "../../../components/modal/modal";
import InputField from "../../../components/inputField";
import Button from "../../../components/buttonComp";
import Toast from "../../../components/toast";
import { apiClient } from "../../../config/api";
import Select from "../../../components/SelectComp";
import TextArea from '../../../components/textareacomp'


export default function AkhiriTransaksiModal({   
  isOpen,
  onClose,
  regularList,
  onAfterSubmit,
}) {
  const [form, setForm] = useState({
    customer_name: "",
    customer_contact: "",
    description: "",
    total_discount: "",
    payment_method: "Cash",
  });

  const [subTotal, setSubTotal] = useState(0);
  const [toast, setToast] = useState(null);

  useEffect(() => {
    if (isOpen) {
      const regularList = JSON.parse(localStorage.getItem("regularList") || "[]");
      const sum = regularList.reduce((acc, item) => acc + item.sub_total, 0);
      setSubTotal(sum);
    }
  }, [isOpen]);

  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

  const handleConfirm = async () => {
    const regularList = JSON.parse(localStorage.getItem("regularList") || "[]");
    if (regularList.length === 0) {
      setToast({ message: "Tidak ada item dalam transaksi.", type: "error" });
      return;
    }

    const user = JSON.parse(localStorage.getItem("user") || "{}");
    const shiftId = localStorage.getItem("shiftregularId");

    const totalDiscount = form.total_discount ? Number(form.total_discount) : 0;
    const totalPay = subTotal - totalDiscount;

    const payload = {
      transaction_date: new Date().toISOString(),
      cashier_name: user.full_name || "Unknown",
      customer_name: form.customer_name || null,
      customer_contact: form.customer_contact || null,
      description: form.description || null,
      sub_total: subTotal,
      total_discount: form.total_discount ? totalDiscount : null,
      total_pay: totalPay,
      payment_method: form.payment_method,
      shift_id: shiftId ? Number(shiftId) : null,
      items: regularList,
    };

    try {
      await apiClient.post("/sales/regular", payload);
      
      const prevBalance = Number(localStorage.getItem("closing_balance") || "0");
      const newBalance = prevBalance + totalPay;
      localStorage.setItem("closing_balance", newBalance.toString());

      const prevSalesCount = Number(localStorage.getItem("regular_sale") || "0");
      const newSalesCount = prevSalesCount + 1;
      localStorage.setItem("regular_sale", newSalesCount.toString());

      onAfterSubmit();

      localStorage.removeItem("regularList");
      onClose();
      setToast({ message: "Transaksi berhasil disimpan.", type: "success" });
    } catch (err) {
      console.error(err);
      setToast({ message: "Gagal menyimpan transaksi.", type: "error" });
    }
  };

  return (
    <Modal isOpen={isOpen} close={onClose}>
      <h2 className="text-xl font-semibold text-center mb-4">
        Akhiri Transaksi
      </h2>

      <div className="space-y-3 max-h-[60vh] overflow-y-auto">
        <InputField
          label="Nama Customer"
          value={form.customer_name}
          onChange={handleChange("customer_name")}
          placeholder="Opsional"
        />
        <InputField
          label="Kontak Customer"
          value={form.customer_contact}
          onChange={handleChange("customer_contact")}
          placeholder="Opsional"
          type="phone"
        />
        <TextArea
          label="Catatan Tambahan"
          value={form.description}
          onChange={handleChange("description")}
          placeholder="Catatan Tambahan (Opsional)"
        />
        <InputField
          label="Diskon"
          value={form.total_discount}
          onChange={handleChange("total_discount")}
          placeholder="Opsional"
          type="number"
        />
        <div>
          <label className="block font-medium mb-1">Metode Pembayaran</label>
          <Select
            value={form.payment_method}
            onChange={handleChange("payment_method")}
            className="border rounded p-2 w-full"
          >
            <option>Cash</option>
            <option>QRIS</option>
            <option>Debit</option>
            <option>Kredit</option>
            <option>Transfer</option>
          </Select>
        </div>

        <div className="font-semibold text-lg text-right p-5">
          Subtotal: Rp. {subTotal.toLocaleString("id-ID")}
        </div>
      </div>

      {toast && (
        <Toast
          message={toast.message}
          type={toast.type}
          onClose={() => setToast(null)}
        />
      )}

      <div className="mt-6 flex gap-4">
        <button onClick={onClose} className="text-black w-full bg-gray-200 border border-black hover:bg-gray-300 rounded-md">
          Batal
        </button>
        <Button onClick={handleConfirm} className="w-full">
          Konfirmasi
        </Button>
      </div>
    </Modal>
  );
}
