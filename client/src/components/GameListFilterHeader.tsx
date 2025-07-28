import styles from "./GameListFilterHeader.module.css"
import { useRef, useState } from "react";
import searchIcon from '../assets/search_icon.png';
export default function GameListFilterHeader({ onSelect, setSearchQuery }: {
    onSelect: React.Dispatch<React.SetStateAction<number>>,
    setSearchQuery: React.Dispatch<React.SetStateAction<string>>
}) {
    const [currentSelection, setCurrentSelecton] = useState(0)
    const [searchActive, setSearchActive] = useState(false)
    const inputRef = useRef<HTMLInputElement>(null);

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
    const handleInput = (event: React.ChangeEvent<HTMLInputElement>) => {
        setSearchQuery(event.currentTarget.value)
    }
    const handleSearchClick = () => {
        setSearchActive(prev => !prev)
        setSearchQuery("")
        if (inputRef.current) {
            inputRef.current.value = ""
        }
    }

    return (
        <div className={styles.container}>
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
            <div className={styles.searchContainer}>
            <input ref={inputRef}type="text" maxLength={20} onChange={handleInput} className={`${styles.input} ${searchActive ? styles.searchActive : styles.searchInActive
                }`} />
            <img src={searchIcon} onClick={handleSearchClick} />
            </div>
        </div>
    )
}