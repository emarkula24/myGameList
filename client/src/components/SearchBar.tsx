import React, { useEffect, useRef, useState } from "react";
import useDebounce from "../hooks/useDebounce";
import styles from "./SearchBar.module.css";
import { Link, useRouter } from "@tanstack/react-router";
import SearchResult from "./SearchResult";
import { fetchGames } from "../game";
import { useQuery, keepPreviousData } from "@tanstack/react-query";

export default function SearchBar() {
    const [searchQuery, setSearchQuery] = useState("");
    const [showResults, setShowResults] = useState(false);
    const debouncedSearchQuery = useDebounce(searchQuery, 200);
    const containerRef = useRef<HTMLDivElement | null>(null);
    const router = useRouter()

    const gameQuery = useQuery({
        queryKey: ["games", debouncedSearchQuery],
        queryFn: () => fetchGames(debouncedSearchQuery),
        enabled: !!debouncedSearchQuery,
        placeholderData: keepPreviousData,
    });


    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>): void => {
        setSearchQuery(event.currentTarget.value)
        setShowResults(true)
    };

    const handleEnterPress = (event: React.KeyboardEvent) => {
        if (event.key === "Enter") {
            setShowResults(false)
            void router.navigate({ to: "/results/$query", params: {query : searchQuery}});
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

    return (
        <div ref={containerRef} className={styles.searchResults}>
            <label>
                <input
                    type="text"
                    name="searchInput"
                    placeholder="Search for games"
                    value={searchQuery}
                    onChange={handleInputChange}
                    onKeyDown={handleEnterPress}
                    onBlur={() => setSearchQuery("")}
                    autoComplete="off"
                    spellCheck="false"
                />
            </label>
            <div>
                {gameQuery.isLoading && <div>loading..</div>}
                {gameQuery.isFetched && gameQuery.data && gameQuery.data.length > 1 && showResults && (
                    <ul
                        className={styles.listContainer}
                    >
                        {gameQuery.data.map((game) => (
                            <SearchResult key={game.id} game={game} />
                        ))}
                        
                        <Link to="/results/$query" params={{query: searchQuery}}><li>View all results for <b>{searchQuery}</b></li></Link>
                    </ul>
                )}
            </div>
        </div>
    );
}
