import React from 'react';

function InputField({
  label,
  value,
  onChange,
  type = "text",
  placeholder,
  className = "",
}) {
  const isNumeric = type === "number";

  const handleKeyDown = (e) => {
    if (isNumeric) {
      const invalidChars = ["e", "E", "+", "-", "."];
      if (invalidChars.includes(e.key)) {
        e.preventDefault();
      }
    }
  };

  const handleInput = (e) => {
    if (isNumeric) {
      let inputValue = e.target.value;

      // Prevent leading zero (unless user is typing "0." which we allow)
      if (/^0\d+/.test(inputValue)) {
        inputValue = inputValue.replace(/^0+/, '');
        e.target.value = inputValue;
        // Manually call onChange to update value in parent
        onChange({ target: { value: inputValue } });
      }
    }
  };

  const handleWheel = (e) => {
    if (isNumeric) {
      e.target.blur(); // Prevent value from changing on scroll
    }
  };

  return (
    <div className="mb-4">
      <label className="block mb-1 text-sm font-medium text-gray-700">{label}</label>
      <input
        type={type}
        value={value}
        onChange={onChange}
        onKeyDown={handleKeyDown}
        onInput={handleInput}
        onWheel={handleWheel}
        min={isNumeric ? 1 : undefined}
        placeholder={placeholder}
        className={`
          border border-[var(--neutral-400,#A1A1A1)]
          bg-[var(--neutral-200,#E5E5E5)]
          rounded-[6px]
          px-[20px]
          text-base
          placeholder-gray-500
          ${className}
        `}
      />
    </div>
  );
}

export default InputField;
