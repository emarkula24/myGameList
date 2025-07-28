import githubLogo from "../assets/github-mark.png"
import styles from "./Footer.module.css"
export default function Footer() {
    return (
        <div className="routeContainer">
            <div className={styles.footerContainer}>
                <footer >
                    <img src={githubLogo} className={styles.image} />
                </footer>
            </div>
        </div>
    )
}