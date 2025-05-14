import { createFileRoute } from '@tanstack/react-router'
import axios from 'axios'
import { useState } from 'react'
import { useMutation } from '@tanstack/react-query'

export const Route = createFileRoute('/login')({
        component: Login,
})

const url = import.meta.env.VITE_BACKEND_URL

function Login() {
        const [loginFormData, setLoginFormData] = useState({
                username: "",
                password: "",
        })
        const loginMutation = useMutation({
                mutationFn: async (loginData: { username: string; password: string }) => {
                        const response = await axios.post(`${url}/login`, loginData);
                        return response.data;
                },
                onSuccess: (data) => {
                        console.log(data)
                        localStorage.setItem("token", data.accessToken)
                        // AXIOS INTERCEPTOR FOR 401 CAN BE USED TO REFRESH
                },
                onError: (error) => {

                        console.log(error)
                },
        });

        const handleSubmit = (event: React.FormEvent) => {
                event.preventDefault();
                loginMutation.mutate({
                        username: loginFormData.username,
                        password: loginFormData.password,
                });
        };
        function handleChange(event: React.ChangeEvent<HTMLInputElement>) {
                setLoginFormData({
                        ...loginFormData,
                        [event.target.name]: event.target.value
                })
        }

        return (
                <div>
                        <form className="formContainer">
                                <label className="Label">Username:</label>
                                <input name="username" value={loginFormData.username} onChange={handleChange} type="username" placeholder="Enter username" />
                                <label className="Label">Password:</label>
                                <input name="password" value={loginFormData.password} onChange={handleChange} type="password" placeholder="Enter password" />
                                <button onClick={handleSubmit} type="submit">Login</button>
                        </form>
                        {loginMutation.isPending ? (
                                "logging in..."
                        ) : (
                                <>
                                        {loginMutation.isError ? (
                                                <div>failed to log in: {loginMutation.error.message}</div>
                                        ) : null}

                                        {loginMutation.isSuccess ? <div>login succeeded</div> : null}
                                        </>
                                        )}
                                </div>
                        )
}
