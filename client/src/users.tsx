import type { Users } from "./types/types"
import axios, { AxiosError } from "axios"

export const fetchUsers = async (): Promise<Users[]> => {
    try {
        const response  = await axios.get<Users[]>(`user/users`)
        return response.data
    } catch(err) {
        if (err instanceof AxiosError) {
            const errStatus = err.response?.status
            throw Error(`failed to load users: ${errStatus}`)
        }
        throw err
    }
}