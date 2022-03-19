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
- Plant plants in farm plots

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
- `GET: /api/my/contracts` returns a list of the player's contracts
- `GET: /api/my/contracts/{uuid}` returns the contract specified by `uuid`
- `GET: /api/my/warehouses` returns a list of the player's warehouses
- `GET: /api/my/warehouses/{uuid}` returns the warehouse specified by `uuid`
- `GET: /api/my/locations` returns the details for any location with an assistant as well as any owned farms
- `GET: /api/my/locations/{name}` returns the details of the location specified by `name` IF the location is an owned farm or holds an assistant
- `GET: /api/my/nearby-locations` returns a list of the names of every nearby location (all locations of every island with atleast one assistant), for navigational purposes
- `GET: /api/my/plots` returns a list of the player's plots
- `GET: /api/my/plots/{uuid}` returns the plot specified by `uuid`
- `PUT: /api/my/plots/{uuid}/clear` returns the plot if successful in attempt to clear plot, no request body expected
- `POST: /api/my/plots/{uuid}/plant` returns the updated warehouse and plot data if successful in attempt to plant specified plant in plot, as well as the info on the next growth stage of the plant
- - **Request Body** Expects `name` of seed, `quantity` of seed, `size` of plant. Example:
```json
{
    "name": "Cabbage Seeds",
    "quantity": 10,
    "size": "Miniature (1)"
}
```

---

### Response Codes

See `responses.go`

## Roadmap & Changelog

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
- ~~Plot interaction endpoints for non-growth-action plot management (plant plot, clear plot)~~
- - ~~`/plant`~~
- - - ~~Implement failure responses for Plant helper~~
- - ~~`/clear`~~
- ~~Skelling and Tritum YAML defined~~
- ~~Convert from UUID to composite string like "username|farmid|plotid" e.g. "Greenitthe|Homestead Farm|Plot-1" for warehouses, farms, plots, assistants, contracts. UUID wont be used~~
- ~~Convert to symbol based sector/island/location~~
- ~~Convert regions to islands~~
- ~~Add sectors~~
- ~~Remove Quality property from goods. Reintroduce for Trophies, use unique items instead~~
- ~~Define goods in YAML rather than as enum?~~
- ~~Remove enchantment property from goods. Generate unique Goods for the limited enchantable list~~
- ~~Warehouses store map[goodName]quantity now instead of Good structs~~
- ~~Move plot storage to farms, rather than separate DB table~~
- ~~Refactor dictionaries to a main struct~~
- ~~Add Sickle, Shade Cloth tool/action, Spectral Grass plant~~
- ~~Plot `/interact` endpoint with switch on body.action (growth actions)~~
- Cooldown/growth time is enforced (written, tested, just need to enable by uncommenting when testing done)
- ~~Add warehouse increment/decrement methods to handle removing the key for goods that are now at 0~~
- ~~Fix bug with path selector (case sensitive)~~
- Harvest functionality in `/interact`
- - ~~Actually, harvests should give Produce, a superset of Good, and go to special warehouse section so they can have size but not every good~~
- - Check logic in `plots:interact()` for returning growthHarvest when harvest is optional and when harvest but not FinalHarvest
- - Produce deposited to warehouse
- ~~Warehouses have sections for tools, produce, goods, seeds~~
- - ~~Farm tools moved to warehouse~~
- ~~Separate YAML definitions for goods into produce, goods, seeds files~~
- Tested growth and harvest of cabbage, potatos, shelvis fig, spectral grass, INCLUDING optional actions (make sure yield properly adjusted)
- Add `GET: /plants/{plantName}/growth-stages/{index}`

---

### Planned: **[v0.4]**

- Assistants can transfer things between warehouses
- Boldor, Yoggoth YAML defined
- Request validation functions
- Validate 1:1 mapping for every seed and plant after loading both
- Look through log entries to ensure all going to correct namespace (debug, important, error, etc.)
- Look through responses and ensure all are using correct response code
- Look through response codes and ensure all are using correct http response code

---

### Planned: **[v0.5]**

- One functional market, Local orders (non-player orders to provide baseline supply/demand)
- Ratelimiting
- Add `GET` endpoints for sectors, `GET` select island endpoint
- Update documentation
- Deploy pre-alpha server and separate dev server once complete, host timed pre-alpha test via discord (no point till can sell produce and buy new seeds)

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
- Add Gigantic, Colossal, and Titanic planting. Harvesting these don't give Goods, but rather Trophies. Trophies are given unique IDs when harvested, yield contributes to Quality rather than Quantity, and hold Grower Username, Grown Date, Size, numeric Quality, Plant Rarity, and optionally (once sold) Sold Price. These are not taken to any old market, but are transported to special Auctioneer locations where they are auctioned for a set amount of time. Size, Quality, and Rarity are combined, +/- 10% from RNG to give a Sold Price. All trophies are stored in the DB for the season, and there are leaderboards for who can get the highest quality and highest sale price.
- - Add a type of plant that can be indefinitely yield boosted to allow more active competing over Trophy leaderboards

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
