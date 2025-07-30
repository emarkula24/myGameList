

export default function RegisterForm({ handleSubmit, handleChange, registerFormData }: {
    handleSubmit: (e: React.FormEvent) => void
    handleChange: (e: React.ChangeEvent<HTMLInputElement>) => void
    registerFormData: { username: string, email: string, password: string, confirmPassword: string }
}) {
    return (
        <>
            <form className='formContainer' onSubmit={(e) => void handleSubmit(e)}>
                <label className='Label'>Username:</label>
                <input name="username" value={registerFormData.username} onChange={handleChange} type="username" placeholder='Enter username' />
                <label className='Label'>Email:</label>
                <input name="email" value={registerFormData.email} onChange={handleChange} type='email' placeholder='Enter Email' />
                <label className='Label'>Password:</label>
                <input name="password" value={registerFormData.password} onChange={handleChange} type='password' placeholder='Enter Password' />
                <label className='Label'>Confirm Password:</label>
                <input name="confirmPassword" value={registerFormData.confirmPassword} onChange={handleChange} type='password' placeholder='Confirm Password' />
                <button type="submit" >Register</button>
            </form>
        </>
    )
}