import React from 'react';

const Topbar = () => {
  return (
    <header className="topbar" style={{
      position: 'sticky',
      top: 0,
      zIndex: 50,
      width: '100%',
      backgroundColor: 'black',
      boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
      padding: '1rem 0',
      marginBottom: '2rem' // Space below the topbar
    }}>
      <div style={{
        maxWidth: '1200px',
        margin: '0 auto',
        padding: '0 1.5rem',
        display: 'flex',
        alignItems: 'center',
        height: '100%'
      }}>
        <h1 style={{
          margin: 0,
          fontSize: '1.5rem',
          fontWeight: 600,
          color: 'white'
        }}>
          Library Management System
        </h1>
      </div>
    </header>
  );
};

export default Topbar;