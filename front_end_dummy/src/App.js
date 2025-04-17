import React, { useState } from 'react';
import { ApolloProvider, useQuery, gql } from '@apollo/client';

// GraphQL Queries
const GET_BOOKS = gql`
  query getBooks {
    getBooks {
      id
      title
      author_name
      date_published
      description
    }
  }
`;

const GET_BOOK_COPIES = gql`
  query GetBookCopiesById($id: String!) {
    getBookCopiesById(id: $id) {
      id
      book_id
      book_status
    }
  }
`;

// Books Table Component
const BooksTable = ({ onBookClick }) => {
  const { loading, error, data } = useQuery(GET_BOOKS);

  if (loading) return <p>Loading books...</p>;
  if (error) return <p>Error loading books: {error.message}</p>;

  return (
    <table border="1" style={{ width: '100%', textAlign: 'left' }}>
      <thead>
        <tr>
          <th>Title</th>
          <th>Author</th>
          <th>Date Published</th>
          <th>Description</th>
        </tr>
      </thead>
      <tbody>
        {data.getBooks.map((book) => (
          <tr key={book.id} onClick={() => onBookClick(book.id)} style={{ cursor: 'pointer' }}>
            <td>{book.title}</td>
            <td>{book.author_name}</td>
            <td>{book.date_published}</td>
            <td>{book.description}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

// Book Copies Component
const BookCopies = ({ bookId }) => {
  const { loading, error, data } = useQuery(GET_BOOK_COPIES, {
    variables: { id: bookId },
    skip: !bookId, // Skip query if no bookId is selected
  });

  if (!bookId) return <p>Select a book to view its copies.</p>;
  if (loading) return <p>Loading book copies...</p>;
  if (error) return <p>Error loading book copies: {error.message}</p>;

  return (
    <table border="1" style={{ width: '100%', textAlign: 'left', marginTop: '20px' }}>
      <thead>
        <tr>
          <th>Copy ID</th>
          <th>Book ID</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        {data.getBookCopiesById.map((copy) => (
          <tr key={copy.id}>
            <td>{copy.id}</td>
            <td>{copy.book_id}</td>
            <td>{copy.book_status}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

// Main App Component
const App = () => {
  const [selectedBookId, setSelectedBookId] = useState(null);

  return (
    <ApolloProvider client={require('./apolloClient').default}>
      <div style={{ padding: '20px' }}>
        <h1>Books</h1>
        <BooksTable onBookClick={setSelectedBookId} />
        <h2>Book Copies</h2>
        <BookCopies bookId={selectedBookId} />
      </div>
    </ApolloProvider>
  );
};

export default App;