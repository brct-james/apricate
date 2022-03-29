# Run go build first
docker compose -f ./docker-compose-live.yml -p apricate-live down; go build -o apricate-live; docker compose -f ./docker-compose-live.yml -p apricate-live up -d --build