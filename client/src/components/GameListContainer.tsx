export default function GameListContainer({ children }: { children: React.ReactNode }) {
  return (
    <div style={{
      display: "flex",
      flexDirection: "column",
      alignItems: "center",
      justifyContent: "center",
      border: "1px solid lightgrey",
      boxSizing: "border-box",
      width: "75%"
    }}>
      {children}
    </div>
  );
}
