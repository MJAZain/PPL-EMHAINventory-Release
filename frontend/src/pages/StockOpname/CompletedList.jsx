import React from 'react';

const CompletedList = ({ completed = [] }) => {
  if (completed.length === 0) {
    return <p className="text-red-800 font-semibold font-[Open Sans] text-center">No Completed Drafts</p>;
  }

  return (
    <div>
      <ul>
        {completed.map((opname) => (
          <li key={opname.id}>{opname.name}</li>
        ))}
      </ul>
    </div>
  );
};

export default CompletedList;
