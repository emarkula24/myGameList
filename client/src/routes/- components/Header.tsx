import HeaderElement from "./HeaderElement"
import { Location } from "../- types/types"
import styles from "./Header.module.css"
import SearchBar from "./SearchBar"
import {
    QueryClient,
    QueryClientProvider,

} from '@tanstack/react-query'

const queryClient = new QueryClient()

// Header component, consists of HeaderElements and searchbar
export default function Header() {
    const locations: Location[] = [
        { address: "/", addressName: "Home" },
        { address: "/about", addressName: "About" },
        { address: "/community", addressName: "Community" },
        { address: "/profile", addressName: "Profile" },
        { address: "/login", addressName: "Login" },
        { address: "/register", addressName: "Sign up" },

    ]

    return (
        <div className={styles.headerContainer}>
            {locations.map((location) => (
                < HeaderElement key={location.address} location={location} />
            ))}
            <QueryClientProvider client={queryClient}>
                <SearchBar />
            </QueryClientProvider>
        </div>
    )
}