validate_version:
ifndef VERSION
	$(error VERSION is undefined)
endif

db:
	cockroach demo --insecure --no-example-database

tables:
	cockroach sql --insecure < examples/many_to_many/create.sql

data_many_to_many:
	go run dg.go -c ./examples/many_to_many/config.yaml -o ./csvs/many_to_many -i import.sql

data_person:
	go run dg.go -c ./examples/person/config.yaml -o ./csvs/person

data_range_test:
	go run dg.go -c ./examples/range_test/config.yaml -o ./csvs/range_test

data_input_test:
	go run dg.go -c ./examples/input_test/config.yaml -o ./csvs/input_test

data_unique_test:
	go run dg.go -c ./examples/unique_test/config.yaml -o ./csvs/unique_test

data_const_test:
	go run dg.go -c ./examples/const_test/config.yaml -o ./csvs/const_test

data_match:
	go run dg.go -c ./examples/match_test/config.yaml -o ./csvs/match -i import.sql

data_each_match:
	go run dg.go -c ./examples/each_match_test/config.yaml -o ./csvs/each_match -i import.sql

data_pattern:
	go run dg.go -c ./examples/pattern_test/config.yaml -o ./csvs/pattern_test -i import.sql

data_cuid2:
	go run dg.go -c ./examples/cuid2_test/config.yaml -o ./csvs/cuid2_test -i import.sql

data_template:
	go run dg.go -c ./examples/gen_templates_test/config.yaml -o ./csvs/gen_templates_test -i import.sql

data_rel_date:
	go run dg.go -c ./examples/rel_date_test/config.yaml -o ./csvs/rel_date_test -i import.sql

data_rand:
	go run dg.go -c ./examples/rand_test/config.yaml -o ./csvs/rand_test -i import.sql

data_expr:
	go run dg.go -c ./examples/expr_test/config.yaml -o ./csvs/expr_test -i import.sql

data: data_many_to_many data_person data_range_test data_input_test data_unique_test data_const_test \
	data_match data_each_match data_pattern data_cuid2 data_template data_rel_date data_rand data_expr
	echo "done"

file_server:
	python3 -m http.server 3000 -d csvs/many_to_many

import:
	cockroach sql --insecure < examples/many_to_many/insert.sql

test:
	go test ./... -v -cover

cover:
	go test -v -coverpkg=./... -coverprofile=profile.cov ./... -count=1
	go tool cover -func profile.cov
	# go tool cover -html coverage.out

profile:
	go run dg.go -c ./examples/many_to_many/config.yaml -o ./csvs/many_to_many -cpuprofile profile.out
	go tool pprof -http=:8080 profile.out

release: validate_version
	# linux
	GOOS=linux go build -ldflags "-X main.version=${VERSION}" -o dg ;\
	tar -zcvf ./releases/dg_${VERSION}_linux.tar.gz ./dg ;\

	# macos (arm)
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=${VERSION}" -o dg ;\
	tar -zcvf ./releases/dg_${VERSION}_macos_arm64.tar.gz ./dg ;\

	# macos (amd)
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o dg ;\
	tar -zcvf ./releases/dg_${VERSION}_macos_amd64.tar.gz ./dg ;\

	# windows
	GOOS=windows go build -ldflags "-X main.version=${VERSION}" -o dg ;\
	tar -zcvf ./releases/dg_${VERSION}_windows.tar.gz ./dg ;\

	rm ./dg