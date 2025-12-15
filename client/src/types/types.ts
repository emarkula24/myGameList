export interface Location {
        address: string;
        addressName: string;
}

export interface HeaderElementProps {
        location: Location;
}
interface Cover {
        id: number;
        url: string ;
}
interface Platform {
        id: number;
        abbreviation: string;
        platform_logo: PlatformLogo;
}
export interface SearchedGame {
        id: number;
        name: string;
        cover: Cover;
        platforms: Platform[];
}

interface Genre {
  id: number;
  name: string;
}

interface PlatformLogo {
  id: number;
  url: string;
}


interface ReleaseDate {
  id: number;
  human: string;
}

export interface Game {
  id: number;
  name: string;
  genres: Genre[];
  platforms: Platform[];
  rating: number;
  rating_count: number;
  release_dates: ReleaseDate[];
  similar_games: number[];
  cover: Cover;
  summary: string
}

export interface GameListEntries {
  id: number;
  name: string;
  cover: string;
  release_dates: ReleaseDate[];
  status: number;
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
