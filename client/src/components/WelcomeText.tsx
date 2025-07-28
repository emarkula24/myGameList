import "./WelcomeText.module.css"

export default function WelcomeText() {
    return (
        <div style={{width: "75%"}}>
            <p><span style={{color: "green", fontWeight: "bold"}}>myGameList </span> 
                is a platform built for gaming enthusiasts.
                Here you can track and manage your game backlog and even compare it with your friends!
                Please login to get started on creating your own gamelist.
            </p>
        </div>
    )
}