import { Link } from "@tanstack/react-router"
import { HeaderElementProps } from "../- types/types"
import styles from "./HeaderElement.module.css";

// Component for elements inside the header
export default function HeaderElement({ location }: HeaderElementProps) {
    return (
        <>
            <Link to={location.address} className={ styles.linkButton}>
                {location.addressName}
            </Link>
        </>
    )
}