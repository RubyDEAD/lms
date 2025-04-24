import React from "react";
import { Link } from "react-router-dom"; // To use routing for links
import { FaBook,FaDiceD6,FaHandHolding, FaHandHoldingHand, FaRegUser } from "react-icons/fa6";

function Sidebar() {
  return (
    <aside className="sidebar">
      <div className="sidebar-container">
        
        {/* <div className="sidebar-title">
          <h1>Libary Management System</h1>
        </div> */}

        
        <div className="sidebar-nav">
          <ul>
            <li>
              <Link to="/" className="sidebar-link">
                <FaDiceD6></FaDiceD6> Dashboard
              </Link>
            </li>
            <li>
              <Link to="/books" className="sidebar-link">
                <FaBook></FaBook> Books
              </Link>
            </li>
            <li>
              <Link to="/borrowed-books" className="sidebar-link">
                <FaHandHolding></FaHandHolding> Borrowed Books
              </Link>
            </li>
            <li>
              <Link to="/add-book" className="sidebar-link">
                <FaHandHoldingHand></FaHandHoldingHand> Returned Books
              </Link>
            </li>
            <li>
              <Link to="/profile" className="sidebar-link">
                <FaRegUser></FaRegUser> Profile
              </Link>
            </li>
          </ul>
        </div>
      </div>
    </aside>
  );
}

export default Sidebar;
