import React, { useEffect, useState } from "react";
import axios from "axios";

function Books() {
    const [books, setBooks] = useState([]);
    const [bookDetails, setBookDetails] = useState(null);
    const [bookCopies, setBookCopies] = useState([]);
    const [searchResults, setSearchResults] = useState([]);
    const [availableCopy, setAvailableCopy] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [newBook, setNewBook] = useState({
        title: "",
        authorName: "",
        datePublished: "",
        description: "",
    });

    const API_URL = "http://localhost:8081/query";

    // Fetch all books
    useEffect(() => {
        const fetchBooks = async () => {
            try {
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
                });
                setBooks(response.data.data.getBooks);
                setLoading(false);
            } catch (err) {
                console.error("Error fetching books:", err);
                setError("Failed to fetch books. Please try again later.");
                setLoading(false);
            }
        };

        fetchBooks();
    }, []);

    // Add a new book
    const addBook = async () => {
        try {
            const response = await axios.post(API_URL, {
                query: `
                    mutation {
                        addBook(
                            title: "${newBook.title}",
                            authorName: "${newBook.authorName}",
                            datePublished: "${newBook.datePublished}",
                            description: "${newBook.description}"
                        ) {
                            id
                            title
                            author_name
                            date_published
                            description
                        }
                    }
                `,
            });

            const addedBook = response.data.data.addBook;

            if (addedBook) {
                setBooks((prevBooks) => [...prevBooks, addedBook]); // Add the new book to the list
                setNewBook({ title: "", authorName: "", datePublished: "", description: "" }); // Reset form
                setError(null); // Clear any previous errors
            }
        } catch (err) {
            console.error("Error adding book:", err);
            setError("Failed to add book. Please try again later.");
        }
    };

    // Fetch book details by ID
    const fetchBookById = async (id) => {
        try {
            const response = await axios.post(API_URL, {
                query: `
                    query {
                        getBookById(id: "${id}") {
                            id
                            title
                            author_name
                            date_published
                            description
                        }
                    }
                `,
            });
            setBookDetails(response.data.data.getBookById);
            setAvailableCopy(null); // Clear availability state
            setBookCopies([]); // Clear book copies state
        } catch (err) {
            console.error("Error fetching book by ID:", err);
            setError("Failed to fetch book details. Please try again later.");
        }
    };

    // Fetch book copies by book ID
    const fetchBookCopiesById = async (id) => {
        try {
            const response = await axios.post(API_URL, {
                query: `
                    query {
                        getBookCopiesById(id: "${id}") {
                            id
                            book_id
                            book_status
                        }
                    }
                `,
            });
            setBookCopies(response.data.data.getBookCopiesById);
            setBookDetails(null); // Clear book details state
            setAvailableCopy(null); // Clear availability state
        } catch (err) {
            console.error("Error fetching book copies:", err);
            setError("Failed to fetch book copies. Please try again later.");
        }
    };

    // Fetch available book copy by book ID
    const fetchAvailableBookCopyById = async (id) => {
        try {
            const response = await axios.post(API_URL, {
                query: `
                    query {
                        getAvailbleBookCopyByID(id: "${id}") {
                            id
                            book_id
                            book_status
                        }
                    }
                `,
            });

            const availableCopy = response.data.data.getAvailbleBookCopyByID;

            if (availableCopy) {
                setAvailableCopy(availableCopy);
                setBookDetails(null); // Clear book details state
                setBookCopies([]); // Clear book copies state
            } else {
                setAvailableCopy(null); // No available copy
            }
        } catch (err) {
            console.error("Error fetching available book copy:", err);
            setAvailableCopy(null); // Reset available copy state
        }
    };

    if (loading) {
        return <div className="container">Loading...</div>;
    }

    if (error && error !== "No available copy for this book.") {
        return <div className="container text-danger">{error}</div>;
    }

    return (
        <div className="body">
            <div className="container">
                <h1>Books</h1>

                {/* Add Book Form */}
                <div className="mt-4">
                    <h2>Add a New Book</h2>
                    <form
                        onSubmit={(e) => {
                            e.preventDefault();
                            addBook();
                        }}
                    >
                        <div className="mb-3">
                            <label htmlFor="title" className="form-label">
                                Title
                            </label>
                            <input
                                type="text"
                                className="form-control"
                                id="title"
                                value={newBook.title}
                                onChange={(e) => setNewBook({ ...newBook, title: e.target.value })}
                                required
                            />
                        </div>
                        <div className="mb-3">
                            <label htmlFor="authorName" className="form-label">
                                Author Name
                            </label>
                            <input
                                type="text"
                                className="form-control"
                                id="authorName"
                                value={newBook.authorName}
                                onChange={(e) =>
                                    setNewBook({ ...newBook, authorName: e.target.value })
                                }
                                required
                            />
                        </div>
                        <div className="mb-3">
                            <label htmlFor="datePublished" className="form-label">
                                Date Published
                            </label>
                            <input
                                type="date"
                                className="form-control"
                                id="datePublished"
                                value={newBook.datePublished}
                                onChange={(e) =>
                                    setNewBook({ ...newBook, datePublished: e.target.value })
                                }
                                required
                            />
                        </div>
                        <div className="mb-3">
                            <label htmlFor="description" className="form-label">
                                Description
                            </label>
                            <textarea
                                className="form-control"
                                id="description"
                                rows="3"
                                value={newBook.description}
                                onChange={(e) =>
                                    setNewBook({ ...newBook, description: e.target.value })
                                }
                                required
                            ></textarea>
                        </div>
                        <button type="submit" className="btn btn-primary">
                            Add Book
                        </button>
                    </form>
                </div>

                <div className="row mt-4">
                    {books.map((book) => (
                        <div className="col-md-4" key={book.id}>
                            <div className="card mb-4 shadow-sm">
                                <div className="card-body">
                                    <h5 className="card-title">{book.title}</h5>
                                    <p className="card-text">
                                        <strong>Author:</strong> {book.author_name}
                                    </p>
                                    <p className="card-text">
                                        <strong>Published:</strong> {book.date_published}
                                    </p>
                                    <p className="card-text">{book.description}</p>
                                    <button
                                        className="btn btn-primary"
                                        onClick={() => fetchBookById(book.id)}
                                    >
                                        View Details
                                    </button>
                                    <button
                                        className="btn btn-secondary"
                                        onClick={() => fetchBookCopiesById(book.id)}
                                    >
                                        View Copies
                                    </button>
                                    <button
                                        className="btn btn-success"
                                        onClick={() => fetchAvailableBookCopyById(book.id)}
                                    >
                                        Check Availability
                                    </button>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>

                {/* Book Details */}
                {bookDetails && (
                    <div className="mt-4">
                        <h2>Book Details</h2>
                        <p><strong>Title:</strong> {bookDetails.title}</p>
                        <p><strong>Author:</strong> {bookDetails.author_name}</p>
                        <p><strong>Published:</strong> {bookDetails.date_published}</p>
                        <p><strong>Description:</strong> {bookDetails.description}</p>
                    </div>
                )}

                {/* Book Copies */}
                {bookCopies.length > 0 && (
                    <div className="mt-4">
                        <h2>Book Copies</h2>
                        <ul>
                            {bookCopies.map((copy) => (
                                <li key={copy.id}>
                                    Copy ID: {copy.id}, Status: {copy.book_status}
                                </li>
                            ))}
                        </ul>
                    </div>
                )}

                {/* Available Copy */}
                {availableCopy ? (
                    <div className="mt-4">
                        <h2>Available Copy</h2>
                        <p>Copy ID: {availableCopy.id}, Status: {availableCopy.book_status}</p>
                    </div>
                ) : (
                    <div className="mt-4 text-warning">
                        <h2>Available Copy</h2>
                        <p>No available copy for this book.</p>
                    </div>
                )}
            </div>
        </div>
    );
}

export default Books;