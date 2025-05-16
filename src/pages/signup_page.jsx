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

  const createPatron = async (mutation) => {
    try {
      const response = await axios.post(API_URL, {query: mutation})
      console.log(response)

      const { data, error } = await supabase.auth.signInWithPassword({
        email: inputs.email,
        password: inputs.password
      });
      if (error) throw error;
  
      localStorage.setItem('user', JSON.stringify(data.user));

      setTimeout(() => {
        window.location.href = '/books';
      }, 1500);

    } catch (err){
      console.error("Error adding user: ", err);
    }
  }

  const handleChange = (event) => {
    const name = event.target.name;
    const value = event.target.value;
    setInputs(values => ({...values, [name]: value}))
  }

  const handleSubmit = (event) => {
    event.preventDefault();
    console.log(inputs);

    // Check if the passwords are the same
    if (inputs.password !== inputs.confirm_password) {
      alert("Passwords do not match");
      setInputs({
        first_name: '',
        last_name: '',
        phoneNumber: '',
        email: '',
        password: '',
        confirm_password: ''
      });

      return;
    }

    //Forward request to API-GATEWAY
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
    `

    createPatron(mutation);


  }

  return (
    <div className="signup-container">
      {/* LMS Section (Left) */}
      <div className="lms-banner">
        <h1>Welcome to LMS</h1>
        <p>Already have an Account? Sign In now.</p>
        <button onClick={() => window.location.href = '/login_page'}>Sign In</button>
      </div>

      {/* Signup Section (Right) */}
      <div className="signup-card-container">
        <div className="signup-card card">
          <h2>Create an Account</h2>
          {/* {error && <div className="alert alert-danger">{error}</div>} */}

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
              <input type="tel" name="phoneNumber" className="form-control" value={inputs.phoneNumber} onChange={handleChange} required />
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
    </div>
  );
}

export default SignUpPage;
