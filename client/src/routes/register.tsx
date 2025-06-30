import { createFileRoute } from '@tanstack/react-router'
import React, { useState } from 'react'
import "../css/userForms.css"
import { postRegister } from '../utils/auth';
import SubmitError from '../components/SubmitError';


export const Route = createFileRoute('/register')({
        component: Register,
})

function Register() {
        const [error, setError] = React.useState<string | null>(null);
        const [registerFormData, setRegisterFormData] = useState({
                username: "",
                email: "",
                password: "",
                confirmPassword: "",

        });

        async function handleSubmit(event: React.FormEvent) {
                event.preventDefault()
                try {
                        await postRegister(registerFormData.email, registerFormData.password, registerFormData.username)
                } catch (err: any) {
                        setError(err.message || "Register failed.")
                }

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
                                <input name="username" value={registerFormData.username} onChange={handleChange} type="username" placeholder='Enter username' />
                                <label className='Label'>Email:</label>
                                <input name="email" value={registerFormData.email} onChange={handleChange} type='email' placeholder='Enter Email' />
                                <label className='Label'>Password:</label>
                                <input name="password" value={registerFormData.password} onChange={handleChange} type='password' placeholder='Enter Password' />
                                <label className='Label'>Confirm Password:</label>
                                <input name="confirmPassword" value={registerFormData.confirmPassword} onChange={handleChange} type='password' placeholder='Confirm Password' />
                                <button onClick={handleSubmit} type="submit" >Register</button>
                        </form>
                        {error && <SubmitError err={error}/>}
                </>
        )
}
