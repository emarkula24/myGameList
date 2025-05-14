import { queryOptions, QueryOptions } from "@tanstack/react-query";
import { fetchGame } from "../../games";

export const gameQueryOptions = (gameGuid: string)  => 
    queryOptions({
        queryKey: ["game", { gameGuid} ],
        queryFn: () => fetchGame(gameGuid),
    })
