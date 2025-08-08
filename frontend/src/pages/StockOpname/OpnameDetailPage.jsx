import React, { useEffect, useState } from "react";
import Button from "../../components/buttonComp";
import { useNavigate, useParams } from "react-router-dom";
import DataTable from "../../components/tableCompo";
import ActionMenu from "../../components/ActionMenu";
import OpnameProductModal from "./OpnameProductModal";
import ConfirmDialog from "../../components/ConfirmDialog";
import Toast from "../../components/toast";
import { apiClient } from "../../config/api";

export default function OpnameDetailPage() {
  const { id } = useParams(); // opnameID
  const navigate = useNavigate();

  const [toast, setToast] = useState(null);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [itemToDelete, setItemToDelete] = useState(null);

  const [modalOpen, setModalOpen] = useState(false);
  const [editingItem, setEditingItem] = useState(null);

  const [opList, setOpList] = useState([]);
  const [isConfirmed, setIsConfirmed] = useState(false); // controls stage
  const [cancelConfirmOpen, setCancelConfirmOpen] = useState(false);

  useEffect(() => {
    const saved = JSON.parse(localStorage.getItem("opList") || "[]");
    setOpList(saved);

    // Simulated fetch — you may replace this with real status checking
    // If the opname is already confirmed, update the state accordingly
    const status = localStorage.getItem("opnameStatus"); // e.g., "draft" or "recording"
    setIsConfirmed(status === "recording");
  }, []);

  const handleEdit = (row) => {
    const initialData = {
      product: {
        id: row.product_id,
        name: row.product_name,
        quantity: row.quantity,
      },
    };
    setEditingItem(initialData);
    setModalOpen(true);
  };

  const handleDeleteRequest = (row) => {
    setItemToDelete(row);
    setConfirmOpen(true);
  };

  const confirmDelete = async () => {
    const opnameId = localStorage.getItem("opnameId");
    if (!opnameId || !itemToDelete?.detail_id) {
      setToast({ message: "Data tidak lengkap untuk menghapus item.", type: "error" });
      setConfirmOpen(false);
      return;
    }

    try {
      await apiClient.delete(`/stock-opname/draft/${opnameId}/products/${itemToDelete.detail_id}`);
      const updated = opList.filter((item) => item.product_id !== itemToDelete.product_id);
      localStorage.setItem("opList", JSON.stringify(updated));
      setOpList(updated);
      setToast({ message: "Produk berhasil dihapus.", type: "success" });
    } catch (err) {
      setToast({ message: "Gagal menghapus produk dari server.", type: "error" });
    } finally {
      setConfirmOpen(false);
      setItemToDelete(null);
    }
  };

  const cancelDelete = () => {
    setConfirmOpen(false);
    setItemToDelete(null);
  };

  const handleConfirmOpname = async () => {
    const opnameId = localStorage.getItem("opnameId");
    if (!opnameId) return setToast({ message: "opnameId tidak ditemukan.", type: "error" });

    try {
      await apiClient.post(`/stock-opname/${opnameId}/start`);
      setIsConfirmed(true);
      localStorage.setItem("opnameStatus", "recording");
      setToast({ message: "Stok opname dimulai. Silakan isi jumlah fisik.", type: "success" });
    } catch (err) {
      setToast({ message: "Gagal memulai stok opname.", type: "error" });
    }
  };

  const handleSaveActualStock = async (detailId, actualStock) => {
    try {
      await apiClient.put(`/stock-opname/details/${detailId}/record`, {
        actual_stock: actualStock,
      });

      const updated = opList.map((item) =>
        item.detail_id === detailId ? {
          ...item,
          actual_stock: actualStock,
          disparency: actualStock - item.quantity,
        } : item
      );

      localStorage.setItem("opList", JSON.stringify(updated));
      setOpList(updated);
      setToast({ message: "Jumlah fisik berhasil disimpan.", type: "success" });
    } catch (err) {
      setToast({ message: "Gagal menyimpan jumlah fisik.", type: "error" });
    }
  };

  const handleCancelOpname = async () => {
    setCancelConfirmOpen(true);
  };

  const confirmCancelOpname = async () => {
    const opnameId = localStorage.getItem("opnameId");
    try {
      await apiClient.post(`/stock-opname/${opnameId}/cancel`);
      localStorage.removeItem("opList");
      localStorage.removeItem("opnameStatus");
      navigate("/stock-opname");
    } catch (err) {
      setToast({ message: "Gagal membatalkan opname.", type: "error" });
    }
  };

  const handleConfirmStockRecord = async () => {
    const opnameId = localStorage.getItem("opnameId");
    try {
      await apiClient.post(`/stock-opname/${opnameId}/complete`);
      setToast({ message: "Stok opname selesai!", type: "success" });
      localStorage.removeItem("opList");
      localStorage.removeItem("opnameStatus");
      navigate("/draft");
    } catch (err) {
      setToast({ message: "Gagal menyelesaikan stok opname.", type: "error" });
    }
  };

  const columnsBefore = [
    { header: "Nama Obat", accessor: "product_name" },
    { header: "Jumlah", accessor: "quantity" },
    {
      header: "Aksi",
      accessor: "actions",
      isAction: true,
      render: (row) => (
        <ActionMenu
          actions={[
            { label: "Hapus", onClick: () => handleDeleteRequest(row) },
          ]}
        />
      ),
    },
  ];

  const columnsAfter = [
    { header: "Nama Obat", accessor: "product_name" },
    { header: "Jumlah dalam Sistem", accessor: "quantity" },
{
  header: "Jumlah Fisik",
  accessor: "actual_stock",
  render: (row) => (
    <input
      type="number"
      min={0}
      className="border rounded px-2 py-1 w-20"
      value={row.actual_stock ?? ""}
      onChange={(e) => {
        let value = parseInt(e.target.value, 10);
        if (isNaN(value) || value < 0) value = 0;

        const updated = opList.map((item) =>
          item.detail_id === row.detail_id
            ? { ...item, actual_stock: value }
            : item
        );
        setOpList(updated);
      }}
    />
  ),
},
    {
      header: "Perbedaan",
      accessor: "disparency",
      render: (row) => (
        <span>{(row.actual_stock ?? 0) - row.quantity}</span>
      ),
    },
    {
      header: "",
      accessor: "save_button",
      render: (row) => (
        <Button
          size="sm"
          onClick={() =>
            handleSaveActualStock(row.detail_id, row.actual_stock || 0)
          }
        >
          Simpan Jumlah Fisik
        </Button>
      ),
    },
  ];

  return (
    <div className="flex min-h-screen bg-gray-100">
      <div className="flex-1 p-8">
        <h1 className="text-2xl font-bold mb-6">Detail Stok Opname</h1>

        {!isConfirmed && (
          <div className="flex justify-between mb-4">
            <Button onClick={() => {
              setEditingItem(null);
              setModalOpen(true);
            }}>
              Tambah Barang
            </Button>
          </div>
        )}

        <div className="border border-gray-300 p-6 rounded-md bg-white min-h-[150px] w-full">
          {opList.length === 0 ? (
            <div className="flex items-center justify-center h-full">
              <p className="text-red-800 font-semibold text-center">
                Mohon masukkan barang terlebih dahulu
              </p>
            </div>
          ) : (
            <DataTable
              columns={isConfirmed ? columnsAfter : columnsBefore}
              data={opList}
              showIndex={true}
            />
          )}
        </div>

        <OpnameProductModal
          isOpen={modalOpen}
          onClose={() => {
            setModalOpen(false);
            setEditingItem(null);
          }}
          onSave={(updatedList) => {
            setOpList(updatedList);
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

        <ConfirmDialog
          isOpen={cancelConfirmOpen}
          title="Batalkan Stok Opname"
          description="Apakah Anda yakin ingin membatalkan stok opname ini?"
          onCancel={() => setCancelConfirmOpen(false)}
          onConfirm={confirmCancelOpname}
        />

        <div className="flex justify-between mt-6 space-x-4">
          {isConfirmed ? (
            <>
              <Button className="w-full bg-red-600 hover:bg-red-700" onClick={handleCancelOpname}>
                Batal
              </Button>
              <Button className="w-full" onClick={handleConfirmStockRecord}>
                Konfirmasi Stok Fisik
              </Button>
            </>
          ) : (
            <>
              <button
                className="w-full bg-gray-200 border border-black text-black rounded-md py-2 hover:bg-gray-300 transition"
                onClick={() => navigate(-1)}
              >
                Kembali
              </button>
              <Button className="w-full" onClick={handleConfirmOpname}>
                Konfirmasi Stok Opname
              </Button>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
