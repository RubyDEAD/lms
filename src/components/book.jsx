import React, { useEffect, useState } from "react";
import axios from "axios";
import { supabase, bookservice, bucketURL } from '../supabaseClient';
import { useNavigate } from "react-router-dom";
import 'bootstrap/dist/css/bootstrap.min.css';

function Books() {
    const [borrowRecords, setBorrowRecords] = useState([]);
    const [showBorrowRecordsModal, setShowBorrowRecordsModal] = useState(false);
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
        imageFile: null
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
                            image
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
    
            let imageUrl = null;
    
            if (newBook.imageFile) {
                const fileName = `${Date.now()}_${newBook.imageFile.name}`;
                const { data, error } = await bookservice.storage
                    .from('test')
                    .upload(fileName, newBook.imageFile);
    
                if (error) {
                    console.error("Error uploading image:", error);
                    setError("Failed to upload image. Please try again later.");
                    return;
                }
    
                imageUrl = `${bucketURL}${fileName}`;
            }
    
            const response = await axios.post(API_URL, {
                query: `
                    mutation AddBook($title: String!, $authorName: String!, $datePublished: String!, $description: String!, $image: String) {
                        addBook(
                            title: $title,
                            authorName: $authorName,
                            datePublished: $datePublished,
                            description: $description,
                            image: $image
                        ) {
                            id
                            title
                            author_name
                            date_published
                            description
                            image
                        }
                    }
                `,
                variables: {
                    ...newBook,
                    image: imageUrl
                }
            }, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
    
            const addedBook = response.data.data.addBook;
    
            if (addedBook) {
                setBooks((prevBooks) => [...prevBooks, addedBook]);
                setNewBook({ 
                    title: "", 
                    authorName: "", 
                    datePublished: "", 
                    description: "", 
                    imageFile: null 
                });
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
                    query getBookById($id: String!) {
                        getBookById(id: $id) {
                            id
                            title
                            author_name
                            date_published
                            description
                            image
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
                    query getBookCopiesById($id: String!) {
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
                    query getAvailableBookCopyByID($id: String!) {
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

    const fetchBorrowRecords = async (bookid) => {
        try {
            const { data: { session } } = await supabase.auth.getSession();
            if (!session) {
                setError("You must be logged in to view borrow records.");
                return;
            }
    
            const token = session.access_token;
    
            const response = await axios.post(
                API_URL,
                {
                    query: `
                        query getrecord($patronId: ID, $bookId: ID) {
                            borrowRecords(patronId: $patronId, bookId: $bookId) {
                                id
                                bookId
                                bookCopyId
                                returnedAt
                                borrowedAt
                                dueDate
                                status
                            }
                        }
                    `,
                    variables: {
                        patronId: session.user.id,
                        bookId: bookid,
                    },
                },
                {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                }
            );
    
            setBorrowRecords(response.data.data.borrowRecords);
            setShowBorrowRecordsModal(true);
        } catch (err) {
            console.error("Error fetching borrow records:", err);
            setError("Failed to fetch borrow records. Please try again later.");
        }
    };

    const borrowBook = async (bookId) => {
        try {
            const { data: { session } } = await supabase.auth.getSession();
            if (!session) {
                navigate('/login');
                return;
            }
    
            const token = session.access_token;
            const patronId = session.user.id;
    
            const response = await axios.post(API_URL, {
                query: `
                    mutation BorrowBook($bookId: ID!, $patronId: ID!) {
                        borrowBook(bookId: $bookId, patronId: $patronId) {
                            id
                            bookId
                            patronId
                            borrowedAt
                            dueDate
                            status
                            bookCopyId
                        }
                    }
                `,
                variables: {
                    bookId: bookId,
                    patronId: patronId
                }
            }, {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            });
    
            const result = response.data.data.borrowBook;
            if (result) {
                alert(`Book borrowed successfully! Due date: ${new Date(result.dueDate).toLocaleDateString()}`);
                // Refresh available copies
                fetchAvailableBookCopyById(bookId);
            }
        } catch (err) {
            console.error("Error borrowing book:", err);
            setError(err.response?.data?.errors?.[0]?.message || "You already borrowed this book.");
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
                    variables: {
                        recordId: recordId,
                    },
                },
                {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                }
            );
    
            const result = response.data.data.returnBook;
            if (result) {
                alert("Book returned successfully!");
                setBorrowRecords((prevRecords) =>
                    prevRecords.map((record) =>
                        record.id === recordId ? { ...record, status: "RETURNED", returnedAt: result.returnedAt } : record
                    )
                );
                // Refresh the books list to update availability
                fetchBooks();
            } else {
                setError("Failed to return the book. Please try again later.");
            }
        } catch (err) {
            console.error("Error returning book:", err);
            setError("Failed to return the book. Please try again later.");
        }
    };

    const handleImageChange = (e) => {
        if (e.target.files && e.target.files[0]) {
            setNewBook({ ...newBook, imageFile: e.target.files[0] });
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

                <div className="d-flex justify-content-between mb-4">
                    <button
                        className="btn btn-primary"
                        onClick={() => setShowAddForm(!showAddForm)}
                    >
                        {showAddForm ? "Cancel" : "+ Add Book"}
                    </button>
                    <button 
                        className="btn btn-outline-secondary"
                        onClick={() => fetchBorrowRecords()}
                    >
                        My Borrow Records
                    </button>
                </div>

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
                                <div className="mb-3">
                                    <label htmlFor="image" className="form-label">Book Cover Image</label>
                                    <input type="file" className="form-control" id="image"
                                        accept="image/*"
                                        onChange={handleImageChange} />
                                </div>
                                <button type="submit" className="btn btn-primary">Add Book</button>
                            </form>
                        </div>
                    </div>
                )}

                <div className="row row-cols-1 row-cols-md-2 row-cols-lg-3 g-4">
                    {books.map((book) => (
                        <div className="col" key={book.id}>
                            <div className="card h-100 shadow-sm">
                                <img
                                    src={
                                        book.image && new RegExp(bucketURL).test(book.image)
                                        ? book.image
                                        : "https://hwkuzfsecehszlftxqpn.supabase.co/storage/v1/object/public/test/default-book.png"
                                    }
                                    className="card-img-top"
                                    alt={book.title}
                                    style={{ height: '200px', objectFit: 'cover' }}
                                />
                                <div className="card-body">
                                    <h5 className="card-title">{book.title}</h5>
                                    <h6 className="card-subtitle mb-2 text-muted">{book.author_name}</h6>
                                    <p className="card-text text-truncate">{book.description}</p>
                                    <p className="card-text">
                                        <small className="text-muted">Published: {book.date_published}</small>
                                    </p>
                                </div>
                                <div className="card-footer bg-transparent d-flex justify-content-between">
                                    <button className="btn btn-sm btn-outline-primary"
                                        onClick={() => fetchBookById(book.id)}>
                                        Details
                                    </button>
                                    <button className="btn btn-sm btn-outline-success"
                                        onClick={() => fetchAvailableBookCopyById(book.id)}>
                                        Check Availability
                                    </button>
                                    <button className="btn btn-sm btn-outline-info"
                                        onClick={() => borrowBook(book.id)}>
                                        Borrow
                                    </button>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>

                {/* Modals */}
                {bookDetails && (
                    <div className="modal fade show d-block" style={{ backgroundColor: 'rgba(0,0,0,0.5)' }}>
                        <div className="modal-dialog modal-lg">
                            <div className="modal-content">
                                <div className="modal-header">
                                    <h5 className="modal-title">{bookDetails.title}</h5>
                                    <button className="btn-close" onClick={() => setBookDetails(null)}></button>
                                </div>
                                <div className="modal-body">
                                    <div className="row">
                                        <div className="col-md-4">
                                            <img
                                                src={
                                                    bookDetails.image && new RegExp(bucketURL).test(bookDetails.image)
                                                    ? bookDetails.image
                                                    : "https://hwkuzfsecehszlftxqpn.supabase.co/storage/v1/object/public/test/default-book.png"
                                                }                                            
                                                alt={bookDetails.title}
                                                className="img-fluid rounded"
                                                style={{ maxHeight: '300px', objectFit: 'cover' }}
                                            />
                                        </div>
                                        <div className="col-md-8">
                                            <p><strong>Author:</strong> {bookDetails.author_name}</p>
                                            <p><strong>Published:</strong> {bookDetails.date_published}</p>
                                            <p><strong>Description:</strong></p>
                                            <p>{bookDetails.description}</p>
                                            <div className="mt-3">
                                     
                            
                                            </div>
                                        </div>
                                    </div>
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
                        <div className="modal-dialog modal-lg">
                            <div className="modal-content">
                                <div className="modal-header">
                                    <h5 className="modal-title">Book Copies</h5>
                                    <button className="btn-close" onClick={() => setBookCopies([])}></button>
                                </div>
                                <div className="modal-body">
                                    <div className="table-responsive">
                                        <table className="table table-hover">
                                            <thead>
                                                <tr>
                                                    <th>Copy ID</th>
                                                    <th>Status</th>
                                                    <th>Actions</th>
                                                </tr>
                                            </thead>
                                            <tbody>
                                                {bookCopies.map((copy) => (
                                                    <tr key={copy.id}>
                                                        <td>#{copy.id}</td>
                                                        <td>
                                                            <span className={`badge ${copy.book_status === 'available' ? 'bg-success' : 'bg-warning'} text-dark`}>
                                                                {copy.book_status}
                                                            </span>
                                                        </td>
                                                        <td>
                                                            {copy.book_status === 'available' && (
                                                                <button 
                                                                    className="btn btn-sm btn-success"
                                                                    onClick={() => {
                                                                        setBookCopies([]);
                                                                        borrowBook(copy.book_id);
                                                                    }}
                                                                >
                                                                    Borrow
                                                                </button>
                                                            )}
                                                        </td>
                                                    </tr>
                                                ))}
                                            </tbody>
                                        </table>
                                    </div>
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
                                        <>
                                            <div className="alert alert-success">
                                                <h6>Available Copy Found!</h6>
                                                <p>Copy ID: {availableCopy.id}</p>
                                                <p>Status: <span className="badge bg-success">{availableCopy.book_status}</span></p>
                                            </div>
                                            <div className="d-flex justify-content-end">
        
                                            </div>
                                        </>
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

                {showBorrowRecordsModal && (
                    <div className="modal fade show d-block" style={{ backgroundColor: 'rgba(0,0,0,0.5)' }}>
                        <div className="modal-dialog modal-lg">
                            <div className="modal-content">
                                <div className="modal-header">
                                    <h5 className="modal-title">My Borrow Records</h5>
                                    <button className="btn-close" onClick={() => setShowBorrowRecordsModal(false)}></button>
                                </div>
                                <div className="modal-body">
                                    {borrowRecords.length > 0 ? (
                                        <div className="table-responsive">
                                            <table className="table table-hover">
                                                <thead>
                                                    <tr>
                                                        <th>Book ID</th>
                                                        <th>Copy ID</th>
                                                        <th>Borrowed At</th>
                                                        <th>Due Date</th>
                                                        <th>Status</th>
                                                        <th>Actions</th>
                                                    </tr>
                                                </thead>
                                                <tbody>
                                                    {borrowRecords.map((record) => (
                                                        <tr key={record.id}>
                                                            <td>{record.bookId}</td>
                                                            <td>{record.bookCopyId}</td>
                                                            <td>{new Date(record.borrowedAt).toLocaleString()}</td>
                                                            <td>{record.dueDate ? new Date(record.dueDate).toLocaleDateString() : 'N/A'}</td>
                                                            <td>
                                                                <span className={`badge ${
                                                                    record.status === 'RETURNED' ? 'bg-success' : 
                                                                    record.status === 'OVERDUE' ? 'bg-danger' : 'bg-warning'
                                                                }`}>
                                                                    {record.status}
                                                                </span>
                                                            </td>
                                                            <td>
                                                                {record.status === 'BORROWED' && (
                                                                    <button 
                                                                        className="btn btn-sm btn-outline-primary"
                                                                        onClick={() => returnBook(record.id)}
                                                                    >
                                                                        Return
                                                                    </button>
                                                                )}
                                                            </td>
                                                        </tr>
                                                    ))}
                                                </tbody>
                                            </table>
                                        </div>
                                    ) : (
                                        <div className="alert alert-info">No borrow records found.</div>
                                    )}
                                </div>
                                <div className="modal-footer">
                                    <button className="btn btn-secondary" onClick={() => setShowBorrowRecordsModal(false)}>Close</button>
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