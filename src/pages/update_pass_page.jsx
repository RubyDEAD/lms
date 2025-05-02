import React, { useState } from 'react';
import { supabase } from '../supabaseClient';
import 'bootstrap/dist/css/bootstrap.min.css';

function UpdatePasswordPage() {
  const [newPassword, setNewPassword] = useState('');
  const [message, setMessage] = useState('');
  const [errorMsg, setErrorMsg] = useState('');

  const handleUpdatePassword = async (e) => {
    e.preventDefault();

    const { error } = await supabase.auth.updateUser({
      password: newPassword
    });

    if (error) {
      setErrorMsg(error.message);
      setMessage('');
    } else {
      setMessage('Password updated successfully! You can now log in with your new password.');
      setErrorMsg('');
    }
  };

  return (
    <div className="container mt-5" style={{ maxWidth: '400px' }}>
      <h3>Reset Your Password</h3>
      <form onSubmit={handleUpdatePassword}>
        <div className="form-group mb-3">
          <label>New Password</label>
          <input
            type="password"
            className="form-control"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            required
          />
        </div>

        <button type="submit" className="btn btn-success w-100">Update Password</button>

        {message && <div className="alert alert-success mt-3">{message}</div>}
        {errorMsg && <div className="alert alert-danger mt-3">{errorMsg}</div>}
      </form>
    </div>
  );
}

export default UpdatePasswordPage;
