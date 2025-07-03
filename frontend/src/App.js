import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import LoginPage from "./pages/LoginPage";
import MasterBarangPage from "./pages/MasterBarangPage";
import BarangMasukPage from "./pages/BarangMasukPage";
import RegisterUserPage from "./pages/RegisterPage";
import AturSatuanPage from "./pages/MasterSatuan/AturSatuanPage";
import AturKategoriPage from "./pages/MasterKategori/AturKategoriPage";
import AturUsersPage from "./pages/UserManagementPage";
import BarangMasukDetailPage from "./pages/BarangMasukDetail";
import BarangKeluarPage from "./pages/BarangKeluarPage";
import BarangKeluarDetailPage from "./pages/BarangKeluarDetail";
import StockOpnamePage from "./pages/StockOpname/StockOpname";
import DraftDetailPage from "./pages/StockOpname/DraftDetailPage";

import StorageLocationPage from "./pages/MasterData/StorageLocationPage";
import BrandPage from "./pages/MasterData/BrandPage";
import AturPatientsPage from "./pages/MasterPasien/AturPatientPage";
import AturSuppliersPage from "./pages/MasterSupplier/AturSupplierPage";
import AturDoctorPage from "./pages/MasterDokter/AturDoctorPage";
import AturGolonganObatPage from "./pages/MasterGolonganObat/AturGolonganPage";

import NonPBFDetailPage from "./pages/BarangMasukNon-PBF/NonPBFDetailPage";
import NonPBFProductListPage from "./pages/BarangMasukNon-PBF/NonPBFProductList";

import PBFDetailPage from "./pages/BarangMasukPBF/PBFDetailPage";
import PBFProductListPage from "./pages/BarangMasukPBF/PBFProductList";

import RiwayatPBFPage from "./pages/RiwayatTransaksi/RiwayatPBF/RiwayatPBFPage";
import RiwayatNonPBFPage from "./pages/RiwayatTransaksi/RiwayatNonPbf/RiwayatNonPBFPage";

import ResepShiftUmumPage from "./pages/Shift/ResepShiftUmumPage";
import TanpaResepShiftUmumPage from "./pages/Shift/TanpaResepShiftUmumPage";

import KoreksiFormPage from "./pages/KoreksiStok/KoreksiFormPage";

import AturKaryawanPage from "./pages/MasterKaryawan/AturKaryawanPage";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<LoginPage />} />
        <Route path="/user" element={<AturKaryawanPage />} />

        <Route path="/doctor" element={<AturDoctorPage />} />
        <Route path="/supplier" element={<AturSuppliersPage />} />
        <Route path="/patients" element={<AturPatientsPage />} />
        <Route path="/register" element={<RegisterUserPage />} />
        <Route path="/dashboard" element={<MasterBarangPage />} />
        <Route path="/satuan" element={<AturSatuanPage />} />
        <Route path="/kategori" element={<AturKategoriPage />} />
        <Route path="/golongan" element={<AturGolonganObatPage />} />
        <Route path="/storage-locations" element={<StorageLocationPage />} />
        <Route path="/brands" element={<BrandPage />} />

        <Route path="/pbf-detail" element={<PBFDetailPage />} />         
        <Route path="/pbf-detail/:id" element={<PBFDetailPage />} />    
        <Route path="/pbf-list" element={<PBFProductListPage />} />
        <Route path="/pbf-list/:id" element={<PBFProductListPage />} />

        <Route path="/non-pbf-detail" element={<NonPBFDetailPage />} />
        <Route path="/non-pbf-detail/:id" element={<NonPBFDetailPage />} />
        <Route path="/non-pbf-list" element={<NonPBFProductListPage />} />
        <Route path="/non-pbf-list/:id" element={<NonPBFProductListPage />} />

        <Route path="/shift-resep" element={<ResepShiftUmumPage />} />
        <Route path="/shift-tanpa-resep" element={<TanpaResepShiftUmumPage />} />

        <Route path="/stock-opname" element={<StockOpnamePage />} />
        <Route
          path="/stock-opname/draft/:draftId"
          element={<DraftDetailPage />}
        />

        <Route path="/koreksi" element={<KoreksiFormPage />} />

        <Route path="/riwayat-pbf" element={<RiwayatPBFPage />} />
        <Route path="/riwayat-non-pbf" element={<RiwayatNonPBFPage />} />
      </Routes>
    </Router>
  );
}

export default App;
