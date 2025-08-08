import React, { useEffect, useState } from "react";
import Modal from "../../../components/modal/modal";
import InputField from "../../../components/inputField";
import Button from "../../../components/buttonComp";
import Toast from "../../../components/toast";
import { apiClient } from "../../../config/api";
import { useNavigate } from "react-router-dom";
import TextArea from "../../../components/textareacomp";
import Select from "../../../components/SelectComp";

export default function AkhiriShiftModal({ isOpen, onClose }) {
  const navigate = useNavigate();

  const [closingBalance, setClosingBalance] = useState(0);
  const [shiftId, setShiftId] = useState(null);
  const [toast, setToast] = useState(null);

  const [correctionType, setCorrectionType] = useState("Tambah");
  const [correctionAmount, setCorrectionAmount] = useState("");
  const [note, setNote] = useState("");

  useEffect(() => {
    if (isOpen) {
      const balance = Number(localStorage.getItem("closing_balance_pres") || "0");
      const shift = localStorage.getItem("presId");

      setClosingBalance(balance);
      setShiftId(shift);
    }
  }, [isOpen]);

  const handleConfirm = async () => {
    if (!shiftId) {
      setToast({ message: "Shift ID tidak ditemukan.", type: "error" });
      return;
    }

    let finalKas = closingBalance;
    if (correctionAmount) {
      const amount = Number(correctionAmount);
      if (correctionType === "Tambah") {
        finalKas += amount;
      } else if (correctionType === "Kurangi") {
        finalKas -= amount;
      }
    }

    const totalSales = Number(localStorage.getItem("pres_sale") || "0");

    const payload = {
      closing_balance: finalKas,
      note: note || null,
      total_sales: totalSales,
    };

    try {
      await apiClient.put(`/shifts/close/${shiftId}`, payload);

      localStorage.removeItem("presId");
      localStorage.removeItem("closing_balance_pres");
      localStorage.removeItem("pres_sale");

      onClose();
      navigate("/shift-resep");
    } catch (err) {
      console.error(err);
      setToast({ message: "Gagal mengakhiri shift.", type: "error" });
    }
  };

  return (
    <Modal isOpen={isOpen} close={onClose}>
      <h2 className="text-xl font-semibold text-center mb-4">
        Akhiri Shift
      </h2>

      <div className="space-y-4 max-h-[60vh] overflow-y-auto">
        <InputField
          label="Kas Penutup"
          value={closingBalance.toLocaleString("id-ID")}
          disabled
        />

        <div>
          <label className="block font-medium mb-1">Koreksi Kas</label>
          <div className="flex gap-2">
            <Select
              value={correctionType}
              onChange={(e) => setCorrectionType(e.target.value)}
              className="border rounded p-2"
            >
              <option>Tambah</option>
              <option>Kurangi</option>
            </Select>
            <input
              type="number"
              value={correctionAmount}
              onChange={(e) => setCorrectionAmount(e.target.value)}
              placeholder="0"
              className="border border-black rounded p-2 flex-1"
            />
          </div>
        </div>

        <div>
          <label className="block font-medium mb-1">Catatan</label>
          <TextArea
            value={note}
            onChange={(e) => setNote(e.target.value)}
            rows={3}
            placeholder="Opsional"
            className="border rounded p-2 w-full"
          />
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
