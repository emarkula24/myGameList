import { createFileRoute, redirect, useRouter, useRouterState } from '@tanstack/react-router'
import React, { useState } from 'react'
import { z } from 'zod'
import SubmitError from '../components/SubmitError'
import { useAuth } from '../utils/auth'
import CommonDivider from '../components/CommonDivider'
import LoginForm from '../components/LoginForm'
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
                        <LoginForm onChange={onChange} onSubmit={onSubmit} loginFormData={loginFormData} />
                        {isLoggingIn && 'Loading...'}
                        {error && <SubmitError err={error} />}
                </div>

        )
}
