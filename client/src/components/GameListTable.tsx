import type { GameListEntries } from '../types/types'
import GameRow from './GameRow'
import GameTableHeaderRow from './GameTableHeaderRow'

export function GameListTable({
  games,
  username,
  editingGameIds,
  startEditing,
  stopEditing,
  loggedInGameIds
}: {
  games: GameListEntries[]
  username: string
  editingGameIds: Set<number>
  startEditing: (id: number) => void
  stopEditing: (id: number) => void
  loggedInGameIds: Set<number>
}) {
  if (games.length === 0) {
    return <p style={{ fontSize: "2em", textAlign: "center" }}>This category is empty</p>
  }

  return (
    <table style={{ borderCollapse: "collapse" }}>
      <thead>
        <GameTableHeaderRow />
      </thead>
      <tbody>
        {games.map((game, index) => (
          <GameRow
            index={index + 1}
            key={game.id}
            game={game}
            username={username}
            isEditing={editingGameIds.has(game.id)}
            startEditing={startEditing}
            stopEditing={stopEditing}
            isMissingFromLoggedInUserList={!loggedInGameIds.has(game.id)}
          />
        ))}
      </tbody>
    </table>
  );
}
