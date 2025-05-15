import React, { useState, useEffect } from "react";
import axios from "axios";
import { supabase } from '../supabaseClient';
import { useNavigate } from "react-router-dom";
import 'bootstrap/dist/css/bootstrap.min.css';

const BorrowedBooks = () => {
    const [borrowRecords, setBorrowRecords] = useState([]);
    const [booksData, setBooksData] = useState({});
    const [loading, setLoading] = useState(true);
    const [authLoading, setAuthLoading] = useState(true);
    const [error, setError] = useState(null);
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const navigate = useNavigate();

    const API_URL = "http://localhost:8081/query";

    useEffect(() => {
        const checkAuth = async () => {
            try {
                setAuthLoading(true);
                const { data: { session }, error } = await supabase.auth.getSession();

                if (error) throw error;

                setIsAuthenticated(!!session);
                if (!session) {
                    navigate('/login');
                    return;
                }

                await fetchBorrowRecords();
            } catch (err) {
                console.error("Authentication error:", err);
                setError("Session expired. Please log in again.");
                navigate('/login');
            } finally {
                setAuthLoading(false);
            }
        };

        checkAuth();
    }, [navigate]);

    const fetchBookDetails = async (bookIds) => {
        try {
            const { data: { session } } = await supabase.auth.getSession();
            if (!session) return;

            const token = session.access_token;
            const uniqueBookIds = [...new Set(bookIds)].filter(id => !!id);

            const bookDetails = {};

            await Promise.all(uniqueBookIds.map(async (bookId) => {
                try {
                    const response = await axios.post(
                        API_URL,
                        {
                            query: `
                                query GetBookById($id: String!) {
                                    getBookById(id: $id) {
                                        id
                                        title
                                        author_name
                                    }
                                }
                            `,
                            variables: { id: bookId }
                        },
                        {
                            headers: {
                                Authorization: `Bearer ${token}`,
                                'Content-Type': 'application/json'
                            }
                        }
                    );

                    if (response.data?.data?.getBookById) {
                        bookDetails[bookId] = response.data.data.getBookById;
                    }
                } catch (innerErr) {
                    console.error(`Error fetching book ID ${bookId}:`, innerErr?.response?.data || innerErr.message);
                }
            }));

            setBooksData(bookDetails);
        } catch (err) {
            console.error("Error fetching book details:", err);
        }
    };

    const fetchBorrowRecords = async () => {
        try {
            setLoading(true);
            const { data: { session } } = await supabase.auth.getSession();

            if (!session) {
                navigate('/login');
                return;
            }

            const token = session.access_token;
            const userId = session.user.id;

            const response = await axios.post(
                API_URL,
                {
                    query: `
                        query PatronBorrowHistory($patronId: ID!) {
                            patronBorrowHistory(patronId: $patronId) {
                                id
                                bookId
                                patronId
                                bookCopyId
                                borrowedAt
                                dueDate
                                returnedAt
                                renewalCount
                                status
                            }
                        }
                    `,
                    variables: { patronId: userId }
                },
                {
                    headers: {
                        Authorization: `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    }
                }
            );

            const records = response.data.data.patronBorrowHistory;
            setBorrowRecords(records);

            const bookIds = records.map(record => record.bookId);
            if (bookIds.length > 0) {
                await fetchBookDetails(bookIds);
            }

            setError(null);
        } catch (err) {
            console.error("Error fetching borrow records:", err);
            setError("Failed to fetch borrow records. Please try again later.");
            if (err.response?.status === 401) {
                navigate('/login');
            }
        } finally {
            setLoading(false);
        }
    };

    const returnBook = async (recordId) => {
        try {
            const { data: { session } } = await supabase.auth.getSession();
            if (!session) {
                setError("You must be logged in to return a book.");
                return;
            }

            const token = session.access_token;

            const response = await axios.post(
                API_URL,
                {
                    query: `
                        mutation ReturnBook($recordId: ID!) {
                            returnBook(recordId: $recordId) {
                                id
                                returnedAt
                                status
                            }
                        }
                    `,
                    variables: { recordId }
                },
                {
                    headers: {
                        Authorization: `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    }
                }
            );

            const result = response.data.data.returnBook;
            if (result) {
                setBorrowRecords(prevRecords =>
                    prevRecords.map(record =>
                        record.id === recordId ? {
                            ...record,
                            status: "RETURNED",
                            returnedAt: result.returnedAt
                        } : record
                    )
                );
                setError(null);
                alert("Book returned successfully!");
            } else {
                setError("Failed to return the book. Please try again later.");
            }
        } catch (err) {
            console.error("Error returning book:", err);
            setError("Failed to return the book. Please try again later.");
        }
    };

    if (authLoading) {
        return <div className="container mt-5">Checking authentication...</div>;
    }

    if (!isAuthenticated) {
        return (
            <div className="container mt-5">
                <div className="alert alert-warning">
                    Please log in to view your borrowed books.
                </div>
            </div>
        );
    }

    if (loading) {
        return <div className="container mt-5">Loading borrowed books...</div>;
    }

    return (
        <div className="container mt-4">
            <h1 className="mb-4">My Borrowed Books</h1>

            {error && (
                <div className="alert alert-danger mb-4">
                    {error}
                    <button
                        type="button"
                        className="btn-close float-end"
                        onClick={() => setError(null)}
                        aria-label="Close"
                    ></button>
                </div>
            )}

            <div className="table-responsive">
                <table className="table table-striped table-hover">
                    <thead className="table-dark">
                        <tr>
                            <th>Book Title</th>
                            <th>Author</th>
                            <th>Copy ID</th>
                            <th>Borrowed At</th>
                            <th>Due Date</th>
                            <th>Returned At</th>
                            <th>Status</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        {borrowRecords.length === 0 ? (
                            <tr>
                                <td colSpan="8" className="text-center py-4">
                                    You don't have any borrowed books.
                                </td>
                            </tr>
                        ) : (
                            borrowRecords.map(record => (
                                <tr key={record.id}>
                                    <td>{booksData[record.bookId]?.title || `Book ID: ${record.bookId}`}</td>
                                    <td>{booksData[record.bookId]?.author_name || '-'}</td>
                                    <td>{record.bookCopyId}</td>
                                    <td>{new Date(record.borrowedAt).toLocaleDateString()}</td>
                                    <td>{new Date(record.dueDate).toLocaleDateString()}</td>
                                    <td>{record.returnedAt ? new Date(record.returnedAt).toLocaleDateString() : '-'}</td>
                                    <td>
                                        <span className={`badge ${
                                            record.status === "ACTIVE" ? "bg-primary" :
                                            record.status === "RETURNED" ? "bg-success" :
                                            record.status === "OVERDUE" ? "bg-danger" : "bg-warning"
                                        }`}>
                                            {record.status}
                                        </span>
                                    </td>
                                    <td>
                                        {record.status !== "RETURNED" && (
                                            <button
                                                className="btn btn-sm btn-outline-danger"
                                                onClick={() => returnBook(record.id)}
                                            >
                                                Return
                                            </button>
                                        )}
                                    </td>
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default BorrowedBooks;
