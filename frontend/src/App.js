import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import LoginPage from "./pages/LoginPage";
import Dashboard from "./pages/Dashboard/DashboardPage";
import AnalyticsDashboard from "./pages/Laporan/LaporanPenjualan";
import ProfitAnalysisPage from "./pages/Laporan/LaporanLabaRugi";

import MasterBarangPage from "./pages/MasterBarangPage";
import RegisterUserPage from "./pages/RegisterPage";
import AturSatuanPage from "./pages/MasterSatuan/AturSatuanPage";
import AturKategoriPage from "./pages/MasterKategori/AturKategoriPage";

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

import TanpaResepShiftUmumPage from "./pages/Shift/TanpaResep/TanpaResepShiftUmumPage";
import TanpaResepDetailPage from "./pages/Shift/TanpaResep/TanpaResepDetail";

import ResepShiftUmumPage from "./pages/Shift/Resep/ResepShiftUmumPage";
import ResepDetailPage from "./pages/Shift/Resep/ResepDetail";

import RiwayatShiftPage from "./pages/Shift/RiwayatShiftPage";
import RiwayatRegularPage from "./pages/Shift/Riwayat/RiwayatRegular";
import RiwayatPresPage from "./pages/Shift/Riwayat/RiwayatPres";

import KoreksiFormPage from "./pages/KoreksiStok/KoreksiFormPage";

import AturKaryawanPage from "./pages/MasterKaryawan/AturKaryawanPage";

import AturJenisPage from "./pages/MasterJenisPengeluaran/AturJenisPage";
import AturPengeluaranPage from "./pages/MasterPengeluaran/AturPengeluaranPage";

import OpnameDraftPage from './pages/StockOpname/OpnameDraftPage'
import OpnameDetailPage from "./pages/StockOpname/OpnameDetailPage";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<LoginPage />} />
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/sales-analytics" element={<AnalyticsDashboard />} />
        <Route path="/profit" element={<ProfitAnalysisPage />} />
        <Route path="/user" element={<AturKaryawanPage />} />

        <Route path="/doctor" element={<AturDoctorPage />} />
        <Route path="/supplier" element={<AturSuppliersPage />} />
        <Route path="/patients" element={<AturPatientsPage />} />
        <Route path="/register" element={<RegisterUserPage />} />
        <Route path="/master-obat" element={<MasterBarangPage />} />
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
        <Route path="/resep-detail" element={<ResepDetailPage />} />

        <Route path="/shift-tanpa-resep" element={<TanpaResepShiftUmumPage />} />
        <Route path="/tanpa-resep-detail" element={<TanpaResepDetailPage />} />

        <Route path="/shift-riwayat" element={<RiwayatShiftPage />} />
        <Route path="/regular-riwayat" element={<RiwayatRegularPage />} />
        <Route path="/pres-riwayat" element={<RiwayatPresPage />} />

        <Route path="/koreksi" element={<KoreksiFormPage />} />

        <Route path="/riwayat-pbf" element={<RiwayatPBFPage />} />
        <Route path="/riwayat-non-pbf" element={<RiwayatNonPBFPage />} />

        <Route path="/atur-jenis" element={<AturJenisPage />} />
        <Route path="/atur-pengeluaran" element={<AturPengeluaranPage />} />

        <Route path="/draft" element={<OpnameDraftPage />} />
        <Route path="/draft-detail" element={<OpnameDetailPage />} />
      </Routes>
    </Router>
  );
}

export default App;
