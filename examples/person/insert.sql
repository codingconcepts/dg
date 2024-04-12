IMPORT INTO "person"(
    "id",
    "full_name",
    "date_of_birth",
    "user_type",
    "favourite_animal"
)
CSV DATA (
    'http://localhost:3000/person.csv'
)
WITH 
    skip='1',
    nullif = '',
    allow_quoted_null;