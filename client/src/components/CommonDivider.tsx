import styles from "./CommonDivider.module.css"
export default function CommonDivider({routeName}: {routeName: string} ) {
    return (
        <div className={styles.divider}>
            {routeName}
        </div>
    )
}