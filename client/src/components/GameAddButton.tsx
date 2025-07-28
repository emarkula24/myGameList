
export interface GameAddButtonProps {
    onNewListEntry: (status: number) => void
}


export default function GameAddButton({ onNewListEntry }: GameAddButtonProps) {

    const handleClick = () => {
        onNewListEntry(1)
    }
    
    return (
        <div>
            <button onClick={handleClick} type="button">Add to GameList</button>
        </div>
    )
}