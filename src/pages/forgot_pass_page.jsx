import React, { useState } from 'react';
import { supabase } from '../supabaseClient';
import 'bootstrap/dist/css/bootstrap.min.css';

function ForgotPasswordPage() {
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState('');
  const [errorMsg, setErrorMsg] = useState('');

  const handlePasswordReset = async (e) => {
    e.preventDefault();

    const { error } = await supabase.auth.resetPasswordForEmail(email, {
      redirectTo: 'http://localhost:3000/update-password', // Update this to your reset handler route
    });

    if (error) {
      setErrorMsg(error.message);
      setMessage('');
    } else {
      setMessage('Check your email for the password reset link.');
      setErrorMsg('');
    }
  };

  return (
    <div className="container mt-5" style={{ maxWidth: '400px' }}>
      <h3>Forgot Password</h3>
      <form onSubmit={handlePasswordReset}>
        <div className="form-group mb-3">
          <label>Email address</label>
          <input
            type="email"
            className="form-control"
            placeholder="Enter your email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>

        <button type="submit" className="btn btn-primary w-100">Send Reset Link</button>

        {message && <div className="alert alert-success mt-3">{message}</div>}
        {errorMsg && <div className="alert alert-danger mt-3">{errorMsg}</div>}
      </form>
    </div>
  );
}

export default ForgotPasswordPage;
