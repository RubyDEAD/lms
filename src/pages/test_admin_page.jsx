import React from 'react';

function AdminDashboard() {
  return (
    <div style={{ 
      display: 'flex', 
      justifyContent: 'center', 
      alignItems: 'center', 
      height: '100vh', 
      backgroundColor: '#f8f9fa' 
    }}>
      <div style={{ 
        padding: '2rem', 
        border: '1px solid #dee2e6', 
        borderRadius: '8px', 
        backgroundColor: '#ffffff', 
        boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)' 
      }}>
        <h1 style={{ marginBottom: '1rem', color: '#343a40' }}>Admin Dashboard</h1>
        <p style={{ fontSize: '1.2rem', color: '#6c757d' }}>You are in the admin page.</p>
      </div>
    </div>
  );
}

export default AdminDashboard;
