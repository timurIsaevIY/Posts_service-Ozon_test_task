CREATE TABLE IF NOT EXISTS posts
(
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    name TEXT NOT NULL,
    author TEXT NOT NULL,
    content TEXT NOT NULL,
    comments_allowed BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS comments (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    author TEXT NOT NULL,
    content TEXT NOT NULL,
    post_id INT NOT NULL,
    reply_to INT DEFAULT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (reply_to) REFERENCES comments(id) ON DELETE SET NULL
);

CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_reply_to ON comments(reply_to);
