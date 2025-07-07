import React, { useEffect, useRef, useState } from "react";
import useDebounce from "../hooks/useDebounce";
import styles from "./SearchBar.module.css";
import { useNavigate } from "@tanstack/react-router";
import SearchResult from "./SearchResult";
import { fetchGames } from "../game";
import { useSearch } from "../hooks/useSearchContext";
import { useQuery, keepPreviousData } from "@tanstack/react-query";

export default function SearchBar() {
    const [searchQuery, setSearchQuery] = useState("");
    const [showResults, setShowResults] = useState(false);
    const navigate = useNavigate({});
    const { setSearchResults } = useSearch();
    const debouncedSearchQuery = useDebounce(searchQuery, 1000);
    const containerRef = useRef<HTMLDivElement | null>(null);

    const gameQuery = useQuery({
        queryKey: ["games", debouncedSearchQuery],
        queryFn: () => fetchGames(debouncedSearchQuery),
        enabled: !!debouncedSearchQuery,
        placeholderData: keepPreviousData,
    });

    useEffect(() => {
        if (gameQuery.isSuccess) {
            setSearchResults(gameQuery.data);
            
        }
    }, [gameQuery.data]);

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>): void => {
        setSearchQuery(event.currentTarget.value);
    };

    const handleEnterPress = (event: React.KeyboardEvent) => {
        if (event.key === "Enter") {
            navigate({ to: "/results" });
        }
    };

    // Detect clicks outside of container
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
                setShowResults(false);
            }
        };

        document.addEventListener("mousedown", handleClickOutside);
        return () => {
            document.removeEventListener("mousedown", handleClickOutside);
        };
    }, []);
    useEffect(() => {
        if (debouncedSearchQuery) {
            setShowResults(true);
        }
    }, [debouncedSearchQuery])
    return (
        <div ref={containerRef} className={styles.searchResults}>
            <label>
                Search:
                <input
                    type="text"
                    name="searchInput"
                    placeholder="Search for games"
                    value={searchQuery}
                    onChange={handleInputChange}
                    onKeyDown={handleEnterPress}
                />
            </label>
            <div>
                {gameQuery.isLoading && "Loading..."}
                {gameQuery.isFetched && gameQuery.data && showResults && (
                    <ul
                        style={{
                            display: "flex",
                            flexDirection: "column",
                            height: "600px",
                            overflowY: "scroll",
                            listStyleType: "none",
                            padding: 0,
                            margin: 0,
                        }}
                    >
                        {gameQuery.data.map((game) => (
                            <SearchResult key={game.id} game={game} />
                        ))}
                    </ul>
                )}
            </div>
        </div>
    );
}
