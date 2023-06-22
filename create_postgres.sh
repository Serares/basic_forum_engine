  docker run -p 5432:5432 -e POSTGRES_PASSWORD=$PG_PASSWORD -e POSTGRES_DB=$PG_DB -d --name=postgres postgres:14.5-alpine

  printf "\Wainting for postgres..."
  until docker exec postgres pg_isready > /dev/null 2> /dev/null; do
    printf "."
    sleep 5
  done

  echo "\nPostgres ready!"
