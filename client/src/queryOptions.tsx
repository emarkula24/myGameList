import { queryOptions} from "@tanstack/react-query";
import { fetchGame, fetchGameList, fetchGameListEntry, fetchGames } from "./game";

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
export const gameListQueryOptions = (username: string) =>
    queryOptions({
        queryKey: ["gamelist", {username}],
        queryFn: () => fetchGameList(username)
    })
export const gameListEntryQueryOptions = (username: string | undefined, gameId: number) =>
    queryOptions({
        queryKey: ["gamelistsingle", { username, gameId}],
        queryFn: () => fetchGameListEntry(username, gameId)
    })