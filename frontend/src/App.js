import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import LoginPage from "./pages/LoginPage";
import MasterBarangPage from "./pages/MasterBarangPage";
import BarangMasukPage from './pages/BarangMasukPage';
import RegisterUserPage from "./pages/RegisterPage";
import AturSatuanPage from "./pages/AturSatuanPage";
import AturKategoriPage from "./pages/AturKategoriPage";
import AturUsersPage from "./pages/UserManagementPage";
import BarangMasukDetailPage from "./pages/BarangMasukDetail";
import BarangKeluarPage from "./pages/BarangKeluarPage";
import BarangKeluarDetailPage from "./pages/BarangKeluarDetail";
import StockOpnamePage from "./pages/StockOpname/StockOpname";
import DraftDetailPage from "./pages/StockOpname/DraftDetailPage";
import LaporanBarangMasukPage from "./pages/LaporanKeuangan/LaporanBarangMasuk";
import LaporanBarangKeluarPage from "./pages/LaporanKeuangan/LaporanBarangKeluar";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<LoginPage />} />
        <Route path="/register" element={<RegisterUserPage />} />
        <Route path="/dashboard" element={<MasterBarangPage />} />
        <Route path="/barang-masuk" element={<BarangMasukPage />} />
        <Route path="/barang-masuk-detail" element={<BarangMasukDetailPage />} />
        <Route path="/barang-keluar" element={<BarangKeluarPage />} />
        <Route path="/barang-keluar-detail" element={<BarangKeluarDetailPage />} />
        <Route path="/satuan" element={<AturSatuanPage />} />
        <Route path="/kategori" element={<AturKategoriPage />} />
        <Route path="/user" element={<AturUsersPage />} />
        <Route path="/stock-opname" element={<StockOpnamePage />} />
        <Route path="/stock-opname/draft/:draftId" element={<DraftDetailPage />} />
        <Route path="/laporan-masuk" element={<LaporanBarangMasukPage />} />
        <Route path="/laporan-terjual" element={<LaporanBarangKeluarPage />} />
      </Routes>
    </Router>
  );
}

export default App;
