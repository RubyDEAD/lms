import React from 'react';

function AdminDashboard() {
  return (
    <div style={{ 
      position: 'fixed', // Ensures the content stays fixed in the viewport
      top: '50%', // Centers vertically
      left: '60%', // Centers horizontally
      transform: 'translate(-50%, -50%)', // Adjusts for the element's size
      width: '100%', // Ensures it spans the full width
      textAlign: 'center', // Centers text
      //backgroundColor: '#f8f9fa', // Background color
      padding: '2rem', // Padding for spacing
    }}>
      <div style={{ 
        display: 'inline-block', 
        padding: '4rem', 
        border: '1px solid #dee2e6', 
        borderRadius: '8px', 
        backgroundColor: '#ffffff', 
        boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)' 
      }}>
        <h1 style={{ marginBottom: '1rem', color: '#343a40' }}>Welcome to LMS</h1>
        <p style={{ fontSize: '1.2rem', color: '#6c757d' }}>Have fun searching for books!</p>
      </div>
    </div>
  );
}

export default AdminDashboard;
