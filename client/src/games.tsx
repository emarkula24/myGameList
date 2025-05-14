import axios from "axios"
import { Game } from "./routes/- types/types"

const url = import.meta.env.VITE_BACKEND_URL

export class GameNotFoundError extends Error {}

export const fetchGame = async (gameGuid: string) => {
    await new Promise((r) => setTimeout(r, 500))
    const game = await axios
    .get<Game>(`${url}/game/${gameGuid}`)
    .then((r) => r.data)
    .catch((err) => {
        if (err.status === 404) {
            throw new GameNotFoundError(`Game with id ${gameGuid} not found!`)
        }
        throw err
    })
    return game
}