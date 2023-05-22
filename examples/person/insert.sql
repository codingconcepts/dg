IMPORT INTO "person"(
    "uuid",
    "string",
    "date",
    "bool",
    "int8",
    "int16",
    "int32",
    "int64"
)
CSV DATA (
    'http://localhost:3000/person.csv'
)
WITH 
    skip='1',
    nullif = '',
    allow_quoted_null;