import { useState } from 'react';
import { apiClient } from '../api';
import './AddPerson.css';

export function AddPerson({ onPersonAdded }) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!name.trim()) {
      setError('Name is required');
      return;
    }

    setLoading(true);
    setError('');
    try {
      const response = await apiClient.createPerson(name, description);
      setName('');
      setDescription('');
      onPersonAdded(response.data);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to add person');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="form-card">
      <h3>Add Person</h3>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          disabled={loading}
        />
        <input
          type="text"
          placeholder="Description (optional)"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          disabled={loading}
        />
        <button type="submit" disabled={loading}>
          {loading ? 'Adding...' : 'Add Person'}
        </button>
        {error && <p className="error">{error}</p>}
      </form>
    </div>
  );
}
