import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import GameAddButton from './GameAddButton'

describe("<GameAddButton />", () => {
    
    test("function is executed when button is clicked", async () => {
        
        const onNewListEntryMock = vi.fn()

        render(
        <GameAddButton onNewListEntry={onNewListEntryMock} />
        )

        const user = userEvent.setup()
        const button = screen.getByRole("button", {name: "Add to GameList"})
        await user.click(button)

        expect(onNewListEntryMock).toHaveBeenCalledWith(1)
        expect(onNewListEntryMock).toHaveBeenCalledOnce()

    })
})