import axios, { AxiosError } from "axios"
import { type GameListEntry, type Game, type Games, type GameListEntries } from "./types/types"
import { UserNotLoggedInError } from "./utils/auth"
import axiosAuthorizationInstance from "./utils/axios"

export class GameNotFoundError extends Error { }
export class GamesNotFoundError extends Error { }
export class GameListNotFoundError extends Error { }
export class GameListEntryNotFoundError extends Error { }
export class GameListEmptyError extends Error { }

export const addGame = async (gameId: number, status: number, username: string | undefined, gamename: string) => {
    const result = axiosAuthorizationInstance
        .post(`list/add`, {
            game_id: gameId,
            status: status,
            username: username,
            gamename: gamename,

        })
        .then((r) => console.log(r))
        .catch((err: unknown) => {
            if (err instanceof AxiosError) {
                const errStatus = err.response?.status
                if (errStatus === 403) {
                    console.log(err.response)
                    throw new UserNotLoggedInError(`user not logged in when trying to add game ${gamename}`)
                } else {
                    console.log(err.response)
                    throw new Error
                }
            }

        })
    return result
}

export const updateGame = async (gameId: number, status: number, username: string | undefined, gamename: string): Promise<void> => {
    try {
        await axiosAuthorizationInstance.put(`/list/update`,
            {
                game_id: gameId,
                status: status,
                username: username,
                gamename: gamename,
            })
    }
    catch (err) {
        if (err instanceof AxiosError) {
            const errStatus = err.response?.status
            if (errStatus === 403) {
                console.log(err.response)
                throw new UserNotLoggedInError(`user not logged in when trying to update game ${gamename}`)
            } else {
                console.log(err.response)
                throw new Error
                
            }
        }
        throw new Error
    }
}
// Fetches info on a game based on a guid
export const fetchGame = async (guid: string): Promise<Game> => {
    try {
        const response = await axios.get<{ results: Game }>(`/games/game?guid=${guid}`);
        return response.data.results;
    } catch (err) {
        if (err instanceof AxiosError) {
            if (err.response?.status === 404) {
                throw new GameNotFoundError(`Game with id ${guid} not found!`);
            }
        }
        throw err;
    }
};

// Fetches a list of games based on query string
export const fetchGames = async (searchQuery: string): Promise<Games[]> => {
    const encodedSearchQuery = encodeURIComponent(searchQuery)
    try {
        // An empty promise is called in order to create a delay between query changes and actual api calls
        // Rate limiting on the backend side will be implemented eventually
        await new Promise((r) => setTimeout(r, 500))

        const response = await axios.get<{ results: Games[] }>(`/games/search?query=${encodedSearchQuery}`)
        return response.data.results
    } catch (err) {
        if (err instanceof AxiosError) {
            if (err.status === 500) {
                throw new GamesNotFoundError(`Games with query ${encodedSearchQuery} not found`)
            }
            if (err.status === 400) {
                throw new GameListEmptyError
            }

        }
        throw err
    }

}

export const fetchGameList = async (username: string | undefined, page = 1, limit = 20): Promise<GameListEntries[]> => {
    try {
        const response = await axios.get<{ results: GameListEntries[] }>(`/list?username=${username}&page=${page}&limit=${limit}`)
        return response.data.results
    } catch (err) {
        if (err instanceof AxiosError) {
            const errStatus = err.response?.status
            if (errStatus === 500 || errStatus == 404) {
                throw new GameListNotFoundError(`no gamelist found for user`)
            }

        }
        throw err
    }
}
export const fetchGameListEntry = async (username: string | undefined, gameId: number): Promise<GameListEntry> => {
    try {
        const response = await axios.get<GameListEntry>(`list/game?username=${username}&gameId=${gameId}`)
        return response.data
    } catch (err) {
        if (err instanceof AxiosError) {
            const errStatus = err.response?.status
            if (errStatus === 500 || errStatus === 404) {
                throw new GameListEntryNotFoundError(`no game entry found for user`)
            }

        }
        throw err
    }
}