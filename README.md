# Apricate

Go-based server for a fantasy-themed capitalism simulator game set on a farm.

## Features

- Basic account functionality
- - Claim Account: `POST: https://apricate.io/api/users/{username}/claim`
- - - Don't forget to save the token from the response to your claim request. You must use this as a bearer token in the auth header for secure `/my/` routes
- - - Usernames must include only letters, numbers, `-`, and `_`. There are some reserved sequences for game-specific prefixes. A slur filter exists, I hope never to need it.
- - Get public user info at `/api/users/{username}` and get private user info including token at `/api/my/account`
- Monitor list of assistants or specific assistant with `/api/my/assistants` and `/api/my/assistants/{uuid}`
- Islands which contain an X-Y grid of Locations from -100 to 100, with each Island separated by sailing lanes connected to Ports in each island
- Locations are farms or towns, and may hold NPCs or a market.
---

### Endpoints

**Public Routes**
<!-- - `GET: /api/v0/leaderboards` list all available leaderboards and their descriptions
- `GET: /api/v0/leaderboards/{board}` get the specified leaderboard rankings -->
- `GET: /api/islands` returns details on every island in the game, including port connections, for navigational purposes
- `GET: /api/users` returns lists of registered usernames with various filters: unique, active, etc.
- `GET: /api/users/{username}` returns the public user data
- `POST: /api/users/{username}/claim` attempts to claim the specified username, returns the user data after creation, including token which users must save to access private routes
- `GET: /api/plants` returns the data on every plant in the game
- `GET: /api/plants/{plantName}` returns the data on the specified plant

**Secure Routes**
- `GET: /api/my/account` returns the private user data (includes token)
- `GET: /api/my/assistants` returns a list of the player's assistants
- `GET: /api/my/assistants/{uuid}` returns the assistant specified by `uuid`
- `GET: /api/my/farms` returns a list of the player's farms
- `GET: /api/my/farms/{uuid}` returns the farm specified by `uuid`
- `GET: /api/my/plots` returns a list of the player's plots
- `GET: /api/my/plots/{uuid}` returns the plot specified by `uuid`
- `GET: /api/my/contracts` returns a list of the player's contracts
- `GET: /api/my/contracts/{uuid}` returns the contract specified by `uuid`
- `GET: /api/my/warehouses` returns a list of the player's warehouses
- `GET: /api/my/warehouses/{uuid}` returns the warehouse specified by `uuid`
- `GET: /api/my/locations` returns the details for any location with an assistant as well as any owned farms
- `GET: /api/my/locations/{name}` returns the details of the location specified by `name` IF the location is an owned farm or holds an assistant
- `GET: /api/my/nearby-locations` returns a list of the names of every nearby location (all locations of every island with atleast one assistant), for navigational purposes

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

### Complete: **[v0.2]**

- ~~Assistants have helper functions and initialize on create, have GET endpoints~~
- ~~Farms have helper functions and initialize on create, have GET endpoints~~
- ~~Contracts have helper functions and initialize on create, have GET endpoints~~
- ~~Warehouses have helper functions and initialize on create, have GET endpoints~~
- ~~Pria and Veldis YAML defined~~

---

### In-Progress: **[v0.3]**

- ~~PlantDefinition YAML defined for at least 3 types of plants excluding Wild Seeds~~
- ~~PlantDefinition information public GET endpoints~~
- ~~Plot helper functions and initialize on create~~
- ~~Plot GET endpoints~~
- ~~Plant struct for plots defined~~
- Plot helper functions for Growth Actions (plant plot, clear plot, progress plot, harvest plot)
- Plot interact endpoint with switch on body.action (growth actions)
- - GA_Wait
- - GA_Clear
- - GA_Trim
- - GA_Dig
- - GA_Weed
- - GA_Fertilize
- - GA_Water
- - GA_Hill
- - GA_Sprout
- Plants grow when planted on plots
- Goods are deposited to warehouse when plants harvested
- Tested at least 3 types of plants excluding Wild Seeds
- ~~Skelling and Tritum YAML defined~~
- Convert from UUID to composite string like "username|farmid|plotid" e.g. "Greenitthe|Homestead Farm|Plot-1" for warehouses, farms, plots, assistants, contracts. UUID wont be used
- ~~Convert to symbol based sector/island/location~~
- ~~Convert regions to islands~~
- ~~Add sectors~~

---

### Planned: **[v0.4]**

- Assistants can transfer things between warehouses
- Boldor, Yoggoth YAML defined

---

### Planned: **[v0.5]**

- One functional market, Local orders (non-player orders to provide baseline supply/demand)
- Ratelimiting
- Add `GET` endpoints for sectors, `GET` select island endpoint

---

### Planned: **[v0.6]**

- Add YAML-defined contract/quest paths from NPCs
- NPCs defined in YAML

---

### Planned: **[v0.7]**

- At least 10 different plants excluding Wild Seeds, buy and sell at the markets of at least 3 towns, all starting town NPCs have quests
- Tyldia YAML defined
- Wild Seeds implemented

---

### Planned: **[v0.8]**

- At least 20 plants excluding Wild Seeds, add at least 2 additional tools for growing some of the new plants, add randomized contracts, consider adding additional markets, all NPCs on starting map have quests

---

### Planned: **[v0.9]**

- Add refining/crafting with at least 8 recipes, add at least 4 new tools to support crafting, add at least 2 new buildings to support crafting
- Add researching plants (with associated building) to reveal full information
- Fog of War and hide unresearched plant information

---

### Planned: **[v1.0]**

- Meta account and progression, leaderboards, full documentation, separate dev partition that won't affect live

---

### Planned: Post-1.0

- Simplified routing (pass a full route, even over oceans, and server will calculate total fare and travel time, rather than requiring manual travel between each intermediate location) - meta progression unlock?

---

### Planned: Unscheduled

- Nothing yet

---

## Build & Run

Ensure resjon container is running on the correct port: `docker run -di -p 6382:6379 --name rejson_apricate redislabs/rejson:latest`

For the first run, ensure `refreshAuthSecret` in `main.go` is true. Make sure to set this to false for second run.

Either run once to generate, or manually create `slur_filter.txt` in root directory. Add words to filter, one per line, case-insensitive.

Build and start with `go build; ./apricate`. Alternatively, `go run main.go`

Listens on port `50250`

redis-cli via `redis-cli -p 6382`

`FLUSHDB` for each database (`select #`)

`KEYS *` to get all keys

`JSON.GET <token>` to get particular entry

Recommend running with screen `screen -S apricate`. If get detached, can forcibly detach the old ssh session and reattach with `screen -Dr apricate`

---

## Changelog

### v0.2

- Add assitants GET
- Add farms GET
- Add warehouses GET
- Add contracts GET

### v0.3

- Add skellig and tritum YAML
- Add plants YAMl
- Add/modify schemas to support plots/plants
- Add ActionToSkip field to Optional GrowthStages for better UX

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
