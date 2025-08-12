import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import LoginForm from './LoginForm'


describe("< LoginForm />", () => {
    const onSubmitMock = vi.fn()
    const onChangeMock = vi.fn()
    const loginFormData = { username: 'test', password: 'secure' }
    const user = userEvent.setup()

    it("correctly updates form data on text input and submits", async () => {
        render(<LoginForm onSubmit={onSubmitMock} onChange={onChangeMock} loginFormData={loginFormData} />)
        const usernameInput = screen.getByPlaceholderText('Enter username') 
        const passwordInput = screen.getByPlaceholderText('Enter password')

        await user.type(usernameInput, 'test')
        await user.type(passwordInput, 'secure123')
        expect(onChangeMock).toHaveBeenCalledTimes(13)

        const submit = screen.getByText("Login")
        await user.click(submit)

        expect(onSubmitMock).toHaveBeenCalledOnce()
    })

})