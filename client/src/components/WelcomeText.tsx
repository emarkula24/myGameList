import "./WelcomeText.module.css"

export default function WelcomeText() {
    return (
        <div style={{width: "75%"}}>
            <p><span style={{color: "green", fontWeight: "bold"}}>myGameList </span> 
                is a platform built for gaming enthusiasts.
                Here you can track and manage your game backlog and even compare it with your friends!
                Please login to get started on creating your own gamelist.
            </p>
            <p>You can copy and share the link to your list with others. If you are logged in, you can also add items from someone else's list to your own if you don't already have that item added.</p>
        </div>
    )
}