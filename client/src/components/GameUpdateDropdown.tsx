import styles from './GameUpdateDropDown.module.css'
import { useState } from "react";

interface GameAddDropdownProps {
    onUpdateListEntry: (status: number) => void
    status: number
}

const statusOptions: Record<number, string> = {
    1: "Playing",
    2: "Completed",
    3: "On-Hold",
    4: "Dropped",
    5: "Plan to Play"
}

export default function GameUpdateDropdown({ onUpdateListEntry, status }: GameAddDropdownProps) {
    const [showDropdown, setShowDropdown] = useState(false);
    const numericStatus = Number(status)
    const handleSelect = (selectedStatus: number) => {
        setShowDropdown(false);

        if (selectedStatus === status) {
            console.log("status is the same so no update");
            return;
        }

        onUpdateListEntry(selectedStatus);
        console.log("tried to update", selectedStatus, statusOptions[selectedStatus]);
    };

    return (
        <div>
            <div
                onClick={() => setShowDropdown(prev => !prev)}
                className={`${styles.dropbtn} ${showDropdown ? styles.active : ""}`}
            >
                {statusOptions[status]}
            </div>
            <div className={styles.dropdownContainer}>
                {showDropdown && (
                    <div className={`${styles.dropdownContent} ${showDropdown ? styles.show : ""}`}>
                        {Object.entries(statusOptions).map(([key, label]) => {
                            const numericKey = Number(key);
                            const isCurrent = numericKey === numericStatus;
                            return (
                                <p
                                    key={numericKey}
                                    onClick={() => {
                                        if (!isCurrent) handleSelect(numericKey)
                                    }}
                                    className={`${styles.option} ${isCurrent ? styles.disabled : ""}`}
                                >
                                    {label}
                                </p>
                            )
                        })}
                    </div>
                )}
            </div>
        </div>
    );
}
