import axios from "axios"
import type { Game } from "./types/types"
const url = import.meta.env.VITE_BACKEND_URL

export class GameNotFoundError extends Error {}

export const fetchGame = async (guid: string) => {
    await new Promise((r) => setTimeout(r, 500))
    const game = await axios
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