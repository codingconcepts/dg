IMPORT INTO "person"(
    "id",
    "full_name",
    "dob"
)
CSV DATA (
    'http://localhost:3000/person.csv'
)
WITH skip='1';

IMPORT INTO "account"(
    "person_id",
    "email"
)
CSV DATA (
    'http://localhost:3000/account.csv'
)
WITH skip='1';


IMPORT INTO "pet"(
    "id",
    "name",
    "person_id"
)
CSV DATA (
    'http://localhost:3000/pet.csv'
)
WITH skip='1';