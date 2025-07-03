import React, { useState } from "react";
import useBrandActions from "../../hooks/useBrandActions";
import InputField from "../inputField";

import Toast from "../toast";
import {getFriendlyErrorMessage} from "../../utils/errorHandler"

const AddBrandModal = ({ isOpen, close, onSuccess }) => {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [errors, setErrors] = useState({});
  const [isLoading, setIsLoading] = useState(false);
  const { addBrand } = useBrandActions();

  const [toast, setToast] = useState(null);

  const validateForm = () => {
    const newErrors = {};
    if (!name.trim()) {
      newErrors.name = "Nama brand tidak boleh kosong.";
    }
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!validateForm()) {
      return;
    }
    setIsLoading(true);
    try {
      await addBrand({ name, description });
      onSuccess({ message: "Brand berhasil ditambahkan." });
      setName("");
      setDescription("");
      setErrors({});
    } catch (error) {
      const message = getFriendlyErrorMessage(error);
      setToast({ message, type: "error" });
    } finally {
      setIsLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50 flex justify-center items-center">
      <div className="bg-white p-6 rounded-lg shadow-xl w-full max-w-md">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold">Tambah Brand Baru</h2>
          <button
            onClick={close}
            className="text-gray-700 hover:text-gray-900 text-2xl"
          >
            &times;
          </button>
        </div>
        <form onSubmit={handleSubmit}>
          <InputField
            label="Nama Brand"
            id="brand-name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            error={errors.name}
            isRequired={true}
            placeholder="Contoh: Brand A"
          />
          <InputField
            label="Deskripsi Brand"
            id="brand-description"
            type="textarea"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            error={errors.description}
            placeholder="Deskripsi singkat brand (opsional)"
          />
          {errors.form && (
            <p className="mt-1 text-xs text-red-500 mb-2">{errors.form}</p>
          )}
          <div className="flex justify-end space-x-3 mt-6">
            <button
              type="button"
              onClick={close}
              disabled={isLoading}
              className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 disabled:opacity-50"
            >
              Batal
            </button>
            <button
              type="submit"
              disabled={isLoading}
              className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
            >
              {isLoading ? "Menyimpan..." : "Simpan"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default AddBrandModal;
