db:
	cockroach demo --insecure --no-example-database

tables:
	cockroach sql --insecure < examples/many_to_many/create.sql

data:
	go run dg.go -c ./examples/many_to_many/config.yaml -o ./csvs/many_to_many

file_server:
	python3 -m http.server 3000 -d csvs/many_to_many

import:
	cockroach sql --insecure < examples/many_to_many/insert.sql