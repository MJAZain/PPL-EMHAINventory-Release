import React from "react";
import InputField from "./inputField";

function SearchBar({ value, onChange, placeholder = "Search..." }) {
  return (
    <div className="max-w-md">
      <InputField
        label=""
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        className="w-full h-10"
      />
    </div>
  );
}

export default SearchBar;
