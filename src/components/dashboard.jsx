import React, { useEffect, useState } from 'react';
import '../App.css';
import { supabase } from '../supabaseClient';
import axios from 'axios';

function Dashboard() {
    const [stats, setStats] = useState({
        booksAvailable: 0,
        booksBorrowed: 0,
        booksReturned: 0,
        violations: 0,
        totalFine: 0
    });
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const API_URL = "http://localhost:8081/query"; // GraphQL endpoint

    useEffect(() => {
        const fetchDashboardStats = async () => {
            try {
                const { data: { session } } = await supabase.auth.getSession();
                if (!session) throw new Error("Not authenticated");

                const response = await axios.post(API_URL, {
                    query: `
                        query DashboardStats {
                            booksAvailable
                            booksBorrowed
                            booksReturned
                            violations
                            totalFine
                        }
                    `
                }, {
                    headers: {
                        Authorization: `Bearer ${session.access_token}`
                    }
                });

                const data = response.data.data;
                setStats({
                    booksAvailable: data.booksAvailable,
                    booksBorrowed: data.booksBorrowed,
                    booksReturned: data.booksReturned,
                    violations: data.violations,
                    totalFine: data.totalFine
                });
                setError(null);
            } catch (err) {
                console.error("Dashboard fetch error:", err);
                setError("Failed to load dashboard data.");
            } finally {
                setLoading(false);
            }
        };

        fetchDashboardStats();
    }, []);

    if (loading) return <div className="dashboard-container">Loading dashboard...</div>;
    if (error) return <div className="dashboard-container">{error}</div>;

    return (
        <div className="dashboard-container">
            <h1 className="dashboard-title">DASHBOARD</h1>
            <div className="dashboard-row">
                <div className="dashboard-card">
                    <div className="dashboard-card-title">Books Available</div>
                    <div className="dashboard-card-value">{stats.booksAvailable}</div>
                </div>
                <div className="dashboard-card">
                    <div className="dashboard-card-title">Books Borrowed</div>
                    <div className="dashboard-card-value">{stats.booksBorrowed}</div>
                </div>
                <div className="dashboard-card">
                    <div className="dashboard-card-title">Books Returned</div>
                    <div className="dashboard-card-value">{stats.booksReturned}</div>
                </div>
                <div className="dashboard-card">
                    <div className="dashboard-card-title">Violations</div>
                    <div className="dashboard-card-value">{stats.violations}</div>
                </div>
                <div className="dashboard-card">
                    <div className="dashboard-card-title">Total Fine</div>
                    <div className="dashboard-card-value">{stats.totalFine}</div>
                </div>
            </div>
        </div>
    );
}

export default Dashboard;
