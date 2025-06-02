import React from "react";
import ActionMenu from "./ActionMenu";

const getNestedValue = (obj, accessor) => {
  if (typeof accessor !== "string") return undefined;

  return accessor
    .split(".")
    .reduce((acc, key) => (acc && acc[key] !== undefined ? acc[key] : undefined), obj);
};


function DataTable({ columns, data, showIndex = false }) {
  return (
    <div className="overflow-x-auto pb-5">
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

            return (
              <tr key={row.id || rowIndex} className="hover:bg-gray-50">
                {showIndex && (
                  <td className="px-4 py-3 border-b border-gray-600 text-center text-[18px] leading-[24px] text-black font-[Open_Sans] font-normal">
                    {rowIndex + 1}
                  </td>
                )}
                {columns.map((col, colIndex) => (
                  <td
                    key={`${rowIndex}-${col.accessor || colIndex}`}
                    className="px-4 py-3 border-b border-gray-600 text-center text-[18px] leading-[24px] text-black font-[Open_Sans] font-normal"
                  >
                    {col.render ? (
                      col.render(row)
                    ) : col.isAction && col.getActions ? (
                      <ActionMenu actions={col.getActions(row)} />
                    ) : (
                      getNestedValue(row, col.accessor) ?? "-"
                    )}
                  </td>
                ))}
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}

export default DataTable;
