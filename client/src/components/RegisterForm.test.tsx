import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import RegisterForm from './RegisterForm'


describe("< RegisterForm />", () => {
    const handleSubmitMock = vi.fn()
    const handleChangeMock = vi.fn()
    const registerFormData = { username: '', password: '', email: '', confirmPassword: '' }
    const user = userEvent.setup()
    beforeEach(() => {
        render(
            <RegisterForm
                handleSubmit={handleSubmitMock}
                handleChange={handleChangeMock}
                registerFormData={registerFormData}
            />
        )
    })

    it("correctly updates form data on text input and submits", async () => {
        await user.type(screen.getByPlaceholderText('Enter username'), 'test')
        await user.type(screen.getByPlaceholderText('Enter email'), 'test@example.com')
        await user.type(screen.getByPlaceholderText('Enter password'), 'secure123')
        await user.type(screen.getByPlaceholderText('Confirm password'), 'secure123')

        expect(handleChangeMock).toHaveBeenCalledTimes('test'.length + "test@example.com".length + 'secure123'.length * 2)

        const submit = screen.getByRole("button")
        await user.click(submit)

        expect(handleSubmitMock).toHaveBeenCalledTimes(1)
    })
    it('renders all input fields and submit button', () => {
        expect(screen.getByPlaceholderText('Enter username')).toBeInTheDocument()
        expect(screen.getByPlaceholderText('Enter email')).toBeInTheDocument()
        expect(screen.getByPlaceholderText('Enter password')).toBeInTheDocument()
        expect(screen.getByPlaceholderText('Confirm password')).toBeInTheDocument()
        expect(screen.getByRole('button', { name: /register/i })).toBeInTheDocument()
    })


})