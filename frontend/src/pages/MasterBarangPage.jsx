import React, { useState, useEffect } from "react";
import useSearch from "../hooks/useSearch";
import useProductActions from "../hooks/useProductsAction";
import { useNavigate } from 'react-router-dom';

import { apiClient } from "../config/api";
import { PlusIcon } from '@heroicons/react/24/solid';

import ActionMenu from "../components/ActionMenu";
import SearchBar from "../components/SearchBar";
import DataTable from "../components/tableCompo";
import Button from "../components/buttonComp";
import AddBarangModal from "../components/modal/addBarangModal";
import EditBarangModal from "../components/modal/editBarangModal";
import Sidebar from "../components/Sidebar";
import ConfirmDialog from "../components/ConfirmDialog";
import Toast from "../components/toast";


function MasterBarangPage() {
  const [isEditOpen, setEditOpen] = useState(false);
  const [editId, setEditId] = useState(null);
  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [deleteTargetId, setDeleteTargetId] = useState(null);
  const [barangList, setBarangList] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [toast, setToast] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const navigate = useNavigate();

  const { getProductById, deleteProduct } = useProductActions();

 const fetchBarang = async () => {
  try {
    const response = await apiClient.get("/products/");
    return response.data;
  } catch (err) {
    setError(err.message || "Gagal mengambil data.");
    console.error("Fetch error:", err);
    return []; // Fallback to empty array
  }
};


  useEffect(() => {
  const loadBarang = async () => {
    try {
      const data = await fetchBarang();

      if (Array.isArray(data)) {
        setBarangList(data);
      } else if (Array.isArray(data?.data)) {
        setBarangList(data.data);
      } else {
        console.error("Unexpected data format:", data);
        setBarangList([]);
      }
    } catch (error) {
      console.error("Error loading barang:", error);
      setBarangList([]);
    } finally {
      setLoading(false);
    }
  };

  loadBarang();
}, []);

  const handleAddBarang = async () => {
  setToast({ message: "Barang berhasil ditambahkan", type: "success" });

  const data = await fetchBarang();
    if (Array.isArray(data)) {
      setBarangList(data);
    } else if (Array.isArray(data?.data)) {
      setBarangList(data.data);
    }
  setIsModalOpen(false);
  };

  const handleEditSuccess = async () => {
  setToast({
    message: "Barang berhasil diperbarui",
    type: "success",
  });

  const data = await fetchBarang();
    if (Array.isArray(data)) {
      setBarangList(data);
    } else if (Array.isArray(data?.data)) {
      setBarangList(data.data);
    }
  setEditOpen(false);
};


  const handleEdit = async (id) => {
  try {
      const product = await getProductById(id);
      setEditId(id);
      setEditOpen(true);
    } catch (err) {
      alert("Gagal mengambil data produk.");
    }
  };

  const handleDeleteRequest = (id) => {
  setDeleteTargetId(id);
  setIsConfirmOpen(true);
  };

  const handleCancelDelete = () => {
    setIsConfirmOpen(false);
    setDeleteTargetId(null);
  };



const handleConfirmDelete = async () => {
  try {
    await deleteProduct(deleteTargetId);
    const data = await fetchBarang();
    if (Array.isArray(data)) {
      setBarangList(data);
    } else if (Array.isArray(data?.data)) {
      setBarangList(data.data);
    }
    setToast({
      message: "Barang berhasil dihapus",
      type: "success",
    });
  } catch (err) {
    alert("Gagal menghapus data.");
  } finally {
    setIsConfirmOpen(false);
    setDeleteTargetId(null);
  }
};

  const { searchTerm, setSearchTerm, filteredData } = useSearch(barangList, [
    "nama",
    "sku",
    "kategori",
    "merk",
  ]);

  const columns = [
    { header: "Nama", accessor: "name" },
    { header: "SKU", accessor: "code" },
    { header: "Barcode", accessor: "barcode" },
    { header: "Kategori", accessor: "category_id" },
    { header: "Satuan", accessor: "unit_id" },
    { header: "Isi", accessor: "package_content" },
    { header: "Harga Beli", accessor: "purchase_price" },
    { header: "Harga Jual", accessor: "selling_price" },
    { header: "Harga Grosir", accessor: "wholesale_price" },
    { header: "Stok Buffer", accessor: "stock_buffer" },
    { header: "Lokasi", accessor: "storage_location" },
    { header: "Merk", accessor: "brand" },
    {
      header: "Pilih Aksi",
      accessor: "actions",
      isAction: true,
      render: (item) => (
        <ActionMenu
          actions={[
            { label: "Edit", onClick: () => handleEdit(item.id) },
            { label: "Delete", onClick: () => handleDeleteRequest(item.id) },
          ]}
        />
      ),
    },
  ];


  if (loading) return <div className="text-center mt-6">Loading...</div>;

  return (
    <div className="flex">
      
      <div className="bg-white min-h-screen">
        <Sidebar />
      </div>
      
      <div className="p-5">
        <h1 className="text-2xl font-bold mb-6">Daftar Barang</h1>
        <SearchBar value={searchTerm} onChange={setSearchTerm} />

        <div className=" border-1 rounded-md border-gray-300 bg-white p-5">
          <div className="flex gap-4 mb-4">
             <button onClick={() => setIsModalOpen(true)}
              className="flex items-center text-blue-700 font-semibold space-x-1 bg-transparent border border-blue-700 py-2 px-4 rounded-md"
            >
              <PlusIcon className="w-4 h-4" />
                <span>Tambah Barang</span>
            </button>
            <Button onClick={() => navigate('/stock-opname')}>
              Stock Opname
            </Button>
          </div>
        
          <div className="max-w-[1121px]">
            <DataTable columns={columns} data={filteredData} showIndex={true} />
          </div>
        
        </div>
     </div>
      

      <AddBarangModal
        isOpen={isModalOpen}
        close={() => setIsModalOpen(false)}
        onSubmit={handleAddBarang}
      />

      <EditBarangModal
        isOpen={isEditOpen}
        close={() => setEditOpen(false)}
        productId={editId}
        onSubmit={handleEditSuccess}
      />

      <ConfirmDialog
        isOpen={isConfirmOpen}
        title="Konfirmasi Penghapusan"
        description="Apakah Anda yakin ingin menghapus obat ini?"
        onCancel={handleCancelDelete}
        onConfirm={handleConfirmDelete}
      />

      {toast && (
        <Toast
          message={toast.message}
          type={toast.type}
          onClose={() => setToast(null)}
        />
      )}
    </div>
  );
}

export default MasterBarangPage;
