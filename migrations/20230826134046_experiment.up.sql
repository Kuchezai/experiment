CREATE TABLE IF NOT EXISTS segments (
    slug VARCHAR(100) PRIMARY KEY 
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    encrypted_pwd VARCHAR(100)  
);