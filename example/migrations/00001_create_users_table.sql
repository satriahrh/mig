-- +mig Up
CREATE TABLE users (
    id int NOT NULL PRIMARY KEY,
    username text,
    name text,
    surname text
);

INSERT INTO users VALUES
(0, 'root', '', ''),
(1, 'sfreud', 'Sigmund', 'Freud');

-- +mig Down
DROP TABLE users;
