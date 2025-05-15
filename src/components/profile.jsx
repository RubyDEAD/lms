import React, { useEffect, useState } from "react";
import axios from "axios";

function Profile() {
  const API_URL = "http://localhost:8081/query";

  const [patron, setPatron] = useState(null);

  useEffect(() => {
    const userOBJ = localStorage.getItem("user");
    const user = JSON.parse(userOBJ);

    const query = `
    query {
      getPatronById(patron_id: "${user.id}") {
        first_name
        last_name
        patron_id
        status {
          unpaid_fees
          patron_status
          warning_count
        }
        membership {
          membership_id
          level
        }
      }
    }
  `;

    const getPatron = async () => {
      try {
        const response = await axios.post(API_URL, { query });
        const data = response.data.data.getPatronById;
        setPatron(data);
      } catch (err) {
        console.error("Error fetching user: ", err);
      }
    };

    getPatron();
  }, []);

  return (
    <div className="body">
      <div className="container">
        <h1>Profile</h1>
        {patron ? (
          <>
            <h3>
              User: {patron.first_name} {patron.last_name}
            </h3>
            <div>
              <p>Status:</p>
              <p>Unpaid Fees: {patron.status.unpaid_fees}</p>
              <p>Status: {patron.status.patron_status}</p>
              <p>Warning Count: {patron.status.warning_count}</p>
            </div>
            <div>
              <p>Membership: </p>
              <p>Membership ID: {patron.membership.membership_id}</p>
              <p>Level: {patron.membership.level}</p>
            </div>
          </>
        ) : (
          <h3>Loading...</h3> 
        )}
      </div>
    </div>
  );
}

export default Profile;