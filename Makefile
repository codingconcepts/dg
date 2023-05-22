db:
	cockroach demo --insecure --no-example-database

tables:
	cockroach sql --insecure < examples/person_pet_create.sql

data:
	go run dg.go -c examples/person_pet.yaml -o csvs/person_pet

file_server:
	python3 -m http.server 3000 -d csvs/person_pet

import:
	cockroach sql --insecure < examples/person_pet_insert.sql