import React, { useState, useEffect } from 'react';
import { apiClient } from '../../config/api';
import InputField from '../../components/inputField';
import Modal from '../../components/modal/modal';
import Button from '../../components/buttonComp';
import { getFriendlyErrorMessage } from '../../utils/errorHandler';
import Toast from '../../components/toast';

const formFields = [
  { label: "Tanggal Stock Opname", key: "opname_date", placeholder: "Pilih Tanggal Stock Opname" },
  { label: "Jenis Stock Opname", key: "jenis_stock_opname", type: "select", optionsKey: "jenis" },
  { label: "Catatan Stock Opname", key: "notes", placeholder: "Catatan Stock Opname" },
];

const CreateDraftModal = ({ isOpen, onClose, onDraftCreated }) => {
  const [opnameDate, setOpnameDate] = useState(new Date().toISOString().substring(0, 10));
  const [notes, setNotes] = useState('');
  const [loading, setLoading] = useState(false);
  const [toast, setToast] = useState(null);

  const [jenis, setJenis] = useState([
    { id: 'reguler', name: 'Reguler' },
    { id: 'harian', name: 'Harian' }
  ]);

  const [form, setForm] = useState({
  opname_date: opnameDate,
  notes: '',
  jenis_stock_opname: '',
});

const handleChange = (key) => (e) => {
  const value = e.target.value;
  setForm((prev) => ({ ...prev, [key]: value }));
};

  const createDraft = async () => {
    if (!opnameDate) return;
    try {
      setLoading(true);
      await apiClient.post('/stock-opname/draft', {
  
        notes: form.notes,
        jenis_stock_opname: form.jenis_stock_opname
      });
      onDraftCreated();
      onClose();
    } catch (error) {
      const message = getFriendlyErrorMessage(error);
      setToast({ message, type: "error" });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal isOpen={isOpen} close={onClose}>

      <h2 className="text-xl font-semibold mb-4 text-center py-5">Buat Draft Stock Opname</h2>

      <div className="max-h-[60vh] overflow-y-auto pr-2">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 w-full">
          {formFields.map(({ label, key, placeholder, type, optionsKey }) => {
            const isTextArea = key.includes("notes");

            if (type === "select") {
              const optionsRaw = {
                jenis,
              }[optionsKey];
              const options = Array.isArray(optionsRaw) ? optionsRaw : [];

              return (
                <div key={key} className="flex flex-col">
                  <label className="text-sm font-medium mb-1">{label}</label>
                  <select
                    value={form[key]}
                    onChange={handleChange(key)}
                    className="border border-gray-300 rounded-md px-3 py-2"
                  >
                    <option value="">Pilih {label}</option>
                    {options.map((opt) => (
                      <option key={opt.id} value={opt.id}>
                        {opt.name}
                      </option>
                    ))}
                  </select>
                </div>
              );
            }

            if (isTextArea) {
              return (
                <div key={key} className="flex flex-col col-span-full">
                  <label className="text-sm font-medium mb-1">{label}</label>
                  <textarea
                    value={form[key]}
                    onChange={handleChange(key)}
                    placeholder={placeholder}
                    rows={4}
                    className="border border-gray-300 rounded-md px-3 py-2 resize-none"
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
        <button onClick={onClose} className='w-1/4 bg-gray-200 border border-black text-black rounded-md py-2 hover:bg-gray-300 transition'>
          Batal
        </button>
        <Button onClick={createDraft} disabled={loading || !opnameDate}>
          {loading ? 'Membuat...' : 'Buat Draft'}
        </Button>
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
};

export default CreateDraftModal;
