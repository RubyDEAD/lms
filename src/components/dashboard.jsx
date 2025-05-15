import React, { useEffect, useState } from "react";
import "../App.css";

function Dashboard() {
    const [stats, setStats] = useState({
        booksAvailable: 0,
        booksBorrowed: 0,
        booksReturned: 0, // Added books returned
        totalFines: 0,
        warningCount: 0,
        unpaidFees: 0, // Added unpaid fees
    });
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchDashboardData = async () => {
            try {
                // Simulate fetching data with mock data
                const books = [
                    { id: 1, title: "Book 1", availableCopies: 3 },
                    { id: 2, title: "Book 2", availableCopies: 0 },
                    { id: 3, title: "Book 3", availableCopies: 5 },
                ];

                const borrowedBooks = [
                    { id: 1, bookId: 1, status: "borrowed" },
                    { id: 2, bookId: 2, status: "returned" },
                ];

                const fines = [
                    { fine_id: 1, amount: 50 },
                    { fine_id: 2, amount: 30 },
                ];

                const profile = {
                    status: {
                        warning_count: 2,
                        unpaid_fees: 100, // Added unpaid fees
                    },
                };

                // Calculate stats
                const booksAvailable = books.filter(book => book.availableCopies > 0).length;
                const booksBorrowed = borrowedBooks.filter(record => record.status === "borrowed").length;
                const booksReturned = borrowedBooks.filter(record => record.status === "returned").length; // Calculate books returned
                const totalFines = fines.reduce((sum, fine) => sum + fine.amount, 0);
                const warningCount = profile.status.warning_count;
                const unpaidFees = profile.status.unpaid_fees;

                setStats({
                    booksAvailable,
                    booksBorrowed,
                    booksReturned, // Set books returned
                    totalFines,
                    warningCount,
                    unpaidFees, // Set unpaid fees
                });

                setError(null);
            } catch (err) {
                console.error("Dashboard fetch error:", err.message || err);
                setError("Failed to load dashboard data.");
            } finally {
                setLoading(false);
            }
        };

        fetchDashboardData();
    }, []);

    if (loading) return <div className="dashboard-container">Loading dashboard...</div>;
    if (error) return <div className="dashboard-container">{error}</div>;

    return (
        <div className="dashboard-container">
            <h1 className="dashboard-title">DASHBOARD</h1>

            {/* Stats Section */}
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
                    <div className="dashboard-card-title">Warning Count</div>
                    <div className="dashboard-card-value">{stats.warningCount}</div>
                </div>
                <div className="dashboard-card">
                    <div className="dashboard-card-title">Total Fines</div>
                    <div className="dashboard-card-value">₱{stats.totalFines.toFixed(2)}</div>
                </div>
                <div className="dashboard-card">
                    <div className="dashboard-card-title">Unpaid Fees</div>
                    <div className="dashboard-card-value">₱{stats.unpaidFees.toFixed(2)}</div>
                </div>
            </div>
        </div>
    );
}

export default Dashboard;