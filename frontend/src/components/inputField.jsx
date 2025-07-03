import React from "react";

function InputField({
  label,
  value,
  onChange, 
  type = "text",
  placeholder,
  className = "",
  id,
  isRequired,
  error,
  min,
}) {
  const isNumeric = type === "number";
  const isPhoneField = type === 'phone';

  const handleKeyDown = (e) => {
    // Block unwanted characters for number fields
    if (isNumeric) {
      const invalidChars = ["e", "E", "+", "-", "."];
      if (invalidChars.includes(e.key)) {
        e.preventDefault();
      }
    }

    // Allow only digits for phone fields
    if (isPhoneField && !/^\d$/.test(e.key) && e.key !== "Backspace" && e.key !== "Tab") {
      e.preventDefault();
    }
  };

  const handleInput = (e) => {
    let inputValue = e.target.value;

    // For phone fields: strip out anything that's not a digit
    if (isPhoneField) {
      inputValue = inputValue.replace(/\D/g, ""); // Remove non-digit characters
      e.target.value = inputValue;
    }

    // For number fields: remove leading zeroes
    if (isNumeric) {
      inputValue = inputValue.replace(/^0+/, "");
      e.target.value = inputValue;
    }

    if (onChange) {
      onChange({ target: { name: e.target.name, value: inputValue } });
    }
  };

  const handleWheel = (e) => {
    if (isNumeric) {
      e.target.blur();
    }
  };

  return (
    <div className="mb-4">
      {label && (
        <label
          htmlFor={id}
          className="block mb-1 text-sm font-medium text-gray-700"
        >
          {label}
          {isRequired && <span className="text-red-500 ml-1">*</span>}
        </label>
      )}
      <input
        id={id}
        name={id || label?.toLowerCase().replace(/\s+/g, "-")}
        type={type}
        value={value}
        onChange={onChange}
        onKeyDown={handleKeyDown}
        onInput={handleInput}
        onWheel={handleWheel}
        min={min || (isNumeric ? 1 : undefined)}
        placeholder={placeholder}
        className={`
          border ${
            error ? "border-red-500" : "border-black"
          }
          rounded-[6px]
          px-[20px] 
          py-2 
          text-base
          placeholder-black
          w-full
          focus:outline-none focus:ring-1 ${
            error ? "focus:ring-red-500" : "focus:ring-blue-500"
          }
          ${className}
        `}
      />
      {error && <p className="mt-1 text-xs text-red-500">{error}</p>}
    </div>
  );
}

export default InputField;
