import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { apiClient } from '../../config/api';

const DraftDetailPage = () => {
  const { draftId } = useParams();
  const [draft, setDraft] = useState(null);
  const [products, setProducts] = useState([]);
  const [actualCounts, setActualCounts] = useState({});

  useEffect(() => {
    fetchDraft();
    fetchProducts();
  }, []);

  const fetchDraft = async () => {
    try {
      const response = await apiClient.get(`/stock-opname/draft/${draftId}`);
      setDraft(response.data);
    } catch (error) {
      console.error('Error fetching draft:', error);
    }
  };

  const fetchProducts = async () => {
    try {
      const response = await apiClient.get('/stock-opname/products');
      setProducts(response.data.products);
    } catch (error) {
      console.error('Error fetching products:', error);
    }
  };

  const addProductToDraft = async (productId) => {
    try {
      await apiClient.post(`/stock-opname/draft/${draftId}/products`, {
        product_id: productId,
      });
      fetchDraft();
    } catch (error) {
      console.error('Error adding product:', error);
    }
  };

  const removeProductFromDraft = async (detailId) => {
    try {
      await apiClient.delete(
        `/stock-opname/draft/${draftId}/products/${detailId}`
      );
      fetchDraft();
    } catch (error) {
      console.error('Error removing product:', error);
    }
  };

  const recordActualStock = async (detailId, actualStock) => {
    try {
      await apiClient.put(`/stock-opname/details/${detailId}/record`, {
        actual_stock: actualStock,
      });
      fetchDraft();
    } catch (error) {
      console.error('Error recording actual stock:', error);
    }
  };

  const startOpname = async () => {
    try {
      await apiClient.post(`/stock-opname/${draftId}/start`);
      fetchDraft();
    } catch (error) {
      console.error('Error starting opname:', error);
    }
  };

  const completeOpname = async () => {
    try {
      await apiClient.post(`/stock-opname/${draftId}/complete`);
      fetchDraft();
    } catch (error) {
      console.error('Error completing opname:', error);
    }
  };

  const cancelOpname = async () => {
    try {
      await apiClient.post(`/stock-opname/${draftId}/cancel`);
      fetchDraft();
    } catch (error) {
      console.error('Error canceling opname:', error);
    }
  };

  if (!draft) return <div>Loading...</div>;

  return (
    <div>
      <h1>{draft.name}</h1>
      <button onClick={startOpname}>Start Opname</button>
      <button onClick={completeOpname}>Complete Opname</button>
      <button onClick={cancelOpname}>Cancel Opname</button>
      <h2>Products in Draft</h2>
      <ul>
        {draft.details.map((detail) => (
          <li key={detail.id}>
            {detail.product.name} - Expected: {detail.expected_stock}
            <input
              type="number"
              value={actualCounts[detail.id] || ''}
              onChange={(e) =>
                setActualCounts({
                  ...actualCounts,
                  [detail.id]: e.target.value,
                })
              }
            />
            <button
              onClick={() =>
                recordActualStock(detail.id, actualCounts[detail.id])
              }
            >
              Record Actual Stock
            </button>
            <button onClick={() => removeProductFromDraft(detail.id)}>
              Remove
            </button>
          </li>
        ))}
      </ul>
      <h2>Add Product to Draft</h2>
      <ul>
        {products.map((product) => (
          <li key={product.id}>
            {product.name}
            <button onClick={() => addProductToDraft(product.id)}>Add</button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default DraftDetailPage;
