import React from "react";
import 'bootstrap/dist/css/bootstrap.min.css';
import { useState } from 'react';
import axios from "axios";

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

      console.log(response);
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
    <form onSubmit={handleSubmit}>
      <label>Enter your First Name:
        <input 
            type="text" 
            name="first_name" 
            value={inputs.first_name} 
            onChange={handleChange}
            required={true}
        />
      </label>

      <label>Enter your Last Name:
        <input 
            type="text" 
            name="last_name" 
            value={inputs.last_name} 
            onChange={handleChange}
            required={true}
        />
      </label>

      <label>Enter your Phone Number:
        <input 
            type="tel"
            name="phoneNumber"
            value={inputs.phoneNumber}
            onChange={handleChange}
            pattern='^[0-9]{10,15}$'  // Default regex format that is used in patrondb
            required={true}
        />
      </label>

      <label>Enter your Email:
        <input 
            type="email" 
            name="email" 
            value={inputs.email} 
            onChange={handleChange}
            required={true}
        />
      </label>

      <label>Enter your Password: 
        <input 
            type="password" 
            name="password" 
            value={inputs.password} 
            onChange={handleChange}
            required={true}
        />
      </label>

      <label>Confirm your Password
        <input 
            type="password" 
            name="confirm_password" 
            value={inputs.confirm_password} 
            onChange={handleChange}
            required={true}
        />
      </label>
      
        <input type="submit" />
    </form>
  )
}


export default SignUpPage;