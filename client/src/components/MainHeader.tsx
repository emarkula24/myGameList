import { Link } from "@tanstack/react-router"
import styles from "./MainHeader.module.css"
import { useAuth } from "../utils/auth"
import LogoutButton from "./LogoutButton"

export default function MainHeader() {
    const auth = useAuth()
    return (
        <div className={styles.headerContainer}>
            <div className={styles.leftContent}>
                <Link to="/" className={styles.sitelogo}>myGameList</Link>
            </div>
            <div className={styles.rightContent}>
                {auth.isAuthenticated ? (
                    <div className={styles.profileContainer}>
                        <Link to="/profile" className={styles.usermenuProfile}>Profile</Link>
                        {auth.isAuthenticated && (
                            <LogoutButton />
                        )}
                    </div>
                ) : (
                    <div className={styles.usermenu}>
                        <Link to="/login" className={styles.usermenuLogin}>Login</Link>
                        <Link to="/register" className={styles.usermenuRegister}>Register</Link>
                    </div>
                )}

            </div>

        </div>
    )
}