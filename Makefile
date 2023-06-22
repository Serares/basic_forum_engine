build:
	@go build -o ./bin/basic_forum_engine

create_db:
	./create_postgres.sh

run: build create_db
	./bin/basic_forum_engine
