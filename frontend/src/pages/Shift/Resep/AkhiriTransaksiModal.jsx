import React, { useState, useEffect } from "react";
import Modal from "../../../components/modal/modal";
import InputField from "../../../components/inputField";
import Button from "../../../components/buttonComp";
import Toast from "../../../components/toast";
import { apiClient } from "../../../config/api";
import Select from "../../../components/SelectComp";
import TextArea from '../../../components/textareacomp';

export default function AkhiriTransaksiModal({
  isOpen,
  onClose,
  presList,
  onAfterSubmit,
}) {
  const [form, setForm] = useState({
    doctor_id: null,
    doctorSearch: "",
    clinic: "",
    patient_id: null,
    patientSearch: "",
    description: "",
    total_discount: "",
    payment_method: "Cash",
  });

  const [doctors, setDoctors] = useState([]);
  const [patients, setPatients] = useState([]);
  const [subTotal, setSubTotal] = useState(0);
  const [toast, setToast] = useState(null);

  useEffect(() => {
    if (isOpen) {
      fetchDoctors();
      fetchPatients();
      const presList = JSON.parse(localStorage.getItem("presList") || "[]");
      const sum = presList.reduce((acc, item) => acc + item.price, 0);
      setSubTotal(sum);
    }
  }, [isOpen]);

  const fetchDoctors = async () => {
    try {
      const res = await apiClient.get("/doctors/");
      setDoctors(res.data?.data || []);
    } catch (err) {
      console.error("Failed to fetch doctors", err);
    }
  };

  const fetchPatients = async () => {
    try {
      const res = await apiClient.get("/patients/");
      setPatients(res.data?.data || []);
    } catch (err) {
      console.error("Failed to fetch patients", err);
    }
  };

  const handleChange = (key) => (e) => {
    setForm({ ...form, [key]: e.target.value });
  };

  const handleSelectDoctor = (doctor) => {
    setForm((prev) => ({
      ...prev,
      doctor_id: doctor.id,
      clinic: doctor.practice_address,
      doctorSearch: doctor.full_name,
    }));
  };

  const handleSelectPatient = (patient) => {
    setForm((prev) => ({
      ...prev,
      patient_id: patient.id,
      patientSearch: patient.full_name,
    }));
  };

  const handleConfirm = async () => {
    const presList = JSON.parse(localStorage.getItem("presList") || "[]");
    if (presList.length === 0) {
      setToast({ message: "Tidak ada item dalam transaksi.", type: "error" });
      return;
    }
    if (!form.doctor_id || !form.patient_id) {
      setToast({ message: "Dokter dan pasien wajib dipilih.", type: "error" });
      return;
    }

    const shiftId = localStorage.getItem("presId");
    const totalDiscount = form.total_discount ? Number(form.total_discount) : 0;

    const payload = {
      prescription_no: `RX-${Date.now()}`,
      prescription_date: new Date().toISOString(),
      doctor_id: form.doctor_id,
      clinic: form.clinic || "",
      diagnosis: form.diagnosis || "",
      patient_id: form.patient_id,
      transaction_date: new Date().toISOString(),
      payment_method: form.payment_method,
      discount_percent: 0,
      discount_amount: totalDiscount,
      shift_id: shiftId ? Number(shiftId) : null,
      items: presList,
    };

    try {
      await apiClient.post("/sales/prescriptions", payload);

      const prevBalance = Number(localStorage.getItem("closing_balance_pres") || "0");
      const newBalance = prevBalance + (subTotal - totalDiscount);
      localStorage.setItem("closing_balance_pres", newBalance.toString());

      const prevSalesCount = Number(localStorage.getItem("pres_sale") || "0");
      const newSalesCount = prevSalesCount + 1;
      localStorage.setItem("pres_sale", newSalesCount.toString());

      onAfterSubmit();
      localStorage.removeItem("presList");
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
        {/* Search Doctor */}
        <InputField
          label="Cari Dokter"
          value={form.doctorSearch}
          onChange={handleChange("doctorSearch")}
          placeholder="Nama Dokter"
        />
        <ul className="border rounded p-2 max-h-32 overflow-y-auto">
          {doctors
            .filter((doc) =>
              (doc.full_name || "").toLowerCase().includes(form.doctorSearch.toLowerCase())
            )
            .map((doc) => (
              <li
                key={doc.id}
                className="p-2 hover:bg-gray-100 cursor-pointer border-b"
                onClick={() => handleSelectDoctor(doc)}
              >
                <div className="font-medium">{doc.full_name}</div>
                <div className="text-sm text-gray-600">
                  {doc.practice_address}
                </div>
              </li>
            ))}
        </ul>
        <InputField
          label="Klinik"
          value={form.clinic || ""}
          disabled
        />

        {/* Search Patient */}
        <div>
          <InputField
            label="Cari Pasien"
            value={form.patientSearch}
            onChange={(e) =>
              setForm((prev) => ({
                ...prev,
                patientSearch: e.target.value,
              }))
            }
            placeholder="Nama pasien"
          />
          <ul className="border p-2 rounded mb-4 max-h-40 overflow-y-auto">
            {patients
              .filter((p) =>
                p.full_name
                  .toLowerCase()
                  .includes(form.patientSearch.toLowerCase())
              )
              .map((patient) => (
                <li
                  key={patient.id}
                  className="p-2 hover:bg-gray-100 cursor-pointer border-b"
                  onClick={() => handleSelectPatient(patient)}
                >
                  {patient.full_name}
                </li>
              ))}
          </ul>
        </div>
        <TextArea
          label="Diagnosis"
          value={form.diagnosis}
          onChange={handleChange("diagnosis")}
          placeholder="Diagnosis pasien"
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
        <button
          onClick={onClose}
          className="text-black w-full bg-gray-200 border border-black hover:bg-gray-300 rounded-md"
        >
          Batal
        </button>
        <Button onClick={handleConfirm} className="w-full">
          Konfirmasi
        </Button>
      </div>
    </Modal>
  );
}
