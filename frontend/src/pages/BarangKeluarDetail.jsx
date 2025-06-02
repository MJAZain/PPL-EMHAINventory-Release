import React, { useEffect, useState } from "react";
import Button from "../components/buttonComp";
import Sidebar from "../components/Sidebar";
import { useNavigate } from "react-router-dom";
import DataTable from "../components/tableCompo";
import ActionMenu from "../components/ActionMenu";
import TambahBarangKeluarModal from "../components/modal/addBarangKeluarModal";
import { apiClient } from "../config/api";
import ConfirmDialog from "../components/ConfirmDialog";
import Toast from "../components/toast";

export default function BarangKeluarDetailPage() {
  const [toast, setToast] = useState(null);
  const [confirmOpen, setConfirmOpen] = useState(false);
const [itemToDelete, setItemToDelete] = useState(null);

  const [modalOpen, setModalOpen] = useState(false);
  const [formData, setFormData] = useState(null);
  const [editingItem, setEditingItem] = useState(null); // New state

const [barangList, setBarangList] = useState(() => {
  const saved = localStorage.getItem("barangList");
  return saved ? JSON.parse(saved) : []; // ← provide default value here
});


const navigate = useNavigate();

useEffect(() => {
  const storedForm = JSON.parse(localStorage.getItem("barangKeluarForm") || "null");
  if (storedForm) {
    setFormData(storedForm);
  }
}, []);

useEffect(() => {
  const storedBarangList = JSON.parse(localStorage.getItem("barangList") || "null");
  if (storedBarangList) {
    setBarangList(storedBarangList);
  }
}, []); 

useEffect(() => {
  localStorage.setItem("barangList", JSON.stringify(barangList));
}, [barangList]);


  const handleKembali = () => {
    localStorage.removeItem("barangKeluarForm");
    localStorage.removeItem("barangList");
    navigate("/barang-keluar");
  };

  
const handleEdit = (row) => {
  setEditingItem(row);
  setModalOpen(true);
};

const handleDeleteRequest = (row) => {
  setItemToDelete(row);
  setConfirmOpen(true);
};

const confirmDelete = () => {
  const updatedList = barangList.filter(
    (item) => item.product.id !== itemToDelete.product.id
  );
  setBarangList(updatedList);
  localStorage.setItem("barangList", JSON.stringify(updatedList));
  setConfirmOpen(false);
  setItemToDelete(null);
};

const cancelDelete = () => {
  setConfirmOpen(false);
  setItemToDelete(null);
};

const handleTambahBarang = (item) => {
  setBarangList((prev) => [...prev, item]);
};

const handleSave = async () => {
  try {
    const form = JSON.parse(localStorage.getItem("barangKeluarForm"));
    const barangList = JSON.parse(localStorage.getItem("barangList"));

    if (!form || !barangList || barangList.length === 0) {
      setToast({
        message: "Mohon lengkapi form dan tambahkan setidaknya satu produk.",
        type: "error",
      });
      return;
    }

    const payload = {
      outgoing_product: {
        date: new Date(form.date).toISOString(),
        customer: form.customer,
        no_faktur: form.no_faktur,
        payment_status: form.payment_status || "LUNAS",
      },
      details: barangList.map((item) => ({
        product_id: item.product.id,
        quantity: item.quantity,
        price: item.price,
        total: item.quantity * item.price,
      })),
    };


  const response = await apiClient.post("/outgoing-products/", payload);


    setToast({
        message: "Barang terjual berhasil disimpan!",
        type: "success",
      });

    localStorage.removeItem("barangKeluarForm");
    localStorage.removeItem("barangList");

    navigate("/barang-keluar");
  } catch (err) {
    console.error("Gagal menyimpan barang terjual:", err);
    setToast({
        message: "Gagal menyimpan barang terjual. Silakan coba lagi.",
        type: "error",
      });
    
  }
};


const columns = [
  { header: "Nama Obat", accessor: "product.name" },
  { header: "Kuantitas", accessor: "quantity" },
  { header: "Harga", accessor: "price",
    render: (row) => {
        const value = row.price;

        const formatted = new Intl.NumberFormat("id-ID", {
          style: "currency",
          currency: "IDR",
          minimumFractionDigits: 0,
        }).format(value);

        return formatted.replace("Rp", "Rp.");
      },
   },
  {
    header: "Aksi",
    accessor: "actions",
    isAction: true,
    render: (row) => (
      <ActionMenu
        actions={[
          { label: "Edit", onClick: () => handleEdit(row) },
          { label: "Hapus", onClick: () => handleDeleteRequest(row) },
        ]}
      />
    ),
  },
];

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="flex-1 p-8">
        <h1 className="text-2xl font-bold mb-6">Barang Terjual</h1>

        <div className="mb-4">
          <Button className="mb-4" 
          onClick={() => {
            setEditingItem(null); 
            setModalOpen(true)
          }} >
            Tambah Barang
          </Button>
        </div>

        <div className="border border-gray-300 p-6 rounded-md bg-white min-h-[150px] flex items-center justify-center text-center">
          {barangList.length === 0 ? (
            <p className="text-red-800 font-semibold font-[Open Sans] text-center">
                Mohon masukkan barang terlebih dahulu
            </p>
            ) : (
            <DataTable columns={columns} data={barangList} showIndex={true} />
            )}
        </div>

        <TambahBarangKeluarModal
          isOpen={modalOpen}
          onClose={() => {
            setModalOpen(false);
            setEditingItem(null);
          }}
          onSave={(updatedList) => {
            setBarangList(updatedList);
          }}
          initialData={editingItem}
        />


          {toast && (
            <Toast
              message={toast.message}
              type={toast.type}
              onClose={() => setToast(null)}
            />
          )}
        <ConfirmDialog
          isOpen={confirmOpen}
          title="Hapus Produk"
          description="Apakah Anda yakin ingin menghapus produk ini dari daftar?"
          onCancel={cancelDelete}
          onConfirm={confirmDelete}
        />

        <div className="flex justify-between mt-6 space-x-4">
          <Button className="w-full" onClick={handleKembali}>
            Kembali
          </Button>
          <Button className="w-full" onClick={handleSave}>
            Simpan
          </Button>
        </div>
      </div>
    </div>
  );
}
