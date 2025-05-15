import React, { useEffect, useState } from "react";
import axios from "axios";
import "../App.css"; // Ensure this file contains the necessary styles

function Profile() {
  const API_URL = "http://localhost:8081/query";
  const [patron, setPatron] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const userOBJ = localStorage.getItem("user");
    const user = JSON.parse(userOBJ);

    const query = `
      query {
        getPatronById(patron_id: "${user.id}") {
          first_name
          last_name
          status {
            unpaid_fees
            patron_status
            warning_count
          }
        }
      }
    `;

    const getPatron = async () => {
      try {
        const response = await axios.post(API_URL, { query });
        const data = response.data.data.getPatronById;
        setPatron(data);
        setError(null);
      } catch (err) {
        console.error("Error fetching user: ", err);
        setError("Failed to load profile data.");
      } finally {
        setLoading(false);
      }
    };

    getPatron();
  }, []);

  if (loading) return <div className="profile-container">Loading profile...</div>;
  if (error) return <div className="profile-container error">{error}</div>;

  return (
    <div className="profile-container">
      <div className="profile-card">
        <h1 className="profile-title">Profile</h1>
        {patron && (
          <>
            <div className="profile-section">
              <h2>Personal Information</h2>
              <p><strong>Name:</strong> {patron.first_name} {patron.last_name}</p>
            </div>

            <div className="profile-section">
              <h2>Status</h2>
              <p><strong>Unpaid Fees:</strong> ₱{patron.status.unpaid_fees.toFixed(2)}</p>
              <p><strong>Status:</strong> {patron.status.patron_status}</p>
              <p><strong>Warning Count:</strong> {patron.status.warning_count}</p>
            </div>
          </>
        )}
      </div>
    </div>
  );
}

export default Profile;