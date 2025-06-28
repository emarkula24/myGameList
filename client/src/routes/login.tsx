import { createFileRoute, useRouter } from '@tanstack/react-router'
import React, { useState } from 'react'
import { z } from 'zod'

export const Route = createFileRoute('/login')({
        validateSearch: z.object({
                redirect: z.string().optional(),
        }),
        component: LoginComponent,
})


function LoginComponent() {
        const router = useRouter()
        const { auth, status } = Route.useRouteContext({
                select: ({ auth }) => ({ auth, status: auth.status })
        })
        const search = Route.useSearch()
        const [loginFormData, setLoginFormData] = useState({
                username: "",
                password: "",
        })

        const onSubmit = (event: React.FormEvent) => {
                event.preventDefault();
                auth.login(loginFormData.username, loginFormData.password)
                router.invalidate()
        };

        React.useLayoutEffect(() => {
                if (status === "loggedIn" && search.redirect) {
                        router.history.push(search.redirect)
                }
        }, [status, search.redirect])

        function onChange(event: React.ChangeEvent<HTMLInputElement>) {
                setLoginFormData({
                        ...loginFormData,
                        [event.target.name]: event.target.value
                })
        }

        return (
                <div>
                        <form className="formContainer">
                                <label className="Label">Username:</label>
                                <input name="username" value={loginFormData.username} onChange={onChange} type="text" placeholder="Enter username" />
                                <label className="Label">Password:</label>
                                <input name="password" value={loginFormData.password} onChange={onChange} type="password" placeholder="Enter password" />
                                <button onClick={onSubmit} type="submit">Login</button>
                        </form>
                        {auth.status === "loggedIn" && <p>logged in!</p>}
                </div>
                
        )
}
