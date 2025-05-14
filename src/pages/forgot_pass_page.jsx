import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { supabase } from '../supabaseClient';
import 'bootstrap/dist/css/bootstrap.min.css';
import '../App.css'; // Ensure styles are applied

function ForgotPasswordPage() {
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState('');
  const [errorMsg, setErrorMsg] = useState('');
  const navigate = useNavigate();

  const handlePasswordReset = async (e) => {
    e.preventDefault();

    const { error } = await supabase.auth.resetPasswordForEmail(email, {
      redirectTo: 'http://localhost:3000/update-password', // Change as needed
    });

    if (error) {
      setErrorMsg(error.message);
      setMessage('');
    } else {
      setMessage('Check your email for the password reset link.');
      setErrorMsg('');
    }
  };

  const handleCancel = () => {
    navigate('/login');
  };

  return (
    <div className="login-container">
      {/* LMS Banner (Left Side) */}
      <div className="lms-banner">
        <h1>Welcome to LMS</h1>
        <p>Empowering minds, one book at a time.</p>
        <button onClick={() => navigate('/login')}>Back to Login</button>
      </div>

      {/* Forgot Password Card (Right Side) */}
      <div className="login-card-container">
        <div className="login-card">
          <h2>Forgot Password</h2>
          <form onSubmit={handlePasswordReset}>
            <div className="form-group mb-3">
              <label className="form-label">Email address</label>
              <input
                type="email"
                className="form-control"
                placeholder="Enter your email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>

            <button type="submit" className="btn btn-primary w-100 mb-2">
              Send Reset Link
            </button>
          
            {message && <div className="alert alert-success mt-3">{message}</div>}
            {errorMsg && <div className="alert alert-danger mt-3">{errorMsg}</div>}
          </form>
        </div>
      </div>
    </div>
  );
}

export default ForgotPasswordPage;
