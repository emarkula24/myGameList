import { useEffect } from 'react'
import { Link } from '@tanstack/react-router'
import GameUpdateDropdown from './GameUpdateDropdown'
import { useUpdateGameMutation, useAddGameMutation } from '../queryOptions'
import { useAuth } from '../utils/auth'
import type { Game } from '../types/types'
import GameAddButton from './GameAddButton'
import styles from './GameRow.module.css'

const statusOptions: { [key: number]: string } = {
    1: "Playing",
    2: "Completed",
    3: "On-Hold",
    4: "Dropped",
    5: "Plan to Play"
}

export default function GameRow({ index, game, username, isEditing, startEditing, stopEditing, onUpdateSuccess, isMissingFromLoggedInUserList }: {
    index: number
    game: Game
    username: string
    isEditing: boolean
    startEditing: (id: number) => void
    stopEditing: (id: number) => void
    onUpdateSuccess: (id: number) => void
    isMissingFromLoggedInUserList: boolean
}) {
    const updateMutation = useUpdateGameMutation(username, game.id, game.name)
    const addMutation = useAddGameMutation(game.id, game.name)
    const auth = useAuth()
    const isLoggedInUser = auth.user?.username === username
    useEffect(() => {
        if (updateMutation.isSuccess) {
            onUpdateSuccess(game.id)  // Exit edit mode on success

        }
    }, [updateMutation.isSuccess])

    return (
            <tr className={styles.block}>
                <td className={styles.item}>{index}</td>
                <td className={styles.item}><img src={game.image.icon_url} /></td>
                <td>
                    <Link to="/games/$guid" params={{ guid: game.guid }} className={styles.title}>
                        {game.name}
                    </Link>
                </td>
                <td className={styles.item}>
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
                            {updateMutation.isSuccess && isEditing}
                        </>
                    ) : (

                        statusOptions[game.status]
                    )}
                </td>
                {isLoggedInUser ? (
                    <td onClick={() => isEditing ? stopEditing(game.id) : startEditing(game.id)} className={styles.editBtn}>
                        {isEditing ? 'Cancel' : 'Edit'}
                    </td>
                ) : isMissingFromLoggedInUserList && (
                    <td>
                        <GameAddButton
                            onNewListEntry={(status) => addMutation.mutate(status)}
                        />
                    </td>
                )
                }
            </tr>
    )
}
