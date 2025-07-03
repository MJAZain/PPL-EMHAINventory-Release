import React, { useState, useEffect } from "react";
import InputField from "../components/inputField";
import Button from "../components/buttonComp";
import useLogin from "../hooks/useUser";
import myImage from '../media/loginDecor.png';
import '../styles/loginStyles.css';
import Toast from "../components/toast";
import { useLocation } from "react-router-dom";

function LoginPage() {
  const {
    email,
    setEmail,
    password,
    setPassword,
    loading,
    error,
    success,
    login,
  } = useLogin();
  
  const [toast, setToast] = useState(null);
  const location = useLocation();

  useEffect(() => {
    if (success) {
      setToast({ message: success, type: "success" });
    }
    if (error) {
      setToast({ message: error, type: "error" });
    }
  }, [success, error]);

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const errorParam = params.get("error");

    if (errorParam === "inactive") {
      setToast({
        message: "Akun Anda telah dinonaktifkan.",
        type: "error",
      });
    }
  }, [location.search]);

  const handleLogin = async () => {
    await login();
  };

  return (
    <div className="loginPage">
      {toast && (
        <Toast
          message={toast.message}
          type={toast.type}
          onClose={() => setToast(null)}
        />
      )}

      <div className="loginDecor">
        <h1 className="loginTitle">EIMS</h1>
        <img src={myImage} alt="Description" className="loginImg" />
      </div>

      <div className="loginContainer">
        <h2 className="loginLogin">Login</h2>

        <InputField
          label="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="Masukkan email"
          className="loginField"
        />

        <InputField
          label="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          type="password"
          placeholder="Masukkan password"
          className="loginField"
        />

        <Button onClick={handleLogin} disabled={loading} className="loginButton">
          {loading ? "Masuk..." : "Masuk"}
        </Button>
      </div>
    </div>
  );
}

export default LoginPage;
