CREATE TABLE IF NOT EXISTS segments (
    slug VARCHAR(100) PRIMARY KEY 
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    encrypted_pwd VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS segments_to_users (
    segment_slug VARCHAR(100) REFERENCES segments(slug) ON DELETE CASCADE NOT NULL,
    user_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    expiration_date TIMESTAMP WITHOUT TIME ZONE,
    PRIMARY KEY (segment_slug, user_id)
);