import './App.css';
import { useState, useCallback } from 'react';
import { AddPerson } from './components/AddPerson';
import { AddConnection } from './components/AddConnection';
import { ManageConnectionTypes } from './components/ManageConnectionTypes';
import { Modal } from './components/Modal';
import { GraphVisualization } from './components/GraphVisualization';

function App() {
  const [refreshKey, setRefreshKey] = useState(0);
  const [connectionRefresh, setConnectionRefresh] = useState(0);
  const [isManageModalOpen, setIsManageModalOpen] = useState(false);

  const handleDataChanged = useCallback(() => {
    setRefreshKey((prev) => prev + 1);
    setConnectionRefresh((prev) => prev + 1);
  }, []);

  const handlePersonAdded = useCallback(() => {
    setConnectionRefresh((prev) => prev + 1);
    handleDataChanged();
  }, [handleDataChanged]);

  return (
    <div className="app">
      <div className="container">
        <aside className="sidebar">
          <div className="sidebar-content">
            <button 
              className="manage-btn"
              onClick={() => setIsManageModalOpen(true)}
            >
              ➕ Manage Connection Types
            </button>
            <AddPerson onPersonAdded={handlePersonAdded} />
            <AddConnection onConnectionAdded={handleDataChanged} refreshTrigger={connectionRefresh} />
          </div>
        </aside>

        <main className="main">
          <GraphVisualization key={refreshKey} onRefresh={handleDataChanged} />
        </main>
      </div>

      <Modal 
        isOpen={isManageModalOpen} 
        onClose={() => setIsManageModalOpen(false)}
        title="Connection Types"
      >
        <ManageConnectionTypes />
      </Modal>
    </div>
  );
}

export default App;
