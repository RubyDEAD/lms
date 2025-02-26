import React from "react";
import { Link } from "react-router-dom"; // To use routing for links

function Sidebar() {
  return (
    <aside className="sidebar">
      <div className="sidebar-container">
        
        <div className="sidebar-title">
          <h1>Libary Management System</h1>
        </div>

        
        <div className="sidebar-nav">
          <ul>
            <li>
              <Link to="/" className="sidebar-link">
                Dashboard
              </Link>
            </li>
            <li>
              <Link to="/books" className="sidebar-link">
                Books
              </Link>
            </li>
            <li>
              <Link to="/borrowed-books" className="sidebar-link">
                Borrowed Books
              </Link>
            </li>
            <li>
              <Link to="/add-book" className="sidebar-link">
                Add Book
              </Link>
            </li>
            <li>
              <Link to="/profile" className="sidebar-link">
                Profile
              </Link>
            </li>
          </ul>
        </div>
      </div>
    </aside>
  );
}

export default Sidebar;
