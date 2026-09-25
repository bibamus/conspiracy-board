import axios from 'axios';

// Use the same domain/origin as the application is running on
// Falls back to localhost:8080/api if VITE_API_URL is set
const API_BASE = import.meta.env.VITE_API_URL || `${window.location.origin}/api`;

export const api = axios.create({
  baseURL: API_BASE,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const apiClient = {
  // Connection Types
  createConnectionType: (name, description, color = '#666') =>
    api.post('/connection-types', { name, description, color }),
  getConnectionTypes: () => api.get('/connection-types'),
  updateConnectionType: (id, name, description, color) =>
    api.put(`/connection-types/${id}`, { name, description, color }),
  deleteConnectionType: (id) =>
    api.delete(`/connection-types/${id}`),

  // People
  createPerson: (name, description) =>
    api.post('/people', { name, description }),
  getPeople: () => api.get('/people'),
  getPerson: (id) => api.get(`/people/${id}`),
  deletePerson: (id) => api.delete(`/people/${id}`),

  // Connections
  createConnection: (fromPersonId, toPersonId, typeId, description) =>
    api.post('/connections', {
      from_person_id: fromPersonId,
      to_person_id: toPersonId,
      type_id: typeId,
      description,
    }),
  getConnections: () => api.get('/connections'),
  getPersonConnections: (id) => api.get(`/people/${id}/connections`),
  deleteConnection: (id) => api.delete(`/connections/${id}`),

  // Graph
  getGraph: () => api.get('/graph'),
};
