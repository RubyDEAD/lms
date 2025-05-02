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
import SignUpPage from './pages/signup_page';
import LoginPage from './pages/login_page';
import ForgotPasswordPage from './pages/forgot_pass_page';
import UpdatePasswordPage from './pages/update_pass_page';

function App() {
  const location = useLocation();
  const noDesignRoutes = ['/sign-up', '/login', '/forgot-password', '/update-password'];

  const isLayoutVisible = !noDesignRoutes.includes(location.pathname);

  return (
    <div className="App">
      {isLayoutVisible && <Topbar />}

      <div style={{ display: 'flex' }}>
        {isLayoutVisible && <Sidebar />}

        <div style={{ marginLeft: isLayoutVisible ? '250px' : 0, padding: '20px', flex: 1 }}>
          <Routes>
            {/* Auth & Onboarding */}
            <Route path="/login" element={<LoginPage />} />
            <Route path="/sign-up" element={<SignUpPage />} />
            <Route path="/forgot-password" element={<ForgotPasswordPage />} />
            <Route path="/update-password" element={<UpdatePasswordPage />} />

            {/* App Pages */}
            <Route path="/" element={<Dashboard />} />
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/books" element={<Books />} />
            <Route path="/borrowed-books" element={<BorrowedBooks />} />
            <Route path="/add-book" element={<AddBook />} />
            <Route path="/profile" element={<Profile />} />
          </Routes>
        </div>
      </div>
    </div>
  );
}

export default App;
