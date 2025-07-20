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
);

CREATE TABLE IF NOT EXISTS games (
    game_id INT UNSIGNED NOT NULL,
    gamename VARCHAR(255)
    PRIMARY KEY (game_id)
);

CREATE TABLE IF NOT EXISTS user_games (
    username VARCHAR(80) NOT NULL,
    game_id INT UNSIGNED NOT NULL,
    status VARCHAR(120) NOT NULL,
    PRIMARY KEY (username, game_id),
    FOREIGN KEY (username) REFERENCES users(username),
    FOREIGN KEY (game_id) REFERENCES games(game_id)
);

INSERT INTO users (username, email, password, created_at)
VALUES ('listaddtest', 'list@example.com', 'hashedpassword123', NOW());