CREATE TABLE IF NOT EXISTS movies (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    origin VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    release_date DATE NOT NULL,
    status VARCHAR(255) NOT NULL,
    tags VARCHAR(255)[]
)