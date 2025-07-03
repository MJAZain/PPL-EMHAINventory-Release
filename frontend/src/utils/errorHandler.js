export function getFriendlyErrorMessage(error) {
  if (error.response) {
    const status = error.response.status;
    const dataMsg = error.response.data?.message;

    switch (status) {
      case 400:
        return "Data duplikat. Silakan periksa kembali data Anda.";
      case 401:
        return "Anda tidak memiliki akses. Silakan login kembali.";
      case 403:
        return "Anda tidak memiliki izin untuk melakukan aksi ini.";
      case 404:
        return "Data tidak ditemukan.";
      case 500:
        return "Terjadi kesalahan pada server. Silakan coba lagi nanti.";
      default:
        return dataMsg || "Terjadi kesalahan yang tidak diketahui.";
    }
  }

  return error.message || "Terjadi kesalahan jaringan. Silakan coba lagi.";
}
