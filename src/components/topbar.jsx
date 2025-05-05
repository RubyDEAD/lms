import React, { useState } from 'react';
import { FaBell, FaCog, FaUserCircle } from 'react-icons/fa';
import '../App.css'; // Make sure to style accordingly

const Topbar = () => {
  const [showNotifications, setShowNotifications] = useState(false);

  const toggleNotifications = () => {
    setShowNotifications(prev => !prev);
  };

  // Dummy notification data
  const notifications = [
    "New book 'React Basics' added.",
    "Overdue: 'Modern JS Guide'.",
    "Library maintenance on Friday.",
  ];

  return (
    <header className="topbar">
      <div className="topbar-container">
        {/* Logo and title */}
        <div className="top-left">
          <h1>Library Management System</h1>
        </div>

        {/* Right side icons and profile */}
        <div className="top-right">
          <div className="icon-wrapper" onClick={toggleNotifications}>
            <FaBell className="topbar-icon" />
            {notifications.length > 0 && <span className="notification-badge">{notifications.length}</span>}
          </div>

          <FaUserCircle className="topbar-icon" />
          <FaCog className="topbar-icon" />

          {/* Notifications Panel */}
          {showNotifications && (
            <div className="notification-panel">
              <h6>Notifications</h6>
              <ul>
                {notifications.map((note, idx) => (
                  <li key={idx}>{note}</li>
                ))}
              </ul>
            </div>
          )}
        </div>
      </div>
    </header>
  );
};

export default Topbar;
