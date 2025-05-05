import { createFileRoute } from '@tanstack/react-router'
import axios from 'axios';
import React, { useState } from 'react'
import "../css/userForms.css"
const url = import.meta.env.VITE_BACKEND_URL


export const Route = createFileRoute('/register')({
  component: Register,
})

function Register() {
  const [registerFormData, setRegisterFormData] = useState({
    username: "",
    email: "",
    password: "",
    confirmPassword: "",

  });

  function registerUser(event: React.FormEvent) {
    event.preventDefault()
    return axios
    .post(`${url}/user/register`, {
      email: registerFormData.email,
      password: registerFormData.password,
      username: registerFormData.username,
    })
  }

  function handleChange(event: React.ChangeEvent<HTMLInputElement>) {
    setRegisterFormData({
        ...registerFormData,
    [event.target.name]: event.target.value
    })
  }

  return (
    <>
      <form className='formContainer'>
            <label className='Label'>Username:</label>
            <input  name="username" value={registerFormData.username} onChange={handleChange} type="username" placeholder='Enter username'/>
            <label className='Label'>Email:</label>
            <input name="email" value={registerFormData.email} onChange={handleChange} type='email' placeholder='Enter Email'/>
            <label className='Label'>Password:</label>
            <input name="password" value={registerFormData.password}onChange={handleChange}  type='password' placeholder='Enter Password' />
            <label className='Label'>Confirm Password:</label>
            <input name="confirmPassword" value={registerFormData.confirmPassword} onChange={handleChange}  type='password' placeholder='Confirm Password' />
            <button onClick={registerUser} type='submit' >Register</button>
      </form>
    </>
  )
}
