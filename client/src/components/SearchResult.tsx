import type { SearchedGame } from "../types/types";
import { useNavigate } from "@tanstack/react-router";
import styles from "./SearchResult.module.css"
interface SearchResultProps {
    game: SearchedGame;
}

const SearchResult: React.FC<SearchResultProps> = ({ game }) => {
    const navigate = useNavigate({})
    const guid = game.id.toString()
    const onMouseClick = () => {
        void navigate({
            // guid is the id value that is used to call for game specific information
            to: `/games/$guid`,
            params: { guid },
        })
    }
    
    const u = game.cover.url.replace("t_thumb", "t_cover_small")
    return (
        <div> 
            <li
                key={game.id}
                onClick={onMouseClick}
                className={styles.container}
                >
                <img src={u}  className={styles.resultImage}
                />
                
                <div className={styles.textContainer}>
                {game.name}
                </div>
                </li>
        
        </div>
    );
};

export default SearchResult;
