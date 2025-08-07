import type { Game } from "../types/types";
import styles from "./GamePlatformHeader.module.css"
export default function GamePlatformHeader({game}: {game: Game}) {
    return (
        <div className={styles.container}>
            {game.platforms.map((platform) =>
            <p>{platform.name}</p>)}
        </div>
    )
}