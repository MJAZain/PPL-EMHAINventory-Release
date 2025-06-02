import React, { useState } from 'react';
import { apiClient } from '../../config/api';
import Modal from '../../components/modal/modal';
import Button from '../../components/buttonComp';

const CreateDraftModal = ({ isOpen, onClose, onDraftCreated }) => {
  const [opnameDate, setOpnameDate] = useState(new Date().toISOString().substring(0, 10));
  const [notes, setNotes] = useState('');
  const [loading, setLoading] = useState(false);

  const createDraft = async () => {
    if (!opnameDate) return;
    try {
      setLoading(true);
      await apiClient.post('/stock-opname/draft', {
        opname_date: new Date(opnameDate).toISOString(),
        notes,
      });
      onDraftCreated();
      onClose();
    } catch (error) {
      console.error('Error creating draft:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal isOpen={isOpen} close={onClose}>
      <h2 className="text-xl font-semibold mb-4">Buat Draft Stock Opname</h2>

      <label className="block mb-2 font-medium">Tanggal Stock Opname</label>
      <input
        type="date"
        value={opnameDate}
        onChange={(e) => setOpnameDate(e.target.value)}
        className="w-full border-gray-300 p-2  mb-4
        border border-[var(--neutral-400,#A1A1A1)]
            bg-[var(--neutral-200,#E5E5E5)]
            rounded-[6px]
            px-[20px]
            text-base
            placeholder-gray-500"
      />

      <label className="block mb-2 font-medium">Catatan (Opsional)</label>
      <textarea
        value={notes}
        onChange={(e) => setNotes(e.target.value)}
        placeholder="Tulis catatan..."
        className="w-full border-gray-300 p-2 mb-4
        border border-[var(--neutral-400,#A1A1A1)]
        bg-[var(--neutral-200,#E5E5E5)]
            rounded-[6px]
            px-[20px]
            text-base
            placeholder-gray-500"
      />

      <div className="flex justify-end gap-3">
        <Button onClick={onClose} variant="secondary">
          Batal
        </Button>
        <Button onClick={createDraft} disabled={loading || !opnameDate}>
          {loading ? 'Membuat...' : 'Buat Draft'}
        </Button>
      </div>
    </Modal>
  );
};

export default CreateDraftModal;
