import { queryOptions, useMutation, useQueryClient } from "@tanstack/react-query";
import { addGame, deleteGame, fetchGame, fetchGameList, fetchGameListEntry, fetchGames, updateGame } from "./game";
import { useAuth } from "./utils/auth";
import { fetchUsers } from "./users";

export const gameQueryOptions = (guid: string) =>
    queryOptions({
        queryKey: ["game", { guid }],
        queryFn: () => fetchGame(guid),
    })

export const gamesQueryOptions = (query: string) =>
    queryOptions({
        queryKey: ["games", { query }],
        queryFn: () => fetchGames(query),
    })
export const gameListQueryOptions = (username: string | undefined) =>
    queryOptions({
        queryKey: ["gamelist", { username }],
        queryFn: () => fetchGameList(username),

    })
export const gameListEntryQueryOptions = (username: string | undefined, gameId: number) =>
    queryOptions({
        queryKey: ["gamelistentry", { username, gameId }],
        queryFn: () => fetchGameListEntry(username, gameId),
    })

export const usersQueryOptions = () =>
    queryOptions({
        queryKey: ["users"],
        queryFn: () => fetchUsers()
    }) 

export function useAddGameMutation(gameId: number, gameName: string) {
    const queryClient = useQueryClient()
    const auth = useAuth()
    const username = auth.user?.username

    return useMutation({
        mutationFn: async (status: number) => {
            return await addGame(gameId, status, username, gameName)
        },
        onSettled: () => queryClient.invalidateQueries({ queryKey: ["gamelistentry"] }),
        onError: (error) => {
            console.error('Add game error:', error)
        }
    })
}
export function useUpdateGameMutation(username: string | undefined, gameId: number, gameName: string) {
    const auth = useAuth()
    const queryClient = useQueryClient()
    return useMutation({

        mutationFn: async (status: number) => {
            if (auth.user?.username !== username) {
                throw new Error("you are not this user")
            }
            return await updateGame(gameId, status, auth.user?.username, gameName)
        },
        onSettled: () => queryClient.invalidateQueries({ queryKey: ["gamelist"] }),
        onError: (error) => {
            console.error('Update game error:', error)
        }

    })
}

export function useDeleteGameMutation (gameId: number, username: string | undefined ) {
    const auth = useAuth()
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async () => {
            if (auth.user?.username !== username) {
                throw new Error("you are not this user")
            }
            return await deleteGame(gameId, username)
        },
        onSettled: () => queryClient.invalidateQueries({queryKey: ["gamelist"]}),
        onError: (error) => {
            console.error("Delete game error:", error)
        }
    })
}
