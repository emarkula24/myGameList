import { createFileRoute, redirect, useRouter, useRouterState } from '@tanstack/react-router'
import React, { useState } from 'react'
import { z } from 'zod'
import SubmitError from '../components/SubmitError'
import { useAuth } from '../utils/auth'
import "../css/userForms.css"
import CommonDivider from '../components/CommonDivider'
const fallback = '/' as const

export const Route = createFileRoute('/login')({
        validateSearch: z.object({
                redirect: z.string().optional(),
        }),
        beforeLoad: ({ context, search }) => {
                if (context.auth.isAuthenticated) {
                        throw redirect({ to: search.redirect || fallback })
                }
        },
        component: LoginComponent,

})


function LoginComponent() {
        const [isSubmitting, setIsSubmitting] = React.useState(false)
        const isLoading = useRouterState({ select: (s) => s.isLoading })
        const auth = useAuth()
        const router = useRouter()
        const [error, setError] = React.useState<string | null>(null);
        const search = Route.useSearch()
        const [loginFormData, setLoginFormData] = useState({
                username: "",
                password: "",
        })

        const onSubmit = async (event: React.FormEvent) => {
                setIsSubmitting(true)
                event.preventDefault();
                try {
                        await auth.login(loginFormData.username, loginFormData.password)
                        await router.invalidate()
                        await router.navigate({ to: search.redirect || fallback })
                } catch (err: unknown) {
                        if (err instanceof Error) {
                                setError(err.message || 'Login failed');
                        }
                        
                } finally {
                        setIsSubmitting(false)
                }
        };

        function onChange(event: React.ChangeEvent<HTMLInputElement>) {
                setLoginFormData({
                        ...loginFormData,
                        [event.target.name]: event.target.value
                })
        }
        const isLoggingIn = isLoading || isSubmitting
        return (
                <div className="routeContainer">
                        <CommonDivider routeName="Login" />
                        <form className="formContainer" onSubmit={(e) =>  void onSubmit(e)}>
                                <label className="Label">Username:</label>
                                <input name="username" value={loginFormData.username} onChange={onChange} type="text" placeholder="Enter username" required />
                                <label className="Label">Password:</label>
                                <input name="password" value={loginFormData.password} onChange={onChange} type="password" placeholder="Enter password" required />
                                <button type="submit">Login</button>
                        </form>
                        {isLoggingIn ? 'Loading...' : 'Login'}
                        {error && <SubmitError err={error} />}
                </div>

        )
}
