CREATE TABLE person (
  id UUID PRIMARY KEY,
  full_name STRING NOT NULL,
  dob DATE NOT NULL
);

CREATE TABLE account (
  person_id UUID PRIMARY KEY REFERENCES person(id),
  email STRING NOT NULL
);

CREATE TABLE pet (
  id UUID PRIMARY KEY,
  name STRING NOT NULL,
  person_id UUID NOT NULL REFERENCES person(id)
);