import './App.css';
import { useState, useCallback } from 'react';
import { AddPerson } from './components/AddPerson';
import { AddConnection } from './components/AddConnection';
import { ManageConnectionTypes } from './components/ManageConnectionTypes';
import { Modal } from './components/Modal';
import { GraphVisualization } from './components/GraphVisualization';

function App() {
  // Bumped after any mutation so every view reloads from the server.
  const [dataVersion, setDataVersion] = useState(0);
  const [isManageModalOpen, setIsManageModalOpen] = useState(false);

  const handleDataChanged = useCallback(() => {
    setDataVersion((prev) => prev + 1);
  }, []);

  return (
    <div className="app">
      <div className="container" inert={isManageModalOpen}>
        <aside className="sidebar">
          <div className="sidebar-content">
            <button 
              type="button"
              className="manage-btn"
              aria-haspopup="dialog"
              onClick={() => setIsManageModalOpen(true)}
            >
              <span aria-hidden="true">➕</span> Manage Connection Types
            </button>
            <AddPerson onPersonAdded={handleDataChanged} />
            <AddConnection onConnectionAdded={handleDataChanged} refreshTrigger={dataVersion} />
          </div>
        </aside>

        <main className="main">
          <GraphVisualization refreshTrigger={dataVersion} onDataChanged={handleDataChanged} />
        </main>
      </div>

      <Modal 
        isOpen={isManageModalOpen} 
        onClose={() => setIsManageModalOpen(false)}
        title="Connection Types"
      >
        <ManageConnectionTypes onTypesChanged={handleDataChanged} />
      </Modal>
    </div>
  );
}

export default App;
