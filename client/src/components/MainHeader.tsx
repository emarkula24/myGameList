import { Link } from "@tanstack/react-router"
import styles from "./MainHeader.module.css"

export default function MainHeader() {
    return (
        <div className={styles.headerContainer}>
            <div className={styles.leftContent}>
                <Link to="/">myGameList</Link>
            </div>
            <div className={styles.rightContent}>

                Profile
            </div>

        </div>
    )
}