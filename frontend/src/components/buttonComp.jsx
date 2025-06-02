import React from "react";

function Button({ onClick, children, className = "" }) {
  return (
    <button
      onClick={onClick}
      className={`
        bg-[#2B7FFF] 
        text-white 
        py-2 px-4 
        rounded-md
        ${className}
      `}
    >
      {children}
    </button>
  );
}

export default Button;
