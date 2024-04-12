CREATE TYPE person_type AS ENUM ('admin', 'regular', 'read-only');
CREATE TYPE animal_type AS ENUM ('rabbit', 'dog', 'cat');

CREATE TABLE person (
  "id" UUID PRIMARY KEY,
  "full_name" STRING NOT NULL,
  "date_of_birth" DATE NOT NULL,
  "user_type" person_type NOT NULL,
  "favourite_animal" animal_type NOT NULL
);