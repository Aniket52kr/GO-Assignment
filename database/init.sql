-- Main user accounts
CREATE TABLE IF NOT EXISTS t_users (
    email       VARCHAR(320)    UNIQUE NOT NULL,
    username    VARCHAR(32)     PRIMARY KEY,
    password    VARCHAR(64)     NOT NULL,
    id          CHAR(36)        UNIQUE NOT NULL,
    verified    BOOLEAN         NOT NULL,
    avatar      TEXT,
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;



-- Special users like owners
CREATE TABLE IF NOT EXISTS o_users (
    id          CHAR(36)        PRIMARY KEY,
    CONSTRAINT fk_id
        FOREIGN KEY(id)
            REFERENCES t_users(id)
            ON DELETE CASCADE
) ENGINE=InnoDB;



-- URL shortener feature (eg. posts, profiles):-
CREATE TABLE IF NOT EXISTS shorturl (
    token       VARCHAR(320)    PRIMARY KEY,
    id          CHAR(36)        UNIQUE NOT NULL
) ENGINE=InnoDB;



--  Stores user-created posts
CREATE TABLE IF NOT EXISTS posts (
    user_id     CHAR(36)        NOT NULL,
    id          CHAR(36)        PRIMARY KEY,
    body        VARCHAR(320)    NOT NULL,
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_id
        FOREIGN KEY(user_id)
            REFERENCES t_users(id)
            ON DELETE CASCADE
) ENGINE=InnoDB;



--  Tracks "follow" relationships between users
CREATE TABLE IF NOT EXISTS follows (
    user_id     CHAR(36)        NOT NULL,
    follow_id   CHAR(36)        NOT NULL,
    PRIMARY KEY (user_id, follow_id),
    CONSTRAINT fk_user_id
        FOREIGN KEY(user_id)
            REFERENCES t_users(id)
            ON DELETE CASCADE,
    CONSTRAINT fk_follow_id
        FOREIGN KEY(follow_id)
            REFERENCES t_users(id)
            ON DELETE CASCADE
) ENGINE=InnoDB;




-- Tracks post upvotes (or likes)
CREATE TABLE IF NOT EXISTS votes (
    user_id     CHAR(36)        NOT NULL,
    id          CHAR(36)        NOT NULL,
    PRIMARY KEY (user_id, id),
    CONSTRAINT fk_id
        FOREIGN KEY(id)
            REFERENCES posts(id)
            ON DELETE CASCADE,
    CONSTRAINT fk_user_id
        FOREIGN KEY(user_id)
            REFERENCES t_users(id)
            ON DELETE CASCADE
) ENGINE=InnoDB;



-- Stores comments on posts
CREATE TABLE IF NOT EXISTS comments (
    user_id     CHAR(36)        NOT NULL,
    post_id     CHAR(36)        NOT NULL,
    id          CHAR(36)        PRIMARY KEY,
    body        VARCHAR(320)    NOT NULL,
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_post_id
        FOREIGN KEY(post_id)
            REFERENCES posts(id)
            ON DELETE CASCADE,
    CONSTRAINT fk_user_id
        FOREIGN KEY(user_id)
            REFERENCES t_users(id)
            ON DELETE CASCADE
) ENGINE=InnoDB;
