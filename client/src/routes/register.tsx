import { createFileRoute, useRouter } from '@tanstack/react-router'
import React, { useState } from 'react'
import { postRegister } from '../utils/auth';
import SubmitError from '../components/SubmitError';
import CommonDivider from '../components/CommonDivider';
import RegisterForm from '../components/RegisterForm';

export const Route = createFileRoute('/register')({
        component: Register,
})

function Register() {
        const router = useRouter()
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
                        await router.navigate({ to: "/login" })
                } catch (err: unknown) {
                        if (err instanceof Error) {
                                setError(err.message || "Register failed.")
                        }

                }

        }

        function handleChange(event: React.ChangeEvent<HTMLInputElement>) {
                setRegisterFormData({
                        ...registerFormData,
                        [event.target.name]: event.target.value
                })
        }

        return (
                <div className="routeContainer">
                        <CommonDivider routeName={"Register"} />
                        <RegisterForm handleSubmit={() => void handleSubmit} handleChange={() => void handleChange} registerFormData={registerFormData} />
                        {error && <SubmitError err={error} />}
                </div>
        )
}
