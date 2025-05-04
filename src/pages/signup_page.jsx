import React, { useState } from "react";
import 'bootstrap/dist/css/bootstrap.min.css';
import axios from "axios";
import { supabase } from '../supabaseClient';
import '../App.css';

function SignUpPage() {
  const API_URL = "http://localhost:8081/query";

  const [inputs, setInputs] = useState({
    first_name: '',
    last_name: '',
    phoneNumber: '',
    email: '',
    password: '',
    confirm_password: ''
  });

  const [error, setError] = useState(null);

  const handleChange = (event) => {
    const { name, value } = event.target;
    setInputs(values => ({ ...values, [name]: value }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (inputs.password !== inputs.confirm_password) {
      setError("Passwords do not match");
      return;
    }

    try {
      const mutation = `
        mutation {
          createPatron(
            first_name: "${inputs.first_name}"
            last_name: "${inputs.last_name}"
            phone_number: "${inputs.phoneNumber}"
            email: "${inputs.email}"
            password: "${inputs.password}"
          ) {
            first_name
            last_name
            phone_number
          }
        }
      `;

      await axios.post(API_URL, { query: mutation });

      const { data, error } = await supabase.auth.signInWithPassword({
        email: inputs.email,
        password: inputs.password
      });

      if (error) throw error;

      localStorage.setItem('user', JSON.stringify(data.user));
      window.location.href = '/dashboard';
    } catch (err) {
      console.error("Signup error:", err);
      setError("Failed to create account. Please try again.");
    }
  };

  return (
    <div className="d-flex justify-content-center align-items-center vh-100 bg-light">
      <div className="card shadow-lg p-4" style={{ width: '100%', maxWidth: '500px' }}>
        <h2 className="text-center mb-4">Create an Account</h2>

        {error && <div className="alert alert-danger">{error}</div>}

        <form onSubmit={handleSubmit}>
          <div className="mb-3">
            <label className="form-label">First Name</label>
            <input type="text" name="first_name" className="form-control" value={inputs.first_name} onChange={handleChange} required />
          </div>

          <div className="mb-3">
            <label className="form-label">Last Name</label>
            <input type="text" name="last_name" className="form-control" value={inputs.last_name} onChange={handleChange} required />
          </div>

          <div className="mb-3">
            <label className="form-label">Phone Number</label>
            <input type="tel" name="phoneNumber" className="form-control" value={inputs.phoneNumber} onChange={handleChange} pattern="^[0-9]{10,15}$" required />
          </div>

          <div className="mb-3">
            <label className="form-label">Email</label>
            <input type="email" name="email" className="form-control" value={inputs.email} onChange={handleChange} required />
          </div>

          <div className="mb-3">
            <label className="form-label">Password</label>
            <input type="password" name="password" className="form-control" value={inputs.password} onChange={handleChange} required />
          </div>

          <div className="mb-3">
            <label className="form-label">Confirm Password</label>
            <input type="password" name="confirm_password" className="form-control" value={inputs.confirm_password} onChange={handleChange} required />
          </div>

          <button type="submit" className="btn btn-primary w-100">Sign Up</button>
        </form>
      </div>
    </div>
  );
}

export default SignUpPage;
