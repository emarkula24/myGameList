import styles from "./GameListFilterHeader.module.css"
import { useState } from "react";
export default function GameListFilterHeader({ onSelect }: { onSelect: React.Dispatch<React.SetStateAction<number>> }) {
    const [currentSelection, setCurrentSelecton] = useState(0)
    const statusItems = [
        { id: 0, label: "All Games" },
        { id: 1, label: "Playing" },
        { id: 2, label: "Completed" },
        { id: 3, label: "On-Hold" },
        { id: 4, label: "Dropped" },
        { id: 5, label: "Plan to Play" }
    ];

    const handleClick = (id: number) => {
        setCurrentSelecton(id)
        onSelect(id)
    }
    return (
        <div className={styles.statusContainer}>
            {statusItems.map((item) =>
                <div key={item.id} onClick={() => handleClick(item.id)}
                    className={`${styles.statusItem} ${currentSelection === item.id ? styles.activeSelection : ""
                        }`}
                >
                    {item.label}
                </div>
            )}
        </div>
    )
}