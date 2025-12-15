import { useNavigate } from "@tanstack/react-router";
import type { SearchedGame } from "../types/types";
import styles from "./GameResultRow.module.css"

export default function GameResultRow({ game }: { game: SearchedGame }) {
    const navigate = useNavigate({})
    const guid = game.id.toString()
    const onMouseClick = () => {
        void navigate({
            // guid is the value that is used to call for game specific information
            to: `/games/$guid`,
            params: { guid },
        })
    }
    const u = game.cover.url.replace("t_thumb", "t_cover_small")
    return (
        <div className={styles.container} onClick={onMouseClick}>
            <div key={game.id}>
                <h3 style={{fontSize: "1.6em"}}>{game.name}</h3>
                <img src={u} style={{maxWidth: "100%", height: "auto"}}/>
                <h4>{game.platforms.map(game => game.abbreviation)}</h4>
            </div>
        </div> 
    )
}
