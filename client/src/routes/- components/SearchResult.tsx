import { Game } from "../- types/types";
import { useNavigate } from "@tanstack/react-router";
interface SearchResultProps {
    game: Game;
}

const SearchResult: React.FC<SearchResultProps> = ({ game }) => {
    const navigate = useNavigate({})

    const handleMouseClick = ()  =>{
        navigate({
            // guid is the value that is used to call for game specific information
            to: `/game/${game.guid}`,
            mask: {
                to: `/game/${game.id}`,
            }
        })
    }
    return (
        <>
            <li 
            key={game.id} 
            onClick={handleMouseClick}
            style={{
                flex: '0 0 auto',
                marginRight: '16px',
            }}><img src={game.image?.icon_url} />{game.name}</li>
        </>
    );
};

export default SearchResult;
