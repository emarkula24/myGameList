import { useAuth } from "../utils/auth"
import styles from "./Header.module.css"
import SearchBar from "./SearchBar"
import { Link, } from "@tanstack/react-router"



// Header component, consists of HeaderElements, searchbar, and logout button
export default function Header() {
        const { user, isAuthenticated } = useAuth()
        const userParam: string = user?.username ?? ""
        const links: [string, string, string][] = [[
                "/community", "Community", ""
        ]]
        if (isAuthenticated) {
                links.push([`/gamelist/$username`, "GameList", userParam])
        }

        return (
                <div className={styles.headerContainer}>
                        <div>

                                {
                                        links.map(([to, label, params]) => {
                                                return (
                                                        <Link
                                                                key={to}
                                                                to={to}
                                                                params={{ username: params }}
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