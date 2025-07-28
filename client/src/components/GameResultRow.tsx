import { useNavigate } from "@tanstack/react-router";
import type { Games } from "../types/types";
import styles from "./GameResultRow.module.css"

export default function GameResultRow({ game }: { game: Games }) {
    const navigate = useNavigate({})
    const guid = game.guid.toString()
    const onMouseClick = () => {
        navigate({
            // guid is the value that is used to call for game specific information
            to: `/games/$guid`,
            params: { guid },
        })
    }
    return (
        <div className={styles.container} onClick={onMouseClick}>
            <div key={game.id}>
                <h3>{game.name}</h3>
                <img src={game.image?.thumb_url} style={{maxWidth: "100%", height: "auto"}}/>
            </div>
        </div>
    )
}
