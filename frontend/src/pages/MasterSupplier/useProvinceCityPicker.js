import { useState, useEffect } from "react";
import provinces from '../../data/provinces.json';
import cities from '../../data/regencies.json';

export function useProvinceCityPicker(initialProvinceId = null) {
  const [selectedProvinceId, setSelectedProvinceId] = useState(initialProvinceId);
  const [cityOptions, setCityOptions] = useState([]);

  useEffect(() => {
    if (selectedProvinceId) {
      const filteredCities = cities.filter(
        (city) => city.province_id === selectedProvinceId
      );
      setCityOptions(filteredCities);
    } else {
      setCityOptions([]);
    }
  }, [selectedProvinceId]);

  return {
    provinces,
    cityOptions,
    selectedProvinceId,
    setSelectedProvinceId,
  };
}
