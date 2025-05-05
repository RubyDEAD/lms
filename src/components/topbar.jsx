import React from 'react';
import { Link } from 'react-router-dom'; // For routing links
import { FaBell, FaCog, FaUserCircle } from 'react-icons/fa'; // Importing icons

const Topbar = () => {
  return (
    <header className="topbar">
      <div className="topbar-container">
        {/* Logo and title */}
        <div className="top-left">
          <h1>Library Management System</h1>
        </div>

        {/* Right side icons and profile */}
        <div className="top-right">

        </div>
      </div>
    </header>
  );
};

export default Topbar;
