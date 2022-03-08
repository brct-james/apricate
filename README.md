# Apricate

Go-based server for a fantasy-themed capitalism simulator game set on a farm.

## Features

- Basic account functionality
- - Claim Account: `POST: https://apricate.io/api/users/{username}/claim`
- - - Don't forget to save the token from the response to your claim request. You must use this as a bearer token in the auth header for secure `/my/` routes
- - - Must include only letters, numbers, `-`, and `_`.
- - Get public user info at `/api/users/{username}` and get private user info including token at `/api/my/account`

---

### Endpoints

- `GET: /api/users` returns lists of registered usernames with various filters: unique, active, etc.
- `GET: /api/users/{username}` returns the public user data
- `POST: /api/users/{username}/claim` attempts to claim the specified username, returns the user data after creation, including token which users must save to access private routes
- `GET: /api/my/account` returns the private user data (includes token)

---

### Request Bodies

- None yet

---

### Response Codes

See `responses.go`

## Roadmap

Versioning Convention: `major.minor.hotfix`

---

### Ongoing

- All routes should use kebab-case
- All json & code should use snake_case

---

### In-Progress

**[v0.1]** MVP

- Basic routes function ("/" for project overview, "/api" for game stuff, "/docs" for auto generated documentation)
- Can register an account, which is stored in DB
- Auth middleware works
- Can GET user info, public and private versions

---

### Planned: v1 MVP

- TODO: This

---

### Planned: Post-1.0

- TODO: This

---

### Planned: Unscheduled

- Nothing yet

---

## Build & Run

Ensure resjon container is running on the correct port: `docker run -di -p 6382:6379 --name rejson_apricate redislabs/rejson:latest`

For the first run, ensure `refreshAuthSecret` in `main.go` is true. Make sure to set this to false for second run.

Build and start with `go build; ./apricate`. Alternatively, `go run main.go`

Listens on port `50250`

redis-cli via `redis-cli -p 6382`

`FLUSHDB` for each database (`select #`)

`KEYS *` to get all keys

`JSON.GET <token>` to get particular entry

Recommend running with screen `screen -S apricate`. If get detached, can forcibly detach the old ssh session and reattach with `screen -Dr apricate`

---

## Changelog

### v0.1

- Basic setup
- Add user claiming endpoint
- Add user info endpoint
- Add secure account info endpoint
- Add active users tracking

## Reference

### Technical

- https://github.com/RedisJSON/RedisJSON
- https://github.com/nitishm/go-rejson
- https://oss.redis.com/redisjson/commands/
- https://tutorialedge.net/golang/go-redis-tutorial/
- https://github.com/go-redis/redis
- https://tutorialedge.net/golang/parsing-json-with-golang/
- https://tutorialedge.net/golang/creating-restful-api-with-golang/
- https://github.com/joho/godotenv
- https://github.com/golang-jwt/jwt
- https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql

### Design

- https://api.spacetraders.io/
- https://spacetraders.io/docs/guide
- (Private) https://docs.google.com/document/d/15d-nC5dpiH19LH1sbWiUOM5Pjgr_Cjop-t_Dmuu2Xtc/edit
- (Private) https://keep.google.com/u/0/#LIST/1AyAhsCulc79U76hQK60tpjy9RaC5uQ6MdjHDYKDGrn8CsEPV56mWNezvrTPRdGA_cCrc9Q
- https://spacetraders.io/docs/ship-design