import styles from "./GameTableHeaderRow.module.css"

export default function GameTableHeaderRow () {
    return (
        <>
            <tr style={{ backgroundColor: "rgb(232, 241, 232)", height: "30px", boxSizing: "border-box"}}>
                <th className={styles.headerTitleStatus}></th>
                <th className={`${styles.headerNumber} ${styles.shortRightBorder}`}>#</th>
                <th className={`${styles.headerImage} ${styles.shortRightBorder}`}>Image</th>
                <th className={`${styles.headerName} ${styles.shortRightBorder}`}>Game Title</th>
                <th className={`${styles.headerStatus} ${styles.shortRightBorder}`}>Status</th>
                <th className={styles.headerAction}>Actions</th>
            </tr>
        </>
    )
}