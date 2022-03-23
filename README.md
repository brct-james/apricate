# Apricate

Go-based server for a fantasy-themed capitalism simulator game set on a farm.

### Documentation

Please visit [the official docs](https://apricate.stoplight.io/docs/apricate/ZG9jOjQ3MDIzNTgw-alpha-guide). Contribute changes in the [documentation repo](https://github.com/brct-james/apricate-docs).

---

### Response Codes

See `responses.go`

## Standards

Versioning Convention: `major.minor.hotfix`
Routes should use kebab-case
Json and code should use snake_case (go uses camelCase?)

---

## Build & Run

Modify the volumes to your local environment in the docker-compose file you want to use, then run the appropriate `run_dev.sh` / `start_live.sh` script.

For the first run, ensure `refreshAuthSecret` in `main.go` is true. Make sure to set this to false for second run

DEV Listens on port `50520`
LIVE Listens on port `50250`

redis-cli via `redis-cli -p 6382` for DEV
redis-cli via `redis-cli -p 6383` for LIVE

`FLUSHDB` for each database (`select #`)

`KEYS *` to get all keys

`JSON.GET <token>` to get particular entry

---

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
