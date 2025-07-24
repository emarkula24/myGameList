import { useEffect } from 'react'
import { Link } from '@tanstack/react-router'
import GameUpdateDropdown from './GameUpdateDropdown'
import { createUpdateGameMutation } from '../queryOptions'
import { useAuth } from '../utils/auth'

const statusOptions: { [key: number]: string } = {
    1: "Playing",
    2: "Completed",
    3: "On-Hold",
    4: "Dropped",
    5: "Plan to Play"
}

export default function GameRow({ game, username, isEditing, toggleEditMode }: {
    game: any
    username: string
    isEditing: boolean
    toggleEditMode: (id: number) => void
}) {
    const updateMutation = createUpdateGameMutation(username)(game.id, game.name)
    const auth = useAuth()
    const isLoggedInUser = auth.user?.username === username
    useEffect(() => {
        if (updateMutation.isSuccess) {
            toggleEditMode(game.id)  // Exit edit mode on success
        }
    }, [updateMutation.isSuccess])

    return (
        <tr>
            <td><img src={game.image.icon_url} /></td>
            <td>
                <Link to="/games/$guid" params={{ guid: game.guid }}>
                    {game.name}
                </Link>
            </td>
            <td>
                {isEditing ? (
                    <>
                        <GameUpdateDropdown
                            onUpdateListEntry={(status) => updateMutation.mutate(status)}
                            status={game.status}
                        />
                        {updateMutation.isPending && <span>Updating...</span>}
                        {updateMutation.isError && (
                            <span style={{ color: 'red' }}>{updateMutation.error.message}</span>
                        )}
                    </>
                ) : (
                    statusOptions[game.status]
                )}
            </td>
            {isLoggedInUser &&
                <td onClick={() => toggleEditMode(game.id)} style={{ cursor: 'pointer' }}>
                    {isEditing ? 'Cancel' : 'Edit'}
                </td>
            }
        </tr>
    )
}
