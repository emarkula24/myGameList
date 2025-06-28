import { queryOptions} from "@tanstack/react-query";
import { fetchGame, fetchGames } from "./game";

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
