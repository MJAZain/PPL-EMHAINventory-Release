import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { publicClient } from "../config/api";

const useLogin = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  const navigate = useNavigate();

const login = async () => {
  setLoading(true);
  setError(null);
  setSuccess(null);

  try {
  const { data } = await publicClient.post("/users/login", { email, password });
  console.log("Response:", data);

    if (!data) {
      throw new Error("Login gagal");
    }

  const token = data.Data?.token || data.data?.token || data.token;
    if (token) {
      localStorage.setItem("token", token);
    } else {
      throw new Error("No token received");
    }

    setSuccess("Login successful!");
    navigate("/dashboard");
  } catch (err) {
    let message = "Terjadi kesalahan saat login. Silakan coba lagi.";

    if (err.response) {
      const status = err.response.status;
      const serverMessage = err.response.data?.message?.toLowerCase?.();

      if (status === 400 && (!email || !password)) {
        message = "Email dan Password harus diisi.";
      } else if (status === 401 || serverMessage?.includes("invalid")) {
        message = "Email atau Password salah.";
      } else if (status === 403) {
        message = "Akun Anda tidak memiliki izin untuk masuk.";
      } else if (status >= 500) {
        message = "Server sedang bermasalah. Silakan coba beberapa saat lagi.";
      } else {
        message = "Login gagal. Periksa kembali data Anda.";
      }
    } else if (err.message === "NO_TOKEN") {
      message = "Token tidak diterima. Hubungi administrator.";
    } else if (err.message === "DATA_NULL") {
      message = "Data login tidak ditemukan.";
    } else {
      message = "Tidak dapat terhubung ke server. Periksa koneksi internet Anda.";
    }

    setError(message);
  } finally {
    setLoading(false);
  }
};

return {
  email,
  setEmail,
  password,
  setPassword,
  loading,
  error,
  success,
  login,
};

};

export default useLogin;
