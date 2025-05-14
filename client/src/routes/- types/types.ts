export interface Location {
        address: string;
        addressName: string;
}

export interface HeaderElementProps {
        location: Location;
}

export interface Game {
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