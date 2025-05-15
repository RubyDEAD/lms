import React, { useEffect, useState } from 'react';
import './App.css';
import { Route, Routes, useLocation, Navigate } from 'react-router-dom';
import Dashboard from './components/dashboard';
import Sidebar from './components/sidebar';
import Books from './components/book';
import BorrowedBooks from './components/BorrowedBooks';
import AddBook from './components/AddBook';
import Profile from './components/profile';
import Fines from "./components/fines";
import Topbar from './components/topbar';
import SignUpPage from './pages/signup_page';
import LoginPage from './pages/login_page';
import ForgotPasswordPage from './pages/forgot_pass_page';
import UpdatePasswordPage from './pages/update_pass_page';
import TestAdminPage from './pages/test_admin_page';
import { supabase } from './supabaseClient';

function App() {
  const location = useLocation();
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  const noDesignRoutes = ['/sign-up', '/login', '/forgot-password', '/update-password'];

  useEffect(() => {
    const getSession = async () => {
      const { data: { session } } = await supabase.auth.getSession();
      setUser(session?.user ?? null);
      setLoading(false);
    };

    getSession();

    const { data: listener } = supabase.auth.onAuthStateChange((_event, session) => {
      setUser(session?.user ?? null);
    });

    return () => {
      listener.subscription.unsubscribe();
    };
  }, []);

  const isLayoutVisible = !noDesignRoutes.includes(location.pathname);

  if (loading) return <div>Loading...</div>;

  return (
    <div className="App">
      {isLayoutVisible && user && <Topbar />}

      <div style={{ display: 'flex' }}>
        {isLayoutVisible && user && <Sidebar />}

        <div style={{ marginLeft: isLayoutVisible && user ? '250px' : 0, padding: '20px', flex: 1 }}>
          <Routes>
            {/* Public Routes */}
            <Route path="/login" element={<LoginPage />} />
            <Route path="/sign-up" element={<SignUpPage />} />
            <Route path="/forgot-password" element={<ForgotPasswordPage />} />
            <Route path="/update-password" element={<UpdatePasswordPage />} />

            {/* Root Redirect */}
            <Route
              path="/"
              element={
                user ? <Navigate to="/dashboard" replace /> : <Navigate to="/login" replace />
              }
            />

            {/* Protected Routes */}
            {user ? (
              <>
                <Route path="/dashboard" element={<Dashboard />} />
                <Route path="/books" element={<Books />} />
                <Route path="/borrowed-books" element={<BorrowedBooks />} />
                <Route path="/add-book" element={<AddBook />} />
                <Route path="/profile" element={<Profile />} />
                <Route path="/fines" element={<Fines />} />

                {/* Testing rani nako, kamo bahala unsaon pag protecteed route sa admin kay di nako kamao ana */}
                <Route path="/admin-test-page" element={<TestAdminPage />} />
              </>
            ) : (
              <Route path="*" element={<Navigate to="/login" replace />} />
            )}
          </Routes>
        </div>
      </div>
    </div>
  );
}

export default App;
