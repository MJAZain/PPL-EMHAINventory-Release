import React from "react";

const Select = ({ error = false, className = "", children, ...props }) => {
  const isError = Boolean(error);
  const errorMessage = typeof error === "string" ? error : null;

  const combinedClassName = `
    border ${isError ? "border-red-500" : "border-black"}
    rounded-[6px]
    px-[20px] 
    py-2 
    text-base
    placeholder-black
    w-full
    focus:outline-none 
    focus:ring-1 ${isError ? "focus:ring-red-500" : "focus:ring-blue-500"}
    ${className}
  `;

  return (
    <div>
      <select className={combinedClassName} {...props}>
        {children}
      </select>
      {errorMessage && (
        <p className="mt-1 text-xs text-red-500">{errorMessage}</p>
      )}
    </div>
  );
};

export default Select;
