import styles from "./Header.module.css"
import SearchBar from "./SearchBar"
import { Link, } from "@tanstack/react-router"



// Header component, consists of HeaderElements, searchbar, and logout button
export default function Header() {
        return (
                <div className={styles.headerContainer}>
                        <div>
                        {/* usually the address would be used as a key, index is for the sake of showing TBA items  */}

                        {(
                                [
                                        ["/community", "Community", 1],
                                        ["", "TBA", 2],
                                        ["", "TBA", 3],
                                        ["", "TBA", 4],
                                        
                                ] as const
                        ).map(([to, label, index]) => {
                                return (
                                        <Link
                                                key={index}
                                                to={to}
                                                activeProps={{ className: styles.active }}
                                                className={styles.linkButton}
                                        >
                                                {label}
                                        </Link>
                                )
                        })}
                        </div>
                        <div>
                        <SearchBar />
                        </div>
                </div>
        )
}