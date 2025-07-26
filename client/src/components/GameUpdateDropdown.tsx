import styles from './GameUpdateDropdown.module.css'
import { useEffect, useState } from "react";

type GameAddDropdownProps = {
    onUpdateListEntry: (status: number) => void
    status: number
}

const statusOptions: { [key: number]: string } = {
    1: "Playing",
    2: "Completed",
    3: "On-Hold",
    4: "Dropped",
    5: "Plan to Play"
}

export default function GameUpdateDropdown({ onUpdateListEntry, status }: GameAddDropdownProps) {
    const [showDropdown, setShowDropdown] = useState(false);
    const [currentStatus, setCurrentStatus] = useState(status);

    useEffect(() => {
        setCurrentStatus(status);
    }, [status]);


    const handleSelect = (selectedStatus: number) => {
        setShowDropdown(false); // close dropdown after selection

        if (selectedStatus === currentStatus) {
            console.log("status is the same so no update")
            return
        }
        setCurrentStatus(selectedStatus)
        onUpdateListEntry(selectedStatus)
        console.log("tried to update", selectedStatus, statusOptions[selectedStatus])

    };

    return (
        <div>
            <button onClick={() => setShowDropdown(prev => !prev)} className={styles.dropbtn}>
                {statusOptions[currentStatus]}
            </button>
            <div className={styles.dropdown}>
            {showDropdown && (
                <div className={`${styles.dropdownContent} ${showDropdown ? styles.show : ""}`}>
                    {Object.entries(statusOptions).map(([key, label]) => {
                        const numericKey = Number(key);
                        const isCurrent = numericKey === currentStatus;
                        return (
                            <p
                                key={numericKey}
                                onClick={() => {
                                    if (!isCurrent) handleSelect(numericKey)
                                }}
                                className={isCurrent ? styles.disabled : undefined}
                            >
                                {label}
                            </p>
                        )
                    })}
                </div>
            )}
        </div>
        </div>
    )
}