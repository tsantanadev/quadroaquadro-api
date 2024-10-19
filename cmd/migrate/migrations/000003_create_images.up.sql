CREATE TABLE images (
    id BIGSERIAL PRIMARY KEY,
    movie_id BIGINT NOT NULL,
    url VARCHAR(255) NOT NULL,
    level INT NOT NULL,
    CONSTRAINT fk_movie FOREIGN KEY (movie_id) REFERENCES movies(id)
);

