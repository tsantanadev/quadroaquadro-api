CREATE TABLE IF NOT EXISTS movies (
    id INTEGER PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    origin VARCHAR(255)[] NOT NULL,
    tags VARCHAR(255)[],
    release_date DATE NOT NULL,
    original_title VARCHAR(255) NOT NULL,
    poster_path VARCHAR(255) NOT NULL,
    genres VARCHAR(255)[]
)