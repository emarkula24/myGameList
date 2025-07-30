export function GameListStatusHeader({ statusText }: { statusText: string }) {
    return (
        <div style={{
            width: '100%',
            height: '40px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            backgroundColor: 'green',
            border: "1px solid grey"
        }}>
            <span style={{ fontWeight: "600", color: 'aliceblue', fontSize: "2.4em" }}>{statusText}</span>
        </div>
    );
}