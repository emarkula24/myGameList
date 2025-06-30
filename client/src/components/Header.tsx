import styles from "./Header.module.css"
import SearchBar from "./SearchBar"
import { Link, } from "@tanstack/react-router"



// Header component, consists of HeaderElements and searchbar
export default function Header() {
        
        return (
                <div className={styles.headerContainer}>
                        {(
                                [
                                        ["/", "Home", true],
                                        ["/about", "About"],
                                        ["/community", "Community"],
                                        ["/profile", "Profile"],
                                        ["/login", "Login"],
                                        ["/register", "Sign up"]
                                ] as const
                        ).map(([to, label, exact]) => {
                                return (
                                        <Link
                                        key={to}
                                        to={to}
                                        activeOptions={{exact}}
                                        activeProps={{className: styles.active}}
                                        className={styles.linkButton}
                                        >
                                        {label}
                                        </Link>
                                )
                        })}
                        <SearchBar/>
                </div>
        )
}