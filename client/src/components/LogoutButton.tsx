import { useRouter, useRouterState } from "@tanstack/react-router"
import { useAuth } from "../utils/auth"
import React from "react";
import styles from "./LogoutButton.module.css"
import { Spinner } from "./Spinner";

export default function LogoutButton() {
    const [isSubmitting, setIsSubmitting] = React.useState(false)
    const isLoading = useRouterState({ select: (s) => s.isLoading })
    const [error, setError] = React.useState<string | null>(null);
    const auth = useAuth()
    const router = useRouter()
    if (!auth.user) {
        console.log("no user to log out")
        return
    }
    const handleClick = async () => {
        setIsSubmitting(true)
        try {
            await auth.logout(auth.user?.username, auth.user?.userId)
            await router.invalidate()
            await router.navigate({to: "/"})
        } catch (error: unknown) {
            if (error instanceof Error) {
                setError(error.message || 'Logout failed');
            }
            
        } finally {
            setIsSubmitting(false)
        }
        
    }
    
    return (
        <div>
            <div onClick={() => {void handleClick()}} className={styles.btn}>Logout</div>
            {isSubmitting && <span></span>}
            {isLoading && < Spinner />}
            {error && <div>{error}</div>}
        </div>
        
    )
}