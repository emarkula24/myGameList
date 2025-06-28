import axios from "axios"
import type { Game, Games } from "./types/types"
const url = import.meta.env.VITE_BACKEND_URL

export class GameNotFoundError extends Error { }
export class GamesNotFoundError extends Error { }

// Fetches info on a game based on a guid
export const fetchGame = async (guid: string) => {
    await new Promise((r) => setTimeout(r, 500))
    const game =  axios
        .get<{ results: Game }>(`${url}/games/game?guid=${guid}`)
        .then((r) => r.data.results)
        .catch((err) => {
            if (err.status === 404) {
                throw new GameNotFoundError(`Game with id ${guid} not found!`)
            }
            throw err
        })
    return game
}
// Fetches a list of games based on query string
export const fetchGames = async (searchQuery: string) => {
    const encodedSearchQuery = encodeURIComponent(searchQuery)
    await new Promise((r) => setTimeout(r, 500))
    return axios
        .get<{ results: Games[] }>(`${url}/games/search?query=${encodedSearchQuery}`)
        .then((response) => response.data.results)
        .catch((err) => {
            if (err.status === 500) {
                throw new GamesNotFoundError(`Games with query ${encodedSearchQuery} not found`)
            }
            throw err
        })

}