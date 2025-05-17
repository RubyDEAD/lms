import React, { useEffect, useState } from "react";
import axios from "axios";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faUserCircle } from "@fortawesome/free-solid-svg-icons";

function Profile() {
  const API_URL = "http://localhost:8081/query";
  const [patron, setPatron] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const userOBJ = localStorage.getItem("user");
    const user = JSON.parse(userOBJ);

    const query = `
      query {
        getPatronById(patron_id: "${user.id}") {
          first_name
          last_name
          patron_id
        }
      }
    `;

    const getPatron = async () => {
      try {
        const response = await axios.post(API_URL, { query });
        const data = response.data.data.getPatronById;
        setPatron(data);
        setTimeout(() => setIsLoading(false), 800); // Smooth loading transition
      } catch (err) {
        console.error("Error fetching user: ", err);
        setIsLoading(false);
      }
    };

    getPatron();
  }, []);

  const getStatusColor = (status) => {
    switch (status.toLowerCase()) {
      case "active":
        return "bg-emerald-100 text-emerald-800";
      case "suspended":
        return "bg-rose-100 text-rose-800";
      case "restricted":
        return "bg-amber-100 text-amber-800";
      default:
        return "bg-blue-100 text-blue-800";
    }
  };

  return (
    <div className="mt-20 min-h-screen bg-gradient-to-br from-gray-50 to-indigo-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md mx-auto transform transition-all duration-300 hover:scale-[1.01]">
        {/* Profile Card */}
        <div className="bg-white rounded-2xl shadow-xl overflow-hidden transition-all duration-500 ease-in-out">
          {/* Profile Header with Floating Effect */}
          <div className="relative h-48 bg-gradient-to-r from-indigo-400 to-purple-500 overflow-hidden">
            <div className="absolute inset-0 bg-gradient-to-t from-black/10 to-transparent"></div>
            <div className="absolute -bottom-16 left-1/2 transform -translate-x-1/2 transition-all duration-500 hover:-translate-y-2">
              <div className="h-32 w-32 rounded-full bg-white/90 backdrop-blur-sm flex items-center justify-center shadow-xl border-4 border-white/30 hover:border-indigo-200/50 transition-all duration-300">
                <FontAwesomeIcon icon={faUserCircle} size="3x" className="text-indigo-500" />
              </div>
            </div>
          </div>

          {/* Profile Content */}
          <div className="px-6 pt-20 pb-8">
            {isLoading ? (
              <div className="space-y-6 animate-pulse">
                <div className="h-8 bg-gray-200 rounded-full w-3/4 mx-auto"></div>
                <div className="h-4 bg-gray-200 rounded-full w-1/2 mx-auto"></div>
                <div className="space-y-4 pt-6">
                  {[...Array(4)].map((_, i) => (
                    <div key={i} className="space-y-2">
                      <div className="h-3 bg-gray-200 rounded-full w-1/4"></div>
                      <div className="h-8 bg-gray-100 rounded-lg"></div>
                    </div>
                  ))}
                </div>
              </div>
            ) : patron ? (
              <div className="space-y-6">
                {/* Name Section */}
                <div className="text-center">
                  <h1 className="text-2xl font-bold text-gray-800 tracking-tight">
                    {patron.first_name} {patron.last_name}
                  </h1>
                  <p className="text-indigo-400 mt-1 text-sm font-medium">Library Patron</p>
                </div>

                {/* Stats Grid */}
                {/* <div className="grid grid-cols-3 gap-4 pt-4">
                  <div className="bg-indigo-50/50 p-3 rounded-xl text-center border border-indigo-100">
                    <p className="text-xs text-indigo-500 font-medium">Status</p>
                    <span className={`text-sm font-semibold ${getStatusColor(patron.status.patron_status).replace('bg-', '')}`}>
                      {patron.status.patron_status}
                    </span>
                  </div>
                  <div className="bg-indigo-50/50 p-3 rounded-xl text-center border border-indigo-100">
                    <p className="text-xs text-indigo-500 font-medium">Fees</p>
                    <p className="text-sm font-semibold text-gray-700">
                      ${patron.status.unpaid_fees.toFixed(2)}
                    </p>
                  </div>
                  <div className="bg-indigo-50/50 p-3 rounded-xl text-center border border-indigo-100">
                    <p className="text-xs text-indigo-500 font-medium">Warnings</p>
                    <p className="text-sm font-semibold text-gray-700">
                      {patron.status.warning_count}
                    </p>
                  </div>
                </div> */}

                {/* Details Section */}
                <div className="space-y-4 pt-2">
                  <div className="p-4 bg-gray-50/50 rounded-xl border border-gray-200">
                    <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">
                      Patron Information
                    </h3>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-600">Member ID</span>
                      <span className="font-mono text-sm text-indigo-600 bg-indigo-50 px-2 py-1 rounded">
                        {patron.patron_id}
                      </span>
                    </div>
                  </div>

                  {/* <div className="p-4 bg-gray-50/50 rounded-xl border border-gray-200">
                    <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">
                      Account Details
                    </h3>
                    <div className="space-y-3">
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-gray-600">Membership</span>
                        <span className={`text-xs font-medium px-2 py-1 rounded-full ${getStatusColor(patron.status.patron_status)}`}>
                          {patron.status.patron_status}
                        </span>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-gray-600">Last Updated</span>
                        <span className="text-xs text-gray-500">
                          {new Date().toLocaleDateString()}
                        </span>
                      </div>
                    </div>
                  </div> */}
                </div>
              </div>
            ) : (
              <div className="text-center py-8">
                <p className="text-gray-500">Failed to load profile data</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default Profile;
