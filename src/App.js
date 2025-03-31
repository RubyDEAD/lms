import React from 'react';
import './App.css';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom'; // Use Routes instead of Switch
import Dashboard from './components/dashboard';
import Sidebar from './components/sidebar';
import Books from './components/book';
import BorrowedBooks from './components/BorrowedBooks';  // Import missing components
import AddBook from './components/AddBook';  // Import missing components
import Profile from './components/profile';  // Import missing components
import Topbar from './components/topbar';

function App() {
  return (
    <Router>
      <div className="App">
        <Topbar /> {/* Add Topbar here to make it visible across all routes */}
        
        <div style={{ display: 'flex' }}>
          <Sidebar /> {/* Sidebar stays fixed on the left */}

          <div style={{ marginLeft: '250px', padding: '20px', flex: 1 }}>
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/books" element={<Books />} />
              <Route path="/borrowed-books" element={<BorrowedBooks />} />
              <Route path="/add-book" element={<AddBook />} />
              <Route path="/profile" element={<Profile />} />
            </Routes>
          </div>
        </div>
      </div>
    </Router>
  );
}

export default App;
