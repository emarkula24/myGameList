import axios from "axios"
import { keepPreviousData, useQuery } from '@tanstack/react-query'

import React, { useEffect, useState } from "react"
import useDebounce from "../../hooks/useDebounce"
import styles from "./SearchBar.module.css"
import { Games } from "../- types/types"
import { useNavigate } from "@tanstack/react-router"
import SearchResult from "./SearchResult"

const url = import.meta.env.VITE_BACKEND_URL

interface SearchBarProps {
        setSearchResults: React.Dispatch<React.SetStateAction<Games[]>>
}

const SearchBar: React.FC<SearchBarProps> = ({ setSearchResults }) => {
        const [searchQuery, setSearchQuery] = useState("")
        const navigate = useNavigate({})
        const debouncedSearchQuery = useDebounce(searchQuery, 500);



        const fetchGames = async (searchQuery: string) => {

                const encodedSearchQuery = encodeURIComponent(searchQuery)
                console.info(`${url}/search?query=${encodedSearchQuery}`)
                await new Promise((r) => setTimeout(r, 500))
                return axios
                        .get<{ results: Games[] }>(`${url}/games?query=${searchQuery}`)
                        .then((response) => response.data.results)
                        .catch((err) => {
                                console.error('Error fetching games:', err);
                                throw err
                        })

        }

        const { isFetched, isPending, isError, isSuccess, data, error } = useQuery({
                queryKey: ['games', debouncedSearchQuery],
                queryFn: () => fetchGames(debouncedSearchQuery),
                enabled: !!debouncedSearchQuery,
                placeholderData: keepPreviousData,
        });

        useEffect(() => {
                if (isSuccess) {
                        setSearchResults(data)
                }
        }, [data])

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
                        { data?.length === 0 && isFetched? (
                                <div>No results. </div>
                        ) : isPending ? (
                                <div>Loading..</div>
                        ) : isError ? (
                                <div>Error: {error.message}</div>
                        ) : isFetched && data && data.length > 0 &&
                        (
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
                                        {
                                        data.map((game) => (
                                                <SearchResult key={game.id}game={game} />
                                        ))}
                                </ul>
                        ) 
                        }
                </div>
        )
}

export default SearchBar