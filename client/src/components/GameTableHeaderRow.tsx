import styles from "./GameTableHeaderRow.module.css"

export default function GameTableHeaderRow () {
    return (
        <>
            <tr style={{ backgroundColor: "lightgray" }}>
                {/* <th className={styles.headerTitleStatus}></th> */}
                <th className={styles.headerNumber}>#</th>
                <th className={styles.headerImage}>Image</th>
                <th className={styles.headerName}>Game Title</th>
                <th className={styles.headerStatus}>Status</th>
                <th className={styles.headerAction}>Actions</th>
            </tr>
        </>
    )
}