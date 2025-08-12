import type { Game } from "../types/types";
import styles from "./GamePlatformHeader.module.css"
export default function GamePlatformHeader({game}: {game: Game}) {
    return (
        <div className={styles.container}>
            {game.platforms.map((platform, index) =>
            <p key={index}>{platform.name}</p>)}
        </div>
    )
}