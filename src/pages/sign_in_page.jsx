import React from "react";
import 'bootstrap/dist/css/bootstrap.min.css';
import { useState } from 'react';

function SignInPage() {
  const [inputs, setInputs] = useState({
    first_name: '',
    last_name: '',
    phoneNumber: '',
    email: '',
    password: '',
    confirm_password: ''
  });

  const handleChange = (event) => {
    const name = event.target.name;
    const value = event.target.value;
    setInputs(values => ({...values, [name]: value}))
  }

  const handleSubmit = (event) => {
    event.preventDefault();
    alert(inputs);
  }

  return (
    <form onSubmit={handleSubmit}>
      <label>Enter your First Name:
        <input 
            type="text" 
            name="first_name" 
            value={inputs.first_name || ""} 
            onChange={handleChange}
        />
      </label>

      <label>Enter your Last Name:
        <input 
            type="text" 
            name="last_name" 
            value={inputs.last_name || ""} 
            onChange={handleChange}
        />
      </label>

      <label>Enter your Phone Number:
        <input 
            type="tel"
            name="phoneNumber"
            value={inputs.phoneNumber || ""}
            onChange={handleChange}
            pattern='^[0-9]{10,15}$'  // Optional: regex pattern for international phone formats
        />
      </label>

      <label>Enter your Email:
        <input 
            type="email" 
            name="email" 
            value={inputs.email || ""} 
            onChange={handleChange}
        />
      </label>

      <label>Enter your Password:
        <input 
            type="password" 
            name="password" 
            value={inputs.password || ""} 
            onChange={handleChange}
        />
      </label>

      <label>Confirm your Password
        <input 
            type="password" 
            name="confirm_password" 
            value={inputs.confirm_password || ""} 
            onChange={handleChange}
        />
      </label>
      
        <input type="submit" />
    </form>
  )
}


export default SignInPage;