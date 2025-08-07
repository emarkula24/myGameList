import styles from "./GameAddButton.module.css"

export interface GameAddButtonProps {
    onNewListEntry: (status: number) => void
}


export default function GameAddButton({ onNewListEntry }: GameAddButtonProps) {

    const handleClick = () => {
        onNewListEntry(1)
    }

    return (
        <div className={styles.container}>
            <button onClick={handleClick} type="button" className={styles.addBtn}>Add to GameList</button>
        </div>
    )
}