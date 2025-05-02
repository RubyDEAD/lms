import React, { useState } from "react";
import 'bootstrap/dist/css/bootstrap.min.css';
import { supabase } from "../supabaseClient"; // make sure you have this file set up

function SignUpPage() {
  const [inputs, setInputs] = useState({
    first_name: '',
    last_name: '',
    phoneNumber: '',
    email: '',
    password: '',
    confirm_password: ''
  });

  const [errorMsg, setErrorMsg] = useState('');
  const [successMsg, setSuccessMsg] = useState('');

  const handleChange = (event) => {
    const { name, value } = event.target;
    setInputs(values => ({ ...values, [name]: value }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (inputs.password !== inputs.confirm_password) {
      setErrorMsg("Passwords do not match");
      return;
    }

    // Sign up with Supabase Auth
    const { data, error } = await supabase.auth.signUp({
      email: inputs.email,
      password: inputs.password,
      options: {
        data: {
          first_name: inputs.first_name,
          last_name: inputs.last_name,
          phone_number: inputs.phoneNumber
        }
      }
    });

    if (error) {
      console.error("Supabase signup error:", error);
      setErrorMsg(error.message);
      setSuccessMsg('');
    } else {
      console.log("User signed up:", data);
      setSuccessMsg("Signup successful! Check your email for confirmation.");
      setErrorMsg('');
      setInputs({
        first_name: '',
        last_name: '',
        phoneNumber: '',
        email: '',
        password: '',
        confirm_password: ''
      });
    }
  };

  return (
    <form onSubmit={handleSubmit} className="p-4">
      <h3>Sign Up</h3>

      <label>First Name:
        <input
          type="text"
          name="first_name"
          value={inputs.first_name}
          onChange={handleChange}
          className="form-control mb-2"
          required
        />
      </label>

      <label>Last Name:
        <input
          type="text"
          name="last_name"
          value={inputs.last_name}
          onChange={handleChange}
          className="form-control mb-2"
          required
        />
      </label>

      <label>Phone Number:
        <input
          type="tel"
          name="phoneNumber"
          value={inputs.phoneNumber}
          onChange={handleChange}
          className="form-control mb-2"
          pattern="^[0-9]{10,15}$"
          required
        />
      </label>

      <label>Email:
        <input
          type="email"
          name="email"
          value={inputs.email}
          onChange={handleChange}
          className="form-control mb-2"
          required
        />
      </label>

      <label>Password:
        <input
          type="password"
          name="password"
          value={inputs.password}
          onChange={handleChange}
          className="form-control mb-2"
          required
        />
      </label>

      <label>Confirm Password:
        <input
          type="password"
          name="confirm_password"
          value={inputs.confirm_password}
          onChange={handleChange}
          className="form-control mb-3"
          required
        />
      </label>

      <button type="submit" className="btn btn-primary">Sign Up</button>

      {errorMsg && <div className="alert alert-danger mt-3">{errorMsg}</div>}
      {successMsg && <div className="alert alert-success mt-3">{successMsg}</div>}
    </form>
  );
}

export default SignUpPage;
