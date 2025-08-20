import { useEffect } from 'react'
import { Link } from '@tanstack/react-router'
import GameUpdateDropdown from './GameUpdateDropdown'
import { useUpdateGameMutation, useAddGameMutation } from '../queryOptions'
import { useAuth } from '../utils/auth'
import type { Game } from '../types/types'
import GameAddButton from './GameAddButton'
import styles from './GameRow.module.css'

const statusOptions: Record<number, string> = {
    1: "Playing",
    2: "Completed",
    3: "On-Hold",
    4: "Dropped",
    5: "Plan to Play"
}
const statusColors: Record<number, string> = {
    1: "green",      // Playing
    2: "blue",       // Completed
    3: "orange",     // On-Hold
    4: "red",        // Dropped
    5: "gray"        // Plan to Play
};

export default function GameRow({ index, game, username, isEditing, startEditing, stopEditing, isMissingFromLoggedInUserList }: {
    index: number
    game: Game
    username: string
    isEditing: boolean
    startEditing: (id: number) => void
    stopEditing: (id: number) => void
    isMissingFromLoggedInUserList: boolean
}) {
    const updateMutation = useUpdateGameMutation(username, game.id, game.name)
    const addMutation = useAddGameMutation(game.id, game.name)
    const auth = useAuth()
    const isLoggedInUser = auth.user?.username === username
    useEffect(() => {
        if (updateMutation.isSuccess) {
            stopEditing(game.id)  // Exit edit mode on success

        }
    }, [updateMutation.isSuccess])
    const statusColor = statusColors[game.status] || "lightgray";
    return (
        <tr className={styles.block}>
            <td><div style={{ backgroundColor: statusColor }} className={styles.titleColor}></div></td>
            <td className={styles.indexNumber}>{index}</td>
            <td className={styles.item}><img src={game.image.icon_url} /></td>
            <td>
                <Link to="/games/$guid" params={{ guid: game.guid }} className={styles.title}>
                    {game.name}
                </Link>
            </td>
            <td className={styles.item}>
                {isEditing ? (
                    <>
                        {updateMutation.isError ? (
                            <span style={{ color: "red" }}>Error</span>
                        ) : updateMutation.isPending ? (
                            <span style={{ color: "blue" }}>Updating...</span>
                        ) : (
                            <GameUpdateDropdown
                                onUpdateListEntry={(status) => updateMutation.mutate(status)}
                                status={game.status}
                            />
                        )}
                    </>
                ) : (
                    statusOptions[game.status]
                )}
            </td>
            {isLoggedInUser ? (

                <td>
                    <div className={styles.actionContainer}>
                        <div onClick={() => isEditing ? stopEditing(game.id) : startEditing(game.id)} className={styles.editBtn}>
                            {isEditing ? 'Cancel' : 'Edit'}
                        </div>
                        <div className={styles.deleteBtn}>Delete</div>
                    </div>
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
