CREATE TABLE images (
    id UUID PRIMARY KEY,
    movie_id BIGINT NOT NULL,
    level INT NOT NULL,
    CONSTRAINT fk_movie FOREIGN KEY (movie_id) REFERENCES movies(id)
);

