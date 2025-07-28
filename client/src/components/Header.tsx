import { useAuth } from "../utils/auth"
import styles from "./Header.module.css"
import SearchBar from "./SearchBar"
import { Link, } from "@tanstack/react-router"



// Header component, consists of HeaderElements, searchbar, and logout button
export default function Header() {
        const { user, isAuthenticated } = useAuth() 
        let userParam: string | undefined
        if (isAuthenticated) {
                userParam = user?.username
        } else {
                userParam = ""
        }
        return (
                <div className={styles.headerContainer}>
                        <div>

                        {(
                                [
                                        ["/community", "Community", ""],
                                        [`/gamelist/$username`, "GameList", userParam],
                                        
                                ] as const
                        ).map(([to, label, params]) => {
                                return (
                                        <Link
                                                key={to}
                                                to={to}
                                                params={{username: params}}
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