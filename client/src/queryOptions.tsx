import { queryOptions} from "@tanstack/react-query";
import { fetchGame } from "./game";

export const gameQueryOptions = (guid: string) =>
    queryOptions({
        queryKey: ["game", { guid }],
        queryFn: () => fetchGame(guid),
    })
