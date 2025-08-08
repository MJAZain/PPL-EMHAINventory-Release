import React, { useEffect, useState } from "react";
import Button from "../../../components/buttonComp";
import { useNavigate, useParams } from "react-router-dom";
import DataTable from "../../../components/tableCompo";
import ActionMenu from "../../../components/ActionMenu";
import ResepProductModal from "./ResepProductModal";
import AkhiriTransaksiModal from "./AkhiriTransaksiModal";
import AkhiriShiftModal from "./AkhiriShiftModal";
import ConfirmDialog from "../../../components/ConfirmDialog";
import Toast from "../../../components/toast";

export default function ResepDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();

  const [toast, setToast] = useState(null);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [itemToDelete, setItemToDelete] = useState(null);

  const [modalOpen, setModalOpen] = useState(false);
  const [editingItem, setEditingItem] = useState(null);

  const [presList, setPresList] = useState([]);
  const [transaksiModalOpen, setTransaksiModalOpen] = useState(false);
  const [shiftModalOpen, setShiftModalOpen] = useState(false);

  useEffect(() => {
    const saved = JSON.parse(localStorage.getItem("presList") || "[]");
    setPresList(saved);
  }, []);

  const handleEdit = (row) => {
    const initialData = {
      product: {
        id: row.product_id,
        name: row.name,
        selling_price: row.price,
        code: row.code,
        unit: row.unit,
      },
      quantity: row.quantity,
    };

    setEditingItem(initialData);
    setModalOpen(true);
  };

  const handleDeleteRequest = (row) => {
    setItemToDelete({ product_id: row.product_id });
    setConfirmOpen(true);
  };

  const confirmDelete = () => {
    const updatedList = presList.filter(
      (item) => item.product_id !== itemToDelete.product_id
    );
    localStorage.setItem("presList", JSON.stringify(updatedList));
    setPresList(updatedList);
    setConfirmOpen(false);
    setItemToDelete(null);
  };

  const cancelDelete = () => {
    setConfirmOpen(false);
    setItemToDelete(null);
  };

  const calculateTotalPembelian = () => {
    return presList.reduce(
      (sum, item) => sum + item.quantity * item.price,
      0
    );
  };

  const columns = [
    { header: "Nama Obat", accessor: "name" },
    { header: "Kode Barang", accessor: "code" },
    {
      header: "Harga Beli Satuan",
      accessor: "price",
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
    { header: "Kuantitas", accessor: "quantity" },
    {
      header: "Total Beli Produk",
      accessor: "computed_total_price",
      render: (row) => {
        const total = row.quantity * row.price;
        return new Intl.NumberFormat("id-ID", {
          style: "currency",
          currency: "IDR",
          minimumFractionDigits: 0,
        }).format(total).replace("Rp", "Rp.");
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
    <div className="flex min-h-screen bg-gray-100">
      <div className="flex-1 p-8">
        <h1 className="text-2xl font-bold mb-6">Kasir dengan Resep</h1>

        <div className="flex items-center justify-between mb-4">
          <Button
            className="mb-4"
            onClick={() => {
              setEditingItem(null);
              setModalOpen(true);
            }}
          >
            Tambah Barang
          </Button>
        </div>

        <div className="border border-gray-300 p-6 rounded-md bg-white min-h-[150px] w-full">
          {presList.length === 0 ? (
            <div className="flex items-center justify-center h-full">
              <p className="text-red-800 font-semibold font-[Open Sans] text-center">
                Mohon masukkan barang terlebih dahulu
              </p>
            </div>
          ) : (
            <>
              <DataTable columns={columns} data={presList} showIndex={true} />

              <div className="flex justify-end mt-4">
                <div className="text-right font-semibold text-lg">
                  Total Pembelian:{" "}
                  {new Intl.NumberFormat("id-ID", {
                    style: "currency",
                    currency: "IDR",
                    minimumFractionDigits: 0,
                  })
                    .format(calculateTotalPembelian())
                    .replace("Rp", "Rp.")}
                </div>
              </div>
            </>
          )}
        </div>

        <ResepProductModal
          isOpen={modalOpen}
          onClose={() => {
            setModalOpen(false);
            setEditingItem(null);
          }}
          onSave={(updatedList) => {
            setPresList(updatedList);
          }}
          initialData={editingItem}
        />

        <AkhiriTransaksiModal
          isOpen={transaksiModalOpen}
          onClose={() => setTransaksiModalOpen(false)}
          presList={presList}
          onAfterSubmit={() => {
            setToast({ message: "Transaksi Berhasil", type: "success" });
            setPresList([]);
          }}
        />

        <AkhiriShiftModal
          isOpen={shiftModalOpen}
          onClose={() => setShiftModalOpen(false)}
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
          <button
            className="w-full bg-gray-200 border border-black text-black rounded-md py-2 hover:bg-gray-300 transition"
            onClick={() => setShiftModalOpen(true)}
          >
            Akhiri Shift
          </button>
          <Button className="w-full" onClick={() => setTransaksiModalOpen(true)}>
            Akhiri Transaksi
          </Button>
        </div>
      </div>
    </div>
  );
}
