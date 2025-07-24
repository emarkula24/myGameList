import styles from "./GameListFilterHeader.module.css"
export default function GameListFilterHeader() {
    const statusItems = [
        { id: 0, label: "All Games"},
        { id: 1, label: "Playing" },
        { id: 2, label: "Completed" },
        { id: 3, label: "On-Hold" },
        { id: 4, label: "Dropped" },
        { id: 5, label: "Plan to Play" }
    ];
    const handleClick = () => {

    }
    return (
        <div className={styles.statusContainer}>
            {statusItems.map((item) =>
                <div key={item.id} onClick={handleClick}className={styles.statusItem}>{item.label}</div>
            )}
        </div>
    )
}