
export default function LoginForm({onSubmit, onChange, loginFormData}: {
    onSubmit: (e: React.FormEvent) => void
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void
    loginFormData: { username: string; password: string }
}) {
    return (
        <>
            <form className="formContainer" onSubmit={(e) => void onSubmit(e)}>
                <label className="Label">Username:</label>
                <input name="username" value={loginFormData.username} onChange={(e) => onChange(e)} type="text" placeholder="Enter username" required />
                <label className="Label">Password:</label>
                <input name="password" value={loginFormData.password} onChange={(e) => onChange(e)} type="password" placeholder="Enter password" required />
                <button type="submit">Login</button>
            </form>
        </>
    )
}