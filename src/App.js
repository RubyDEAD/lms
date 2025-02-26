import React from 'react';
import './App.css';
import { BrowserRouter as Router } from 'react-router-dom';  // Import BrowserRouter
import Dashboard from './components/dashboard';
import Sidebar from './components/sidebar';

function App() {
  return (
    <Router> 
      <div className="App">
        <div style={{ display: 'flex' }}>
          <Sidebar />

          <div style={{ marginLeft: '250px', padding: '20px', flex: 1 }}>
            <Dashboard />
          </div>
        </div>
      </div>
    </Router> 
  );
}

export default App;
