import { queryOptions } from "@tanstack/react-query";
import { fetchGame } from "../../games";

export const gameQueryOptions = (guid: string)  => 
    queryOptions({
        queryKey: ["game", { guid} ],
        queryFn: () => fetchGame(guid),
    })
