CREATE TABLE "person" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid()
);

CREATE TABLE "event" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid()
);

CREATE TABLE "person_type" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "name" STRING NOT NULL
);

CREATE TABLE "person_event" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "person_type" UUID NOT NULL REFERENCES "person_type"("id"),
  "person_id" UUID NOT NULL REFERENCES "person"("id"),
  "event_id" UUID NOT NULL REFERENCES "event"("id")
);