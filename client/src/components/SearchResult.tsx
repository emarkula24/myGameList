import type { Games } from "../types/types";
import { useNavigate } from "@tanstack/react-router";
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
            mask: {
                to: `/games/${game.id}`,
            },

        })
    }
    return (
        <> 
            <li
                key={game.id}
                onClick={onMouseClick}
                style={{
                    flex: '0 0 auto',
                    marginRight: '16px',
                }}><img src={game.image?.icon_url} />{game.name}</li>
        
        </>
    );
};

export default SearchResult;
