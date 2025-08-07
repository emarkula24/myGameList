// GameListTable.test.tsx
import { render, screen } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import { GameListTable } from './GameListTable'
import type { Game } from '../types/types'

// You can mock GameRow and GameTableHeaderRow if you want to isolate logic
vi.mock('./GameRow', () => ({
  default: ({ game }: any) => <tr data-testid="game-row"><td>{game.name}</td></tr>
}))
vi.mock('./GameTableHeaderRow', () => ({
  default: () => <tr><th>Header</th></tr>
}))

const mockGames: Game[] = [
  { id: 1, name: 'Game One' },
  { id: 2, name: 'Game Two' }
] as Game[]

describe('<GameListTable />', () => {
  it('renders empty message when no games are provided', () => {
    render(
      <GameListTable
        games={[]}
        username="testUser"
        editingGameIds={new Set()}
        startEditing={() => {}}
        stopEditing={() => {}}
        loggedInGameIds={new Set()}
      />
    )

    expect(screen.getByText(/this category is empty/i)).toBeInTheDocument()
  })

  it('renders the correct number of game rows', () => {
    render(
      <GameListTable
        games={mockGames}
        username="testUser"
        editingGameIds={new Set()}
        startEditing={() => {}}
        stopEditing={() => {}}
        loggedInGameIds={new Set([1, 2])}
      />
    )

   const rows = screen.getAllByTestId('game-row')
    console.log(rows)
    expect(rows.length).toBe(2)
    expect(rows[0]).toHaveTextContent("Game One")
    expect(rows[1]).toHaveTextContent("Game Two")
  })
})
