import type { Games } from "../types/types";
import { useNavigate } from "@tanstack/react-router";
import styles from "./SearchResult.module.css"
interface SearchResultProps {
    game: Games;
}

const SearchResult: React.FC<SearchResultProps> = ({ game }) => {
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
        <> 
            <li
                key={game.id}
                onClick={onMouseClick}
                className={styles.resultItem}><img src={game.image?.icon_url}  className={styles.resultImage}/>{game.name}</li>
        
        </>
    );
};

export default SearchResult;
