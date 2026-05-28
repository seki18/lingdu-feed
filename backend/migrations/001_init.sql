-- ============================================================
-- 001_init.sql — Initial schema for Lingdu Feed (v2)
-- Uses final table names: likes, favorites, states
-- Uses final column names: like_count, favorite_count
-- ============================================================

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id              SERIAL          PRIMARY KEY,
    username        VARCHAR(50)     NOT NULL UNIQUE,
    password        VARCHAR(255)    NOT NULL,
    email           VARCHAR(100)    NOT NULL UNIQUE,
    following_count INT             NOT NULL DEFAULT 0,
    follower_count  INT             NOT NULL DEFAULT 0,
    created_time    TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- Posts table (content metadata only; stats are in post_stats)
CREATE TABLE IF NOT EXISTS posts (
    id               SERIAL          PRIMARY KEY,
    user_id          INT             NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title            VARCHAR(200)    NOT NULL,
    content          TEXT            NOT NULL,
    created_time     TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_time     TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_created_time ON posts(created_time DESC);
CREATE INDEX idx_posts_user_created_time ON posts(user_id, created_time DESC);

-- Post stats table (1:1 with posts, high-write counters separated for caching)
CREATE TABLE IF NOT EXISTS post_stats (
    id               INT              PRIMARY KEY REFERENCES posts(id) ON DELETE CASCADE,
    like_count       INT              NOT NULL DEFAULT 0,
    comment_count    INT              NOT NULL DEFAULT 0,
    favorite_count   INT              NOT NULL DEFAULT 0,
    view_count       INT              NOT NULL DEFAULT 0,
    expose_count     INT              NOT NULL DEFAULT 0,
    score            DOUBLE PRECISION NOT NULL DEFAULT 0,
    updated_time     TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_post_stats_score ON post_stats(score DESC);

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
CREATE INDEX idx_comments_post_created_time ON comments(post_id, created_time ASC);
CREATE INDEX idx_comments_reply_id ON comments(reply_id);

-- Likes table
CREATE TABLE IF NOT EXISTS likes (
    user_id      INT             NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id      INT             NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_time TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, post_id)
);

CREATE INDEX idx_likes_post_id ON likes(post_id);

-- Favorites table
CREATE TABLE IF NOT EXISTS favorites (
    user_id      INT             NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id      INT             NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_time TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, post_id)
);

CREATE INDEX idx_favorites_post_id ON favorites(post_id);

-- Follows table
CREATE TABLE IF NOT EXISTS follows (
    follower_id  INT             NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id INT             NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_time TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id),
    CHECK (follower_id <> following_id)
);

CREATE INDEX idx_follows_following_id ON follows(following_id);

-- States table (feed tracking pipeline: 0=unknown, 1=delivered, 2=exposed, 3=clicked)
CREATE TABLE IF NOT EXISTS states (
    user_id      INT             NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id      INT             NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    status       SMALLINT        NOT NULL DEFAULT 0,
    updated_time TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, post_id),
    CHECK (status BETWEEN 0 AND 3)
);

CREATE INDEX idx_states_user_status ON states(user_id, status);

-- Post images table (supports multiple images per post with ordering)
CREATE TABLE post_images (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    sort_order INT DEFAULT 0,
    created_time TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_post_images_post_id ON post_images(post_id);
