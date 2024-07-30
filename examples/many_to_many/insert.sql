IMPORT INTO "person" (
    "id"
)
CSV DATA (
    'http://localhost:3000/person.csv'
)
WITH 
    skip='1',
    nullif = '',
    allow_quoted_null;

IMPORT INTO "event" (
    "id"
)
CSV DATA (
    'http://localhost:3000/event.csv'
)
WITH 
    skip='1',
    nullif = '',
    allow_quoted_null;

IMPORT INTO "person_type" (
    "id",
    "name"
)
CSV DATA (
    'http://localhost:3000/person_type.csv'
)
WITH 
    skip='1',
    nullif = '',
    allow_quoted_null;

IMPORT INTO "person_event" (
    "person_id",
    "event_id",
    "id",
    "person_type"
)
CSV DATA (
    'http://localhost:3000/person_event.csv'
)
WITH 
    skip='1',
    nullif = '',
    allow_quoted_null;
