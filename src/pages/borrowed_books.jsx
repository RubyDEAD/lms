import React, { useEffect, useState } from "react";
import { gql, useQuery } from "@apollo/client";

const PATRON_BORROW_HISTORY = gql`
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
      previousDueDate
    }
  }
`;

const BorrowedBooks = () => {
  const patronId = "your-patron-id"; // Replace this with actual ID (from auth/session)
  const { loading, error, data } = useQuery(PATRON_BORROW_HISTORY, {
    variables: { patronId },
  });

  if (loading) return <p>Loading borrowed books...</p>;
  if (error) return <p>Error fetching borrow history: {error.message}</p>;

  const borrowedBooks = data.patronBorrowHistory;

  return (
    <div className="borrowed-books-page">
      <h1>Borrowed Books</h1>
      <table className="borrowed-books-table">
        <thead>
          <tr>
            <th>Book ID</th>
            <th>Borrowed At</th>
            <th>Due Date</th>
            <th>Returned At</th>
            <th>Status</th>
            <th>Renewals</th>
          </tr>
        </thead>
        <tbody>
          {borrowedBooks.length === 0 ? (
            <tr>
              <td colSpan="6" style={{ textAlign: "center" }}>
                No borrowed books found.
              </td>
            </tr>
          ) : (
            borrowedBooks.map((record) => (
              <tr key={record.id}>
                <td>{record.bookId}</td>
                <td>{new Date(record.borrowedAt).toLocaleDateString()}</td>
                <td>{new Date(record.dueDate).toLocaleDateString()}</td>
                <td>
                  {record.returnedAt
                    ? new Date(record.returnedAt).toLocaleDateString()
                    : "Not returned"}
                </td>
                <td>{record.status}</td>
                <td>{record.renewalCount}</td>
              </tr>
            ))
          )}
        </tbody>
      </table>

      <style jsx>{`
        .borrowed-books-page {
          padding: 20px;
        }
        .borrowed-books-table {
          width: 100%;
          border-collapse: collapse;
          margin-top: 20px;
        }
        .borrowed-books-table th,
        .borrowed-books-table td {
          border: 1px solid #ccc;
          padding: 8px;
        }
        .borrowed-books-table th {
          background-color: #f5f5f5;
        }
      `}</style>
    </div>
  );
};

export default BorrowedBooks;
