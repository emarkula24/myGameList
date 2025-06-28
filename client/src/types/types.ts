export interface Location {
        address: string;
        addressName: string;
}

export interface HeaderElementProps {
        location: Location;
}

export interface Games {
        id: number;
        guid: number;
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
        name: string;
        deck: string;
        description: string;
        image: {
                medium_url: string;
        }
        original_release_data: string;
        platforms: Array<object>;
        publishers: Array<object>;
        
}
export interface LoginResponse {
        accessToken: string;
        userId: number;
}