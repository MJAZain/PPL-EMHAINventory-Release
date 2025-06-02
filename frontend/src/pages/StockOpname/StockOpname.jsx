import React, { useEffect, useState } from 'react';
import { apiClient } from '../../config/api';
import DraftList from './DraftList';
import CompletedList from './CompletedList';
import CreateDraftModal from './CreateDraftModal';
import Sidebar from '../../components/Sidebar';
import { PlusIcon } from '@heroicons/react/24/solid';

const StockOpnamePage = () => {
  const [drafts, setDrafts] = useState([]);
  const [completed, setCompleted] = useState([]);
  const [isModalOpen, setIsModalOpen] = useState(false);

  useEffect(() => {
    fetchDrafts();
    fetchCompleted();
  }, []);

  const fetchDrafts = async () => {
    try {
      const response = await apiClient.get('/stock-opname');
      console.log(response.data);
      setDrafts(response.data.drafts);
    } catch (error) {
      console.error('Error fetching drafts:', error);
    }
  };

  const fetchCompleted = async () => {
    try {
      const response = await apiClient.get('/stock-opname/history');
      setCompleted(response.data.completed);
    } catch (error) {
      console.error('Error fetching completed opnames:', error);
    }
  };

  return (
    <div className="flex min-h-screen bg-gray-100">
      <div className="bg-white min-h-screen">
        <Sidebar />
      </div>
      <div className="p-5 w-full py-10">
        <h1 className="text-2xl font-bold mb-6">Stock Opname</h1>

        <div className='border border-gray-400 rounded-xl bg-white p-5'>

          <div className="flex justify-start py-5">
            <button
              onClick={() => setIsModalOpen(true)}
              className="flex items-center text-blue-700 font-semibold space-x-1 bg-transparent border border-blue-700 py-2 px-4 rounded-md"
            >
              <PlusIcon className="w-4 h-4" />
              <span>Create New Draft</span>
            </button>
          </div>

          <div>
            <h2 className='font-bold mb-5'>Drafts</h2>
            <div className='border border-gray-400 rounded-xl bg-white p-5 mb-5'>
              <DraftList drafts={drafts} />
            </div>
          </div>
          <div>
            <h2 className='font-bold mb-5'>Completed Drafts</h2>
            <div className='border border-gray-400 rounded-xl bg-white p-5'>
              <CompletedList completed={completed} />
            </div>
          </div>
          
          <CreateDraftModal
            isOpen={isModalOpen}
            onClose={() => setIsModalOpen(false)}
            onDraftCreated={fetchDrafts}
          />

        </div>
        
      </div>
      
    </div>
  );
};

export default StockOpnamePage;
