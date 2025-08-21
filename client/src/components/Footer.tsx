import { Link } from "@tanstack/react-router"
import githubLogo from "../assets/github-mark.png"
import styles from "./Footer.module.css"
export default function Footer() {
    return (
        <>
            <div className={styles.footerContainer}>
                <div className={styles.linkContainer}>
                    <a 
                    href="https://github.com/emarkula24/"
                    target="_blank"
                    rel="noopener noreferrer"
                    >
                    <img src={githubLogo} className={styles.image} />
                    </a>
                    <Link to="/about" className={styles.link}>About</Link>
                </div>
                <footer style={{fontSize: "1.2em"}}>
                    @{new Date().getFullYear()} partavesipirtelo@protonmail.com
                </footer>
            </div>
        </>
    )
}