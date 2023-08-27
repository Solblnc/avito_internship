CREATE TABLE IF NOT EXISTS segments (
    segment_id serial PRIMARY KEY,
    segment_name VARCHAR(255) NOT NULL
);



CREATE TABLE IF NOT EXISTS users (
    user_id serial PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS segments_user (
  user_id INT REFERENCES users(user_id),
  segment_id INT REFERENCES segments(segment_id),
  time timestamp NOT NULL
);

---- create above / drop below ----
DROP TABLE segments;
DROP TABLE users;
DROP TABLE segments_user;
