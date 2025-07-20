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
        image?: {
                medium_url: string;
                icon_url: string;
                tiny_url: string;
                thumb_url: string;
                small_url: string;
                super_url: string;
        };
        original_release_date?: string;
        platforms?: { abbreviation: string }[];
        deck?: string;
        site_detail_url: string;
}

export interface Game {
        id: number;
        name: string;
        deck: string;
        description: string;
        image: {
                medium_url: string;
        }
        original_release_data: string;
        platforms: Array<object>;
        publishers: Array<object>;
        similar_games: Array<object>;
        genres: Array<object>;

        
}
export interface GameListEntry {
        status: string;
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
        platforms: Array<object>;
        publishers: Array<object>;
        similar_games: Array<object>;
        genres: Array<object>;

}
export interface LoginResponse {
        accessToken: string;
        userId: number;
}

export interface RegisterResponse {
        user_id: number;
}

export type SubmitErrorProps = {
  err: string | null;
};
