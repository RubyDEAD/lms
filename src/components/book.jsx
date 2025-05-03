import React, { useEffect, useState } from "react";
import axios from "axios";
import { supabase } from '../supabaseClient';
import { useNavigate } from "react-router-dom";
import 'bootstrap/dist/css/bootstrap.min.css';

function Books() {
    const [books, setBooks] = useState([]);
    const [bookDetails, setBookDetails] = useState(null);
    const [bookCopies, setBookCopies] = useState([]);
    const [availableCopy, setAvailableCopy] = useState(null);
    const [loading, setLoading] = useState(true);
    const [authLoading, setAuthLoading] = useState(true);
    const [error, setError] = useState(null);
    const [showAddForm, setShowAddForm] = useState(false);
    const [newBook, setNewBook] = useState({
        title: "",
        authorName: "",
        datePublished: "",
        description: "",
    });
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
                await fetchBooks();
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

    const fetchBooks = async () => {
        try {
            setLoading(true);
            const { data: { session } } = await supabase.auth.getSession();
            const token = session?.access_token;

            const response = await axios.post(API_URL, {
                query: `
                    query {
                        getBooks {
                            id
                            title
                            author_name
                            date_published
                            description
                        }
                    }
                `,
            }, {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            });

            setBooks(response.data.data.getBooks);
            setError(null);
        } catch (err) {
            console.error("Error fetching books:", err);
            setError("Failed to fetch books. Please try again later.");
            if (err.response?.status === 401) {
                navigate('/login');
            }
        } finally {
            setLoading(false);
        }
    };

    const addBook = async () => {
        try {
            const { data: { session } } = await supabase.auth.getSession();
            if (!session) {
                navigate('/login');
                return;
            }

            const token = session.access_token;

            if (!newBook.title.trim() || !newBook.authorName.trim()) {
                setError("Title and Author Name are required");
                return;
            }

            const response = await axios.post(API_URL, {
                query: `
                    mutation AddBook($title: String!, $authorName: String!, $datePublished: String!, $description: String!) {
                        addBook(
                            title: $title,
                            authorName: $authorName,
                            datePublished: $datePublished,
                            description: $description
                        ) {
                            id
                            title
                            author_name
                            date_published
                            description
                        }
                    }
                `,
                variables: newBook
            }, {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            });

            const addedBook = response.data.data.addBook;
            if (addedBook) {
                setBooks((prevBooks) => [...prevBooks, addedBook]);
                setNewBook({ title: "", authorName: "", datePublished: "", description: "" });
                setShowAddForm(false);
                setError(null);
            }
        } catch (err) {
            console.error("Error adding book:", err);
            setError(err.response?.data?.errors?.[0]?.message || "Failed to add book. Please try again later.");
        }
    };

    const fetchBookById = async (id) => {
        try {
            const { data: { session } } = await supabase.auth.getSession();
            if (!session) return navigate('/login');
            const token = session.access_token;

            const response = await axios.post(API_URL, {
                query: `
                    query GetBookById($id: ID!) {
                        getBookById(id: $id) {
                            id
                            title
                            author_name
                            date_published
                            description
                        }
                    }
                `,
                variables: { id }
            }, {
                headers: { Authorization: `Bearer ${token}` }
            });

            setBookDetails(response.data.data.getBookById);
            setAvailableCopy(null);
            setBookCopies([]);
            setError(null);
        } catch (err) {
            console.error("Error fetching book by ID:", err);
            setError("Failed to fetch book details. Please try again later.");
        }
    };

    const fetchBookCopiesById = async (id) => {
        try {
            const { data: { session } } = await supabase.auth.getSession();
            if (!session) return navigate('/login');
            const token = session.access_token;

            const response = await axios.post(API_URL, {
                query: `
                    query GetBookCopiesById($id: ID!) {
                        getBookCopiesById(id: $id) {
                            id
                            book_id
                            book_status
                        }
                    }
                `,
                variables: { id }
            }, {
                headers: { Authorization: `Bearer ${token}` }
            });

            setBookCopies(response.data.data.getBookCopiesById);
            setBookDetails(null);
            setAvailableCopy(null);
            setError(null);
        } catch (err) {
            console.error("Error fetching book copies:", err);
            setError("Failed to fetch book copies. Please try again later.");
        }
    };

    const fetchAvailableBookCopyById = async (id) => {
        try {
            const { data: { session } } = await supabase.auth.getSession();
            if (!session) return navigate('/login');
            const token = session.access_token;

            const response = await axios.post(API_URL, {
                query: `
                    query GetAvailableBookCopyByID($id: ID!) {
                        getAvailbleBookCopyByID(id: $id) {
                            id
                            book_id
                            book_status
                        }
                    }
                `,
                variables: { id }
            }, {
                headers: { Authorization: `Bearer ${token}` }
            });

            const availableCopy = response.data.data.getAvailbleBookCopyByID;

            if (availableCopy) {
                setAvailableCopy(availableCopy);
                setBookDetails(null);
                setBookCopies([]);
                setError(null);
            } else {
                setAvailableCopy(null);
                setError("No available copy for this book.");
            }
        } catch (err) {
            console.error("Error fetching available book copy:", err);
            setAvailableCopy(null);
            setError("Failed to check availability. Please try again later.");
        }
    };

    if (authLoading) return <div className="container mt-5">Checking authentication...</div>;

    if (!isAuthenticated) {
        return (
            <div className="container mt-5">
                <div className="alert alert-warning">Please log in to access the books library.</div>
            </div>
        );
    }

    if (loading) return <div className="container mt-5">Loading books...</div>;

    return (
        <div className="body">
            <div className="container mt-4">
                <h1 className="mb-4">Books Library</h1>

                {error && (
                    <div className="alert alert-danger mb-4">
                        {error}
                        <button className="btn-close float-end" onClick={() => setError(null)}></button>
                    </div>
                )}

                {/* Toggle Add Book Form */}
                <button
                    className="btn btn-outline-primary mb-3"
                    onClick={() => setShowAddForm(!showAddForm)}
                >
                    {showAddForm ? "Cancel" : "Add Book"}
                </button>

                {/* Add Book Form */}
                {showAddForm && (
                    <div className="card mb-4">
                        <div className="card-body">
                            <h2 className="card-title">Add a New Book</h2>
                            <form onSubmit={(e) => { e.preventDefault(); addBook(); }}>
                                <div className="mb-3">
                                    <label htmlFor="title" className="form-label">Title *</label>
                                    <input type="text" className="form-control" id="title"
                                        value={newBook.title}
                                        onChange={(e) => setNewBook({ ...newBook, title: e.target.value })}
                                        required />
                                </div>
                                <div className="mb-3">
                                    <label htmlFor="authorName" className="form-label">Author Name *</label>
                                    <input type="text" className="form-control" id="authorName"
                                        value={newBook.authorName}
                                        onChange={(e) => setNewBook({ ...newBook, authorName: e.target.value })}
                                        required />
                                </div>
                                <div className="mb-3">
                                    <label htmlFor="datePublished" className="form-label">Date Published</label>
                                    <input type="date" className="form-control" id="datePublished"
                                        value={newBook.datePublished}
                                        onChange={(e) => setNewBook({ ...newBook, datePublished: e.target.value })} />
                                </div>
                                <div className="mb-3">
                                    <label htmlFor="description" className="form-label">Description</label>
                                    <textarea className="form-control" id="description" rows="3"
                                        value={newBook.description}
                                        onChange={(e) => setNewBook({ ...newBook, description: e.target.value })}></textarea>
                                </div>
                                <button type="submit" className="btn btn-primary">Add Book</button>
                            </form>
                        </div>
                    </div>
                )}

                {/* Books List */}
                <div className="row row-cols-1 row-cols-md-2 row-cols-lg-3 g-4">
                    {books.map((book) => (
                        <div className="col" key={book.id}>
                            <div className="card h-100 shadow-sm">
                                <div className="card-body">
                                    <h5 className="card-title">{book.title}</h5>
                                    <h6 className="card-subtitle mb-2 text-muted">{book.author_name}</h6>
                                    <p className="card-text text-truncate">{book.description}</p>
                                    <p className="card-text">
                                        <small className="text-muted">Published: {book.date_published}</small>
                                    </p>
                                </div>
                                <div className="card-footer bg-transparent">
                                    <button className="btn btn-sm btn-outline-primary me-2"
                                        onClick={() => fetchBookById(book.id)}>Details</button>
                                    <button className="btn btn-sm btn-outline-secondary me-2"
                                        onClick={() => fetchBookCopiesById(book.id)}>Copies</button>
                                    <button className="btn btn-sm btn-outline-success"
                                        onClick={() => fetchAvailableBookCopyById(book.id)}>Availability</button>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>

                {/* Modals */}
                {bookDetails && (
                    <div className="modal fade show d-block" style={{ backgroundColor: 'rgba(0,0,0,0.5)' }}>
                        <div className="modal-dialog">
                            <div className="modal-content">
                                <div className="modal-header">
                                    <h5 className="modal-title">{bookDetails.title}</h5>
                                    <button className="btn-close" onClick={() => setBookDetails(null)}></button>
                                </div>
                                <div className="modal-body">
                                    <p><strong>Author:</strong> {bookDetails.author_name}</p>
                                    <p><strong>Published:</strong> {bookDetails.date_published}</p>
                                    <p><strong>Description:</strong> {bookDetails.description}</p>
                                </div>
                                <div className="modal-footer">
                                    <button className="btn btn-secondary" onClick={() => setBookDetails(null)}>Close</button>
                                </div>
                            </div>
                        </div>
                    </div>
                )}

                {bookCopies.length > 0 && (
                    <div className="modal fade show d-block" style={{ backgroundColor: 'rgba(0,0,0,0.5)' }}>
                        <div className="modal-dialog">
                            <div className="modal-content">
                                <div className="modal-header">
                                    <h5 className="modal-title">Book Copies</h5>
                                    <button className="btn-close" onClick={() => setBookCopies([])}></button>
                                </div>
                                <div className="modal-body">
                                    <ul className="list-group">
                                        {bookCopies.map((copy) => (
                                            <li key={copy.id} className="list-group-item d-flex justify-content-between align-items-center">
                                                Copy #{copy.id}
                                                <span className={`badge ${copy.book_status === 'available' ? 'bg-success' : 'bg-warning'} text-dark`}>
                                                    {copy.book_status}
                                                </span>
                                            </li>
                                        ))}
                                    </ul>
                                </div>
                                <div className="modal-footer">
                                    <button className="btn btn-secondary" onClick={() => setBookCopies([])}>Close</button>
                                </div>
                            </div>
                        </div>
                    </div>
                )}

                {availableCopy !== null && (
                    <div className="modal fade show d-block" style={{ backgroundColor: 'rgba(0,0,0,0.5)' }}>
                        <div className="modal-dialog">
                            <div className="modal-content">
                                <div className="modal-header">
                                    <h5 className="modal-title">Copy Availability</h5>
                                    <button className="btn-close" onClick={() => setAvailableCopy(null)}></button>
                                </div>
                                <div className="modal-body">
                                    {availableCopy ? (
                                        <div className="alert alert-success">
                                            <h6>Available Copy Found!</h6>
                                            <p>Copy ID: {availableCopy.id}</p>
                                            <p>Status: <span className="badge bg-success">{availableCopy.book_status}</span></p>
                                        </div>
                                    ) : (
                                        <div className="alert alert-warning">
                                            No available copies for this book at the moment.
                                        </div>
                                    )}
                                </div>
                                <div className="modal-footer">
                                    <button className="btn btn-secondary" onClick={() => setAvailableCopy(null)}>Close</button>
                                </div>
                            </div>
                        </div>
                    </div>
                )}

            </div>
        </div>
    );
}

export default Books;
