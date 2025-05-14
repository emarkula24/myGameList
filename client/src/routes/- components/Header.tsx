import HeaderElement from "./HeaderElement"
import { Location, Game } from "../- types/types"
import styles from "./Header.module.css"
import SearchBar from "./SearchBar"
import { useState } from "react"



// Header component, consists of HeaderElements and searchbar
export default function Header() {
        const [searchResults, setSearchResults] = useState<Game[]>([])
        const locations: Location[] = [
                { address: "/", addressName: "Home" },
                { address: "/about", addressName: "About" },
                { address: "/community", addressName: "Community" },
                { address: "/profile", addressName: "Profile" },
                { address: "/login", addressName: "Login" },
                { address: "/register", addressName: "Sign up" },

        ]
        console.log(searchResults)
        return (
                <div className={styles.headerContainer}>
                        {locations.map((location) => (
                                < HeaderElement key={location.address} location={location} />
                        ))}
                        <SearchBar setSearchResults={setSearchResults}/>
                </div>
        )
}