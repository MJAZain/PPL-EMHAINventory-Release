import React from "react";
import ActionMenu from "./ActionMenu";

const getNestedValue = (obj, accessor) => {
  if (typeof accessor !== "string") return undefined;

  const value = accessor
    .split(".")
    .reduce(
      (acc, key) =>
        acc && typeof acc === "object" && acc[key] !== undefined
          ? acc[key]
          : undefined,
      obj
    );
  return value;
};

const formatDate = (dateStr) => {
  if (!dateStr) return "-";
  const date = new Date(dateStr);
  const day = String(date.getDate()).padStart(2, "0");
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const year = date.getFullYear();
  return `${day}-${month}-${year}`;
};


function DataTable({
  columns,
  data,
  showIndex = false,
  currentPage = 0,
  pageSize = 0,
}) {
  if (!data) {
    console.warn("DataTable: data prop is undefined or null.");
    return null;
  }
  if (!columns) {
    console.warn("DataTable: columns prop is undefined or null.");
    return null;
  }

  return (
    <div className="pb-5 overflow-x-auto">
      <table className="min-w-full text-sm border-collapse">
        <thead>
          <tr>
            {showIndex && (
              <th className="px-4 py-3 border-b-2 border-gray-600 text-center font-bold text-[18px] leading-[24px] text-black font-[Open_Sans]">
                No.
              </th>
            )}
            {columns.map((col) => (
              <th
                key={col.accessor || col.header}
                className="px-4 py-3 border-b-2 border-gray-600 text-center font-bold text-[18px] leading-[24px] text-black font-[Open_Sans]"
              >
                {col.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.map((row, rowIndex) => {
            if (!row) return null;
            const itemNumber =
              showIndex && pageSize > 0
                ? rowIndex + 1 + currentPage * pageSize
                : rowIndex + 1;

            return (
              <tr key={row.id || rowIndex} className="hover:bg-gray-100">
                {showIndex && (
                  <td className="px-4 py-3 border-b border-gray-400 text-center text-[18px] leading-[24px] text-black font-[Open_Sans] font-normal">
                    {itemNumber}
                  </td>
                )}
                {columns.map((col, colIndex) => {
                  let cellValue;
                  if (typeof col.accessor === "function" && !col.render) {
                    cellValue = col.accessor(row, rowIndex);
                  } else if (typeof col.accessor === "string" && !col.render) {
                    cellValue = getNestedValue(row, col.accessor);
                  }

                  return (
                    <td
                      key={`${row.id || rowIndex}-${col.accessor || col.header || colIndex}`}
                      className="px-4 py-3 border-b border-gray-400 text-center text-[18px] leading-[24px] text-black font-[Open_Sans] font-normal"
                    >
                      {col.render ? (
                        col.render(row, rowIndex)
                      ) : col.isAction && col.getActions ? (
                        <ActionMenu actions={col.getActions(row)} />
                      ) : col.isDate ? (
                        formatDate(cellValue)
                      ) : (
                        cellValue ?? "-"
                      )}
                    </td>
                  );
                })}
              </tr>
            );
          })}
          {data.length === 0 && (
            <tr>
              <td
                colSpan={showIndex ? columns.length + 1 : columns.length}
                className="px-4 py-5 border-b border-gray-400 text-center text-gray-500"
              >
                Tidak ada data untuk ditampilkan.
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}

export default DataTable;
