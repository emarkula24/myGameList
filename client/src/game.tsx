import axios from "axios"
import { type GameListEntry, type Game, type Games } from "./types/types"
import { UserNotLoggedInError } from "./utils/auth"
import axiosAuthorizationInstance from "./utils/axios"

export class GameNotFoundError extends Error { }
export class GamesNotFoundError extends Error { }
export class GameListNotFoundError extends Error { }

export const addGame = async (gameId: number, status: string, username: string | undefined, gamename: string) => {
    await new Promise((r) => setTimeout(r, 500))
    const result = axiosAuthorizationInstance
    .post(`list/add`, {
        game_id: gameId,
        status: status,
        username: username,
        gamename: gamename,

    })
    .then((r) => console.log(r))
    .catch((err) => {
        const errStatus = err.response?.status
        if (errStatus === 403 ) {
            console.log(err.response)
            throw new UserNotLoggedInError(`user not logged in when trying to add game ${gamename}`)
        } else {
            console.log(err.response)
            throw new Error
        }
    })
    return result
}
// Fetches info on a game based on a guid
export const fetchGame = async (guid: string) => {
    await new Promise((r) => setTimeout(r, 500))
    const game = axios
        .get<{ results: Game }>(`/games/game?guid=${guid}`)
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
        .get<{ results: Games[] }>(`/games/search?query=${encodedSearchQuery}`)
        .then((response) => response.data.results)
        .catch((err) => {
            if (err.status === 500) {
                throw new GamesNotFoundError(`Games with query ${encodedSearchQuery} not found`)
            }
            throw err
        })

}

export const fetchGameList = async (username: string, page: number = 1, limit: number = 20) => {
    
    await new Promise((r) => setTimeout(r, 500))
    return axios
        .get<{results: GameListEntry[]}>(`/list?username=${username}&page=${page}&limit=${limit}`)
        .then((response) => response.data.results)
        .catch((err) => {
            const errStatus = err.response?.status
            if (errStatus === 500 || errStatus == 404) {
                throw new GameListNotFoundError(`no gamelist found for user`)
            }
            throw err
        })
}