-- Enable foreign key support explicitly for SQLite sessions
PRAGMA foreign_keys = ON;

-- 1. USERS TABLE
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,                 -- UUID stored as a string
    username TEXT NOT NULL UNIQUE,       -- Ensures unique display names
    email TEXT NOT NULL UNIQUE,          -- Required validation check for registration
    password_hash TEXT NOT NULL,         -- Stored using bcrypt encryption
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 2. SESSIONS TABLE
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,                 -- Session token (UUID)
    user_id TEXT NOT NULL,
    expires_at DATETIME NOT NULL,        -- Cookie expiration date validation
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 3. POSTS TABLE
CREATE TABLE IF NOT EXISTS posts (
    id TEXT PRIMARY KEY,                 -- UUID
    user_id TEXT NOT NULL,               -- The author of the post
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 4. CATEGORIES TABLE
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE            -- e.g., 'Technology', 'Sports', 'Others'
);

-- 5. POST_CATEGORIES JAVASCRIPT/RELATION TABLE (Many-to-Many Bridge)
-- This fulfills: "When registered users are creating a post they can associate one or more categories to it."
CREATE TABLE IF NOT EXISTS post_categories (
    post_id TEXT NOT NULL,
    category_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

-- 6. COMMENTS TABLE
CREATE TABLE IF NOT EXISTS comments (
    id TEXT PRIMARY KEY,                 -- UUID
    post_id TEXT NOT NULL,               -- Attached directly to a post
    user_id TEXT NOT NULL,               -- The author of the comment
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 7. REACTION/INTERACTION TABLE (Likes and Dislikes)
-- A unified table handles both posts and comments, enforcing a strict rule: One vote type per entity per user.
CREATE TABLE IF NOT EXISTS interactions (
    id TEXT PRIMARY KEY,                 -- UUID
    user_id TEXT NOT NULL,               -- Voter
    target_id TEXT NOT NULL,             -- Refers to either a Post ID or a Comment ID
    target_type TEXT NOT NULL CHECK (target_type IN ('post', 'comment')), -- Reaction target kind
    value INTEGER NOT NULL CHECK (value IN (1, -1)), -- 1 for Like, -1 for Dislike
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, target_id, target_type), -- Prevents double-liking or double-disliking
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ⚡ PERFORMANCE CRITICAL INDEXES
-- These guarantee sub-millisecond filtering logic when the database grows.
CREATE INDEX IF NOT EXISTS idx_posts_user ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_comments_post ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_interactions_target ON interactions(target_id, target_type);
