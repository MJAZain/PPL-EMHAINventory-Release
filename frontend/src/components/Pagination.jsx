import React from "react";

const Pagination = ({
  currentPage,
  totalPages,
  onPageChange,
  itemsPerPage,
  totalData,
}) => {
  if (totalData === 0) {
    return null;
  }

  const pageNumbers = [];
  for (let i = 1; i <= totalPages; i++) {
    pageNumbers.push(i);
  }

  const firstItemIndex =
    totalData > 0 ? (currentPage - 1) * itemsPerPage + 1 : 0;
  const lastItemIndex =
    totalData > 0 ? Math.min(currentPage * itemsPerPage, totalData) : 0;

  return (
    <div className="flex flex-col sm:flex-row justify-between items-center mt-6 py-3">
      <div className="text-sm text-gray-700 mb-2 sm:mb-0">
        {totalData > 0 ? (
          <>
            Menampilkan <span className="font-medium">{firstItemIndex}</span> -{" "}
            <span className="font-medium">{lastItemIndex}</span> dari{" "}
            <span className="font-medium">{totalData}</span> hasil
          </>
        ) : (
          "Tidak ada data"
        )}
      </div>
      {totalPages > 0 && (
        <nav aria-label="Pagination">
          <ul className="inline-flex items-center -space-x-px">
            <li>
              <button
                onClick={() => onPageChange(currentPage - 1)}
                disabled={currentPage === 1}
                className={`px-3 py-2 ml-0 leading-tight text-gray-500 bg-white border border-gray-300 rounded-l-lg hover:bg-gray-100 hover:text-gray-700 disabled:opacity-50 disabled:cursor-not-allowed`}
              >
                Sebelumnya
              </button>
            </li>
            {pageNumbers.map((number) =>
              number === currentPage ||
              Math.abs(number - currentPage) < 2 ||
              number === 1 ||
              number === totalPages ? (
                <li key={number}>
                  <button
                    onClick={() => onPageChange(number)}
                    className={`px-3 py-2 leading-tight border ${
                      currentPage === number
                        ? "text-blue-600 bg-blue-50 border-blue-300 hover:bg-blue-100 hover:text-blue-700"
                        : "text-gray-500 bg-white border-gray-300 hover:bg-gray-100 hover:text-gray-700"
                    }`}
                  >
                    {number}
                  </button>
                </li>
              ) : Math.abs(number - currentPage) === 2 &&
                number !== 1 &&
                number !== totalPages ? (
                <li key={`ellipsis-${number}`}>
                  <span className="px-3 py-2 leading-tight text-gray-500 bg-white border border-gray-300">
                    ...
                  </span>
                </li>
              ) : null
            )}
            <li>
              <button
                onClick={() => onPageChange(currentPage + 1)}
                disabled={currentPage === totalPages || totalPages === 0}
                className={`px-3 py-2 leading-tight text-gray-500 bg-white border border-gray-300 rounded-r-lg hover:bg-gray-100 hover:text-gray-700 disabled:opacity-50 disabled:cursor-not-allowed`}
              >
                Berikutnya
              </button>
            </li>
          </ul>
        </nav>
      )}
    </div>
  );
};

export default Pagination;
