# Run go build first
# go build; docker compose -f ./docker-compose-dev.yml -p apricate-dev up -d --build
# Run attached
docker compose -f ./docker-compose-dev.yml -p apricate-dev down; go build; docker compose -f ./docker-compose-dev.yml -p apricate-dev up -d --build