export interface Location {
        address: string;
        addressName: string;
}

export interface HeaderElementProps {
        location: Location;
}

export interface Games {
        id: number;
        guid: string;
        name: string;
        image: {
                medium_url: string;
                icon_url: string;
                tiny_url: string;
                thumb_url: string;
                small_url: string;
                super_url: string;
        };
        original_release_date: string;
        platforms: { name: string }[];
        deck: string;
        site_detail_url: string;
}

export interface Game {
        id: number;
        name: string;
        guid: string;
        status: number;
        deck: string;
        description: string;
        image: {
                medium_url: string;
                icon_url: string;
                tiny_url: string;
                thumb_url: string;
                small_url: string;
        }
        original_release_data: string;
        platforms: { name: string }[];
        publishers: object[];
        similar_games: object[];
        genres: object[];


}
export interface GameListEntries {
        status: number;
        id: number;
        guid: string;
        name: string;
        deck: string;
        description: string;
        image: {
                icon_url: string;
                medium_url: string;
                thumb_url: string;
                tiny_url: string;
                small_url: string;
        }
        original_release_data: string;
        platforms: { name: string }[];
        publishers: object[];
        similar_games: object[];
        genres: object[];

}
export interface GameList {
        gamelist: GameListEntries
}
export interface GameData {
        id?: number;
        status?: number;
}
export interface GameListEntry {
        gamedata: GameData | null
}
export interface RegisterResponse {
        user_id: number;
}

export interface SubmitErrorProps {
        err: string | null;
}
