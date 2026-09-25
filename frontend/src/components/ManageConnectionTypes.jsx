import { useEffect, useState } from 'react';
import { apiClient } from '../api';
import './ManageConnectionTypes.css';

export function ManageConnectionTypes({ onTypesChanged }) {
  const [types, setTypes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [editingId, setEditingId] = useState(null);
  const [editData, setEditData] = useState({ name: '', description: '', color: '#666' });
  const [formData, setFormData] = useState({ name: '', description: '', color: '#FF6B6B' });
  const [formError, setFormError] = useState('');
  const [formLoading, setFormLoading] = useState(false);

  useEffect(() => {
    loadTypes();
  }, []);

  const loadTypes = async () => {
    try {
      setLoading(true);
      const response = await apiClient.getConnectionTypes();
      setTypes(response.data || []);
      setError('');
    } catch (err) {
      setError('Failed to load connection types');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateSubmit = async (e) => {
    e.preventDefault();
    if (!formData.name.trim()) {
      setFormError('Name is required');
      return;
    }

    setFormLoading(true);
    setFormError('');
    try {
      await apiClient.createConnectionType(formData.name, formData.description, formData.color);
      setFormData({ name: '', description: '', color: '#FF6B6B' });
      loadTypes();
      onTypesChanged?.();
    } catch (err) {
      setFormError(err.response?.data?.error || 'Failed to create connection type');
    } finally {
      setFormLoading(false);
    }
  };

  const handleEdit = (type) => {
    setEditingId(type.id);
    setEditData({
      name: type.name,
      description: type.description,
      color: type.color || '#666',
    });
  };

  const handleSave = async (id) => {
    try {
      await apiClient.updateConnectionType(
        id,
        editData.name,
        editData.description,
        editData.color
      );
      setEditingId(null);
      setError('');
      loadTypes();
      onTypesChanged?.();
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to update connection type');
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this connection type?')) {
      return;
    }
    try {
      await apiClient.deleteConnectionType(id);
      setError('');
      loadTypes();
      onTypesChanged?.();
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to delete connection type');
    }
  };

  if (loading) {
    return <div className="manage-types-loading">Loading connection types...</div>;
  }

  return (
    <div className="manage-types-modal-content">
      {error && <p className="error">{error}</p>}
      
      <div className="create-form-section">
        <h3>Create New Type</h3>
        <form onSubmit={handleCreateSubmit} className="create-form">
          <input
            type="text"
            placeholder="Type name (e.g., 'knows', 'works-with')"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            disabled={formLoading}
            className="form-input"
          />
          <input
            type="text"
            placeholder="Description (optional)"
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            disabled={formLoading}
            className="form-input"
          />
          <div className="form-color-row">
            <label htmlFor="create-color">Color:</label>
            <input
              id="create-color"
              type="color"
              value={formData.color}
              onChange={(e) => setFormData({ ...formData, color: e.target.value })}
              disabled={formLoading}
              className="form-color-input"
            />
          </div>
          <button type="submit" disabled={formLoading} className="btn-create">
            {formLoading ? 'Creating...' : 'Create Type'}
          </button>
          {formError && <p className="error">{formError}</p>}
        </form>
      </div>

      <div className="divider" />

      <div className="list-section">
        <h3>Connection Types</h3>
        <div className="types-list">
          {types.length === 0 ? (
            <p className="empty-state">No connection types yet.</p>
          ) : (
            types.map((type) => (
              <div key={type.id} className="type-item">
                {editingId === type.id ? (
                  <div className="type-editor">
                    <div className="editor-row">
                      <input
                        type="text"
                        value={editData.name}
                        onChange={(e) => setEditData({ ...editData, name: e.target.value })}
                        placeholder="Name"
                        className="editor-input"
                      />
                      <input
                        type="text"
                        value={editData.description}
                        onChange={(e) => setEditData({ ...editData, description: e.target.value })}
                        placeholder="Description"
                        className="editor-input"
                      />
                    </div>
                    <div className="editor-row">
                      <div style={{ display: 'flex', gap: '10px', alignItems: 'center', flex: 1 }}>
                        <label htmlFor={`color-${type.id}`} style={{ marginBottom: 0, minWidth: '50px' }}>Color:</label>
                        <input
                          id={`color-${type.id}`}
                          type="color"
                          value={editData.color}
                          onChange={(e) => setEditData({ ...editData, color: e.target.value })}
                          className="editor-color"
                        />
                      </div>
                      <div className="editor-buttons">
                        <button onClick={() => handleSave(type.id)} className="btn-save">Save</button>
                        <button onClick={() => setEditingId(null)} className="btn-cancel">Cancel</button>
                      </div>
                    </div>
                  </div>
                ) : (
                  <div className="type-view">
                    <div className="type-info">
                      <h4>{type.name}</h4>
                      {type.description && <p>{type.description}</p>}
                    </div>
                    <div
                      className="type-color-preview"
                      style={{ backgroundColor: type.color || '#666' }}
                      title={type.color}
                    />
                    <div className="type-actions">
                      <button onClick={() => handleEdit(type)} className="btn-edit">
                        Edit
                      </button>
                      <button onClick={() => handleDelete(type.id)} className="btn-delete">
                        Delete
                      </button>
                    </div>
                  </div>
                )}
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}
