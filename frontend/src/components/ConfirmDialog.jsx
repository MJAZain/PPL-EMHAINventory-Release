import React from "react";

function ConfirmDialog({ isOpen, title, description, onCancel, onConfirm }) {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm">
      <div className="bg-white rounded-xl p-6 w-full max-w-md shadow-xl">
        <h2 className="text-xl font-semibold text-gray-900 mb-2 text-center">{title}</h2>
        <p className="text-gray-700 mb-6 text-center">{description}</p>
        <div className="flex justify-center gap-3">
          <button
            onClick={onCancel}
            className="px-4 py-2 rounded-md bg-gray-100 text-black hover:bg-gray-200 w-full"
          >
            Batal
          </button>
          <button
            onClick={onConfirm}
            className="px-4 py-2 rounded-md bg-red-600 text-white hover:bg-red-700 w-full"
          >
            Hapus
          </button>
        </div>
      </div>
    </div>
  );
}

export default ConfirmDialog;
