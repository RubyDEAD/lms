import React from "react";
import { Link, useNavigate } from "react-router-dom"; // useNavigate for redirect
import { FaBook, FaDiceD6, FaHandHolding, FaHandHoldingHand, FaMoneyBill, FaRegUser, FaRightFromBracket } from "react-icons/fa6";
import { supabase } from "../supabaseClient"; // adjust path if necessary

function Sidebar() {
  const navigate = useNavigate();

  const handleLogout = async () => {
    const { error } = await supabase.auth.signOut();

    if (error) {
      console.error("Error during logout:", error.message);
    } else {
      // Optionally clear local storage or other state
      // localStorage.clear();

      // Refresh the page
      window.location.href = "/login";
    }
  };

  return (
    <aside className="sidebar">
      <div className="sidebar-container">
        <div className="sidebar-nav">
          <ul>
            <li>
              <Link to="/books" className="sidebar-link">
                <FaBook /> Books
              </Link>
            </li>
            <li>
              <Link to="/borrowed-books" className="sidebar-link">
                <FaHandHolding /> Borrowed Books
              </Link>
            </li>
            {/* <li>
              <Link to="/add-book" className="sidebar-link">
                <FaHandHoldingHand /> Returned Books
              </Link>
            </li> */}
            <li>
              <Link to="/fines" className="sidebar-link">
                <FaMoneyBill /> Fines
              </Link>
            </li>
            <li>
              <Link to="/profile" className="sidebar-link">
                <FaRegUser /> Profile
              </Link>
            </li>
            <li>
              <button className="sidebar-link btn btn-link text-danger" onClick={handleLogout}>
                <FaRightFromBracket /> Logout
              </button>
            </li>
          </ul>
        </div>
      </div>
    </aside>
  );
}

export default Sidebar;
