import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import GameUpdateDropdown from './GameUpdateDropdown'

describe("<GameUpdateDropdown />", () => {
    const onUpdateListEntryMock = vi.fn()

  it('renders with initial status text', () => {
    render(<GameUpdateDropdown status={3} onUpdateListEntry={onUpdateListEntryMock} />)
    expect(screen.getByText('On-Hold')).toBeInTheDocument()
  })

  it('toggles dropdown on button click', async () => {
    render(<GameUpdateDropdown status={1} onUpdateListEntry={onUpdateListEntryMock} />)
    const button = screen.getByText('Playing')

    // Initially dropdown content should not be visible
    expect(screen.queryByText('Completed')).not.toBeInTheDocument()

    // Click to open dropdown
    await userEvent.click(button)
    expect(screen.getByText('Completed')).toBeVisible()

    // Click again to close dropdown
    await userEvent.click(button)
    expect(screen.queryByText('Completed')).not.toBeInTheDocument()
  })


  it('calls onUpdateListEntry with correct status when selecting a new option', async () => {
    const onUpdateListEntry = onUpdateListEntryMock
    render(<GameUpdateDropdown status={1} onUpdateListEntry={onUpdateListEntry} />)
    const button = screen.getByText('Playing')
    await userEvent.click(button)

    const newOption = screen.getByText('Dropped')
    await userEvent.click(newOption)

    expect(onUpdateListEntry).toHaveBeenCalledTimes(1)
    expect(onUpdateListEntry).toHaveBeenCalledWith(4)

    // Dropdown should close after selection
    expect(screen.queryByText('Completed')).not.toBeInTheDocument()
  })

  it('does not call onUpdateListEntry when clicking the current status', async () => {
    const onUpdateListEntry = onUpdateListEntryMock
    render(<GameUpdateDropdown status={3} onUpdateListEntry={onUpdateListEntry} />)
    const button = screen.getAllByText('On-Hold')
    await userEvent.click(button[0])

    const currentOption = screen.getAllByText('On-Hold')
    await userEvent.click(currentOption[1])

    expect(onUpdateListEntry).not.toHaveBeenCalled()
  })

  it('closes dropdown when clicking outside', async () => {
    render(
      <>
        <GameUpdateDropdown status={1} onUpdateListEntry={onUpdateListEntryMock} />
        <div data-testid="outside">Outside Element</div>
      </>
    )

    const button = screen.getByText('Playing')
    await userEvent.click(button)
    expect(screen.getByText('Completed')).toBeVisible()

    // Click outside
    await userEvent.click(screen.getByTestId('outside'))

    expect(screen.queryByText('Completed')).not.toBeInTheDocument()
  })
})