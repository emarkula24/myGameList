import { useRouter, useRouterState } from "@tanstack/react-router"
import { useAuth } from "../utils/auth"
import React from "react";


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
        } catch (error: any) {
            setError(error.message || 'Logout failed');
        } finally {
            setIsSubmitting(false)
        }
        
    }
    
    return (
        <div>
            <div onClick={handleClick}>Logout</div>
        </div>
    )
}