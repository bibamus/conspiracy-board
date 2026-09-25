import { useState } from 'react';
import { apiClient } from '../api';
import './AddConnectionType.css';

export function AddConnectionType({ onTypeAdded }) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [color, setColor] = useState('#FF6B6B');
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
      const response = await apiClient.createConnectionType(name, description, color);
      setName('');
      setDescription('');
      setColor('#FF6B6B');
      onTypeAdded(response.data);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to add connection type');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="form-card">
      <h3>Add Connection Type</h3>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="Type name (e.g., 'knows', 'works-with')"
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
        <div style={{ display: 'flex', gap: '10px', alignItems: 'center' }}>
          <label htmlFor="color" style={{ marginBottom: 0 }}>Color:</label>
          <input
            id="color"
            type="color"
            value={color}
            onChange={(e) => setColor(e.target.value)}
            disabled={loading}
            style={{ width: '50px', cursor: 'pointer' }}
          />
        </div>
        <button type="submit" disabled={loading}>
          {loading ? 'Adding...' : 'Add Type'}
        </button>
        {error && <p className="error">{error}</p>}
      </form>
    </div>
  );
}
