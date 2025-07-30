import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import GameListFilterHeader from './GameListFilterHeader'

describe("<GameListFilterHeader />", () => {
    const onSelectMock = vi.fn()
    const setSearchQueryMock = vi.fn()

    test("clicking on statusItem with text triggers function with matching id", async () => {
        render(<GameListFilterHeader onSelect={onSelectMock} setSearchQuery={setSearchQueryMock} />)
        const user = userEvent.setup()
        const statusItem = screen.getByText("Completed") //id: 2
        const droppedItem = screen.getByText("Dropped");   // id: 4
        await user.click(statusItem)
        await user.click(droppedItem)

        expect(onSelectMock).toHaveBeenCalledTimes(2)
        expect(onSelectMock).toHaveBeenNthCalledWith(1,2)
        expect(onSelectMock).toHaveBeenNthCalledWith(2,4)

    })
    test("writing in the search bar calls setSearchQuery with expected value", async () => {
        render(<GameListFilterHeader onSelect={onSelectMock} setSearchQuery={setSearchQueryMock} />)
        const user = userEvent.setup();
        const searchInput = screen.getByPlaceholderText("Search list");
        await user.type(searchInput, 'metroid');

        expect(setSearchQueryMock).toHaveBeenCalledTimes(7)
        expect(setSearchQueryMock).toHaveBeenLastCalledWith("metroid")
    })
    test("clicking search icon in active state clears search input", async () => {
        render(<GameListFilterHeader onSelect={onSelectMock} setSearchQuery={setSearchQueryMock} />)
        const user = userEvent.setup();
        const searchInput = screen.getByPlaceholderText("Search list");
        const img = screen.getByRole("img")
        await user.click(img)
        await user.type(searchInput, "mario")
        expect(setSearchQueryMock).toHaveBeenLastCalledWith('mario');
        await user.click(img)
        expect(setSearchQueryMock).toHaveBeenCalledWith("")


    })
})