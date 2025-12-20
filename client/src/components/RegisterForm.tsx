import React from "react"
import styles from "./RegisterForm.module.css"

export default function RegisterForm({ handleSubmit, handleChange, registerFormData }: {
    handleSubmit: (e: React.FormEvent) => void
    handleChange: (e: React.ChangeEvent<HTMLInputElement>) => void
    registerFormData: { username: string, email: string, password: string, confirmPassword: string }
}) {
    const match =
        registerFormData.password === registerFormData.confirmPassword &&
        registerFormData.password.length >= 6
    return (
        <>
            <div style={{fontSize: "2em"}}>Password needs to be atleast 6 characters long. </div>
            <form className={styles.formContainer} onSubmit={(e) => void handleSubmit(e)}>
                <label className='Label'>Username:</label>
                <input name="username" value={registerFormData.username} onChange={handleChange} type="text" placeholder='Enter username' />
                <label className='Label'>Email:</label>
                <input name="email" value={registerFormData.email} onChange={handleChange} type='email' placeholder='Enter email' />
                <label className='Label'>Password:</label>
                <input name="password" value={registerFormData.password} onChange={handleChange} type='password' placeholder='Enter password' />
                <label className='Label'>Confirm Password:</label>
                <input name="confirmPassword" value={registerFormData.confirmPassword} onChange={handleChange} type='password' placeholder='Confirm password' />
                <button type="submit" name="register" disabled={!match}>Register</button>
                { !match && <span style={{fontSize: "1.5em"}}>Passwords do not match</span>}
            </form>
        </>
    )
}