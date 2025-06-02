import React from "react";

function Modal({
  isOpen,
  close,
  onSubmit,
  children,
  containerClassName = "",
  contentClassName = "",
}) {
  if (!isOpen) return null;

  return (
    <div
      className={`fixed inset-0 bg-gray-500 bg-opacity-50 flex justify-center items-center z-50 ${containerClassName}`}
    >
      <div
        className={`relative bg-white p-6 rounded-lg shadow-lg w-[90%] max-w-3xl ${contentClassName}`}
        onClick={(e) => e.stopPropagation()}
      >
        <button
          className="absolute top-4 right-4 text-gray-500"
          onClick={close}
        >
          &times;
        </button>
        {children}
      </div>
    </div>
  );
}

export default Modal;
