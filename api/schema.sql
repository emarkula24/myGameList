-- Init file for testcontainer used in integration tests
CREATE TABLE IF NOT EXISTS users (
    user_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    username VARCHAR(80) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password CHAR(60) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id)
);

CREATE TABLE IF NOT EXISTS refreshtokens (
    user_id INT UNSIGNED NOT NULL,
    jti CHAR(36) NOT NULL UNIQUE,
    refresh_token VARCHAR(512) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
)
