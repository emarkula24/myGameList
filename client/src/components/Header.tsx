import { useAuth } from "../utils/auth"
import styles from "./Header.module.css"
import LogoutButton from "./LogoutButton"
import SearchBar from "./SearchBar"
import { Link, } from "@tanstack/react-router"



// Header component, consists of HeaderElements, searchbar, and logout button
export default function Header() {
        const auth = useAuth()
        return (
                <div className={styles.headerContainer}>

                        {(
                                [
                                        ["/about", "About"],
                                        ["/community", "Community"],
                                        ["/profile", "Profile"],
                                        ["/login", "Login"],
                                        ["/register", "Sign up"]
                                ] as const
                        ).map(([to, label]) => {
                                return (
                                        <Link
                                                key={to}
                                                to={to}
                                                
                                                activeProps={{ className: styles.active }}
                                                className={styles.linkButton}
                                        >

                                                {label}

                                        </Link>
                                )
                        })}
                        <SearchBar />
                        {auth.isAuthenticated && (
                                <LogoutButton />
                        )}

                </div>
        )
}