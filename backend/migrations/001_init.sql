-- ============================================================
-- 001_init.sql — Initial schema for Lingdu Feed
-- ============================================================

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id              SERIAL          PRIMARY KEY,
    username        VARCHAR(50)     NOT NULL,
    password        VARCHAR(255)    NOT NULL,
    email           VARCHAR(100)    NOT NULL UNIQUE,
    following_count INT             NOT NULL DEFAULT 0,
    follower_count  INT             NOT NULL DEFAULT 0,
    created_time    TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- Posts table
CREATE TABLE IF NOT EXISTS posts (
    id               SERIAL          PRIMARY KEY,
    user_id          INT             NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title            VARCHAR(200)    NOT NULL,
    content          TEXT            NOT NULL,
    praise_count     INT             NOT NULL DEFAULT 0,
    comment_count    INT             NOT NULL DEFAULT 0,
    collection_count INT             NOT NULL DEFAULT 0,
    view_count       INT             NOT NULL DEFAULT 0,
    created_time     TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_time     TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_posts_user_id      ON posts(user_id);
CREATE INDEX idx_posts_created_time ON posts(created_time DESC);

-- Comments table
CREATE TABLE IF NOT EXISTS comments (
    id             SERIAL          PRIMARY KEY,
    post_id        INT             NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id        INT             NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reply_id       INT             REFERENCES comments(id) ON DELETE SET NULL,
    reply_username VARCHAR(50),
    content        TEXT            NOT NULL,
    created_time   TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_comments_post_id ON comments(post_id);

-- Praises (likes) table
CREATE TABLE IF NOT EXISTS praises (
    user_id      INT             NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id      INT             NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_time TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, post_id)
);

CREATE INDEX idx_praises_post_id ON praises(post_id);

-- Collections (bookmarks) table
CREATE TABLE IF NOT EXISTS collections (
    user_id      INT             NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id      INT             NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_time TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, post_id)
);

CREATE INDEX idx_collections_post_id ON collections(post_id);

-- Follows table
CREATE TABLE IF NOT EXISTS follows (
    follower_id  INT             NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id INT             NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_time TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id)
);

CREATE INDEX idx_follows_following_id ON follows(following_id);

-- Interaction status table (feed tracking)
CREATE TABLE IF NOT EXISTS interaction_status (
    post_id      INT             NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id      INT             NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status       SMALLINT        NOT NULL DEFAULT 0,   -- 0=unknown, 1=delivered, 2=displayed, 3=clicked
    updated_time TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    PRIMARY KEY (post_id, user_id)
);

CREATE INDEX idx_interaction_status_user_id ON interaction_status(user_id, status);
