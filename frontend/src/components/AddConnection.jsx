import { useState, useEffect } from 'react';
import { apiClient } from '../api';
import './AddConnection.css';

export function AddConnection({ onConnectionAdded, refreshTrigger }) {
  const [people, setPeople] = useState([]);
  const [types, setTypes] = useState([]);
  const [fromPersonId, setFromPersonId] = useState('');
  const [toPersonId, setToPersonId] = useState('');
  const [typeId, setTypeId] = useState('');
  const [description, setDescription] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    loadData();
  }, [refreshTrigger]);

  const loadData = async () => {
    try {
      const [peopleRes, typesRes] = await Promise.all([
        apiClient.getPeople(),
        apiClient.getConnectionTypes(),
      ]);
      setPeople(peopleRes.data || []);
      setTypes(typesRes.data || []);
    } catch (err) {
      setError('Failed to load data');
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!fromPersonId || !toPersonId || !typeId) {
      setError('Please fill in all required fields');
      return;
    }

    setLoading(true);
    setError('');
    try {
      const response = await apiClient.createConnection(
        parseInt(fromPersonId),
        parseInt(toPersonId),
        parseInt(typeId),
        description
      );
      setFromPersonId('');
      setToPersonId('');
      setTypeId('');
      setDescription('');
      onConnectionAdded(response.data);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to add connection');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="form-card">
      <h3>Add Connection</h3>
      <form onSubmit={handleSubmit}>
        <select
          value={fromPersonId}
          onChange={(e) => setFromPersonId(e.target.value)}
          disabled={loading}
        >
          <option value="">From Person</option>
          {people.map((person) => (
            <option key={person.id} value={person.id}>
              {person.name}
            </option>
          ))}
        </select>

        <select
          value={toPersonId}
          onChange={(e) => setToPersonId(e.target.value)}
          disabled={loading}
        >
          <option value="">To Person</option>
          {people.map((person) => (
            <option key={person.id} value={person.id}>
              {person.name}
            </option>
          ))}
        </select>

        <select
          value={typeId}
          onChange={(e) => setTypeId(e.target.value)}
          disabled={loading}
        >
          <option value="">Connection Type</option>
          {types.map((type) => (
            <option key={type.id} value={type.id}>
              {type.name}
            </option>
          ))}
        </select>

        <input
          type="text"
          placeholder="Description (optional)"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          disabled={loading}
        />

        <button type="submit" disabled={loading}>
          {loading ? 'Adding...' : 'Add Connection'}
        </button>
        {error && <p className="error">{error}</p>}
      </form>
    </div>
  );
}
