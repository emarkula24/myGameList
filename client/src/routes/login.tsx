import { createFileRoute } from '@tanstack/react-router'
import axios from 'axios'
import { useState } from 'react'

export const Route = createFileRoute('/login')({
        component: Login,
})

const url = import.meta.env.VITE_BACKEND_URL

function Login() {
        const [loginFormData, setLoginFormData] = useState({
                username: "",
                password: "",
        })

        function loginUser(event: React.FormEvent) {
                event.preventDefault()
                return axios
                .post(`${url}/user/login`, {
                        username: loginFormData.username,
                        password: loginFormData.password
                })
        }
        function handleChange(event: React.ChangeEvent<HTMLInputElement>) {
                setLoginFormData({
                        ...loginFormData,
                        [event.target.name]: event.target.value
                })
        }

        return (
                <>
                        <form className="formContainer">
                                <label className="Label">Username:</label>
                                <input name="username" value={loginFormData.username} onChange={handleChange} type="username" placeholder="Enter username" />
                                <label className="Label">Password:</label>
                                <input name="password" value={loginFormData.password} onChange={handleChange} type="password" placeholder="Enter password" />
                                <button onClick={loginUser} type="submit">Login</button>
                        </form>
                </>
        )
}
