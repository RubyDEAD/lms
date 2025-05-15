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
    const [borrowRecords, setBorrowRecords] = useState([]);
    const [books, setBooks] = useState([]);
    const [fines, setFines] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const BOOKS_API = "http://localhost:8081/books"; // Replace with your books API endpoint
    const BORROWED_BOOKS_API = "http://localhost:8081/borrowed-books"; // Replace with your borrowed books API endpoint
    const FINES_API = "http://localhost:8081/fines"; // Replace with your fines API endpoint

    useEffect(() => {
        const fetchDashboardData = async () => {
            try {
                const { data: { session } } = await supabase.auth.getSession();
                if (!session) throw new Error("Not authenticated");

                const token = session.access_token;
                const userId = session.user.id;

                // Fetch all data in parallel
                const [booksResponse, borrowedBooksResponse, finesResponse] = await Promise.all([
                    axios.get(BOOKS_API, {
                        headers: { Authorization: `Bearer ${token}` }
                    }),
                    axios.get(`${BORROWED_BOOKS_API}?patronId=${userId}`, {
                        headers: { Authorization: `Bearer ${token}` }
                    }),
                    axios.get(`${FINES_API}?patronId=${userId}`, {
                        headers: { Authorization: `Bearer ${token}` }
                    })
                ]);

                // Update state with fetched data
                const books = booksResponse.data;
                const borrowedBooks = borrowedBooksResponse.data;
                const fines = finesResponse.data;

                setBooks(books);
                setBorrowRecords(borrowedBooks);
                setFines(fines);

                // Calculate stats
                const booksAvailable = books.filter(book => book.availableCopies > 0).length;
                const booksBorrowed = borrowedBooks.length;
                const booksReturned = borrowedBooks.filter(record => record.status === "returned").length;
                const violations = fines.length;
                const totalFine = fines.reduce((sum, fine) => sum + fine.amount, 0);

                setStats({
                    booksAvailable,
                    booksBorrowed,
                    booksReturned,
                    violations,
                    totalFine
                });

                setError(null);
            } catch (err) {
                console.error("Dashboard fetch error:", err.response?.data || err.message);
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
                    <div className="dashboard-card-title">Violations</div>
                    <div className="dashboard-card-value">{stats.violations}</div>
                </div>
                <div className="dashboard-card">
                    <div className="dashboard-card-title">Total Fine</div>
                    <div className="dashboard-card-value">₱{stats.totalFine.toFixed(2)}</div>
                </div>
            </div>

            {/* Borrowed Books Section */}
            <div className="dashboard-section">
                <h2>Your Borrowed Books</h2>
                <ul>
                    {borrowRecords.map(record => (
                        <li key={record.id}>
                            Book ID: {record.bookId} | Status: {record.status} | Due: {new Date(record.dueDate).toLocaleDateString()}
                        </li>
                    ))}
                </ul>
            </div>

            {/* Books Section */}
            <div className="dashboard-section">
                <h2>All Books</h2>
                <ul>
                    {books.map(book => (
                        <li key={book.id}>
                            <strong>{book.title}</strong> by {book.authorName}
                        </li>
                    ))}
                </ul>
            </div>

            {/* Fines Section */}
            <div className="dashboard-section">
                <h2>Your Fines</h2>
                <ul>
                    {fines.map(fine => (
                        <li key={fine.fine_id}>
                            Book ID: {fine.bookId} | Amount: ₱{fine.amount.toFixed(2)} | Days Late: {fine.daysLate}
                        </li>
                    ))}
                </ul>
            </div>
        </div>
    );
}

export default Dashboard;