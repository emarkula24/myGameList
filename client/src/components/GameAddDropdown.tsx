import styles from './GameAddDropdown.module.css'
import { useState } from "react";
export default function GameAddDropdown({ onSelect }: { onSelect: (status: string) => void }) {
    const [showDropdown, setShowDropdown] = useState(false);
    const [defaultText, setDefaultText] = useState("Select Status")
    const [selectedStatus, setSelectedStatus] = useState<string | null>(null);

    const handleSelect = (gameStatus: string, label: string) => {
        setDefaultText(label)
        setShowDropdown(false); // close dropdown after selection
        setSelectedStatus(gameStatus)
    };
    const handleClick = () => {
        if (selectedStatus) {
            onSelect(selectedStatus)
        }

    }
    return (
        <div className={styles.dropdown}>
            <button onClick={() => setShowDropdown(prev => !prev)} className={styles.dropbtn}>
                {defaultText}
            </button>
            <div className={`${styles.dropdownContent} ${showDropdown ? styles.show : ""}`}>
                {(
                    [
                        ["playing", "Playing"],
                        ["completed", "Completed"],
                        ["on-hold", "On-Hold"],
                        ["dropped", "Dropped"],
                        ["plan to play", "Plan to Play"],
                    ] as const
                ).map(([gameStatus, label]) => {
                    return (<p
                        key={label}
                        onClick={() => handleSelect(gameStatus, label)}
                    >
                        {label}
                    </p>)
                })
                }

            </div>
                        <button
                onClick={handleClick}
                disabled={!selectedStatus}
                className={styles.addButton}
                >
                Add to List
            </button>
        </div>
    )
}