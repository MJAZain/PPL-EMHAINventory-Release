import React from 'react';
import { Link } from 'react-router-dom';

const DraftList = ({ drafts }) => {
  if (!drafts || drafts.length === 0) {
    return <p className="text-red-800 font-semibold font-[Open Sans] text-center">
      No drafts available.
      </p>;
  }

  return (
    <div>
      <ul>
        {drafts.map((draft) => (
          <li key={draft.id}>
            <Link to={`/stock-opname/draft/${draft.id}`}>{draft.name}</Link>
          </li>
        ))}
      </ul>
    </div>
  );
};


export default DraftList;
