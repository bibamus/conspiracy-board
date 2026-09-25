import axios from 'axios';

// Determine API base URL
const getApiBase = () => {
  // If VITE_API_URL is explicitly set, use it
  if (import.meta.env.VITE_API_URL) {
    return import.meta.env.VITE_API_URL;
  }

  // Local development: use localhost:8080
  if (window.location.hostname === 'localhost' && window.location.port !== '8000') {
    return 'http://localhost:8000/api';
  }

  // Docker/Production: use same origin
  return `${window.location.origin}/api`;
};

const API_BASE = getApiBase();

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
