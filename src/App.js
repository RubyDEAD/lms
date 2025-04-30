import React from 'react';
import './App.css';
import { Route, Routes, useLocation } from 'react-router-dom';
import Dashboard from './components/dashboard';
import Sidebar from './components/sidebar';
import Books from './components/book';
import BorrowedBooks from './components/BorrowedBooks';
import AddBook from './components/AddBook';
import Profile from './components/profile';
import Topbar from './components/topbar';
import SignInPage from './pages/sign_in_page';

function App() {
  const location = useLocation();
  const noDesignRoutes = ['/sign-in'];

  return (
    <div className="App">
      {!noDesignRoutes.includes(location.pathname) && <Topbar />}

      <div style={{ display: 'flex' }}>
        {!noDesignRoutes.includes(location.pathname) && <Sidebar />}

        {!noDesignRoutes.includes(location.pathname) ? (
          <div style={{ marginLeft: '250px', padding: '20px', flex: 1 }}>
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/sign-in" element={<SignInPage />} />
              <Route path="/books" element={<Books />} />
              <Route path="/borrowed-books" element={<BorrowedBooks />} />
              <Route path="/add-book" element={<AddBook />} />
              <Route path="/profile" element={<Profile />} />
            </Routes>
          </div>
        ) : (
          <Routes>
            <Route path="/sign-in" element={<SignInPage />} />
          </Routes>
        )}
      </div>
    </div>
  );
}

export default App;
