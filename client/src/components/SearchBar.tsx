
import { keepPreviousData, useQuery } from '@tanstack/react-query'

import React, { useEffect, useState } from "react"
import useDebounce from "../hooks/useDebounce"
import styles from "./SearchBar.module.css"
import { useNavigate } from "@tanstack/react-router"
import SearchResult from "./SearchResult"
import { fetchGames } from "../game"
import { useSearch } from '../hooks/useSearchContext'




export default function SearchBar() {
        const [searchQuery, setSearchQuery] = useState("")
        const navigate = useNavigate({})
        const { setSearchResults } = useSearch()
        const debouncedSearchQuery = useDebounce(searchQuery, 150);

        const gameQuery = useQuery({
                queryKey: ['games', debouncedSearchQuery],
                queryFn: () => fetchGames(debouncedSearchQuery),
                enabled: !!debouncedSearchQuery,
                placeholderData: keepPreviousData,
        });

        useEffect(() => {
                if (gameQuery.isSuccess) {
                        setSearchResults(gameQuery.data)
                }
                console.log("new content loaded")
        }, [gameQuery.data])

        const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>): void => {
                setSearchQuery(event.currentTarget.value)
        }

        const handleEnterPress = (event: React.KeyboardEvent) => {
                if (event.key === "Enter") {
                        navigate({ to: "/results" })
                }

        }

        return (
                <div className={styles.searchResults}>
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
                                {gameQuery.isLoading ? (
                                        "Loading..."
                                ) : gameQuery.isFetched && gameQuery.data ?  (
                                        <>
                                                <ul
                                                        style={{
                                                                display: "flex",
                                                                flexDirection: "column",
                                                                height: "600px",
                                                                overflowY: "scroll",
                                                                listStyleType: 'none',
                                                                padding: 0,
                                                                margin: 0,
                                                        }}
                                                >
                                                        {gameQuery.data?.map((game) => (
                                                                <SearchResult key={game.id} game={game} />
                                                        ))}
                                                </ul>
                                        </>
                                ) : !gameQuery.isFetching && (<></>)}
                        </div>
                </div>
        )
}
