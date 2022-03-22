# Apricate

Go-based server for a fantasy-themed capitalism simulator game set on a farm.

### Endpoints

**Public Routes**
<!-- - `GET: /api/v0/leaderboards` list all available leaderboards and their descriptions
- `GET: /api/v0/leaderboards/{board}` get the specified leaderboard rankings -->
- `GET: /api/islands` returns details on every island in the game, including port connections, for navigational purposes
- `GET: /api/islands/location-symbol}` returns details on specified island
- `GET: /api/users` returns lists of registered usernames with various filters: unique, active, etc.
- `GET: /api/users/{username}` returns the public user data
- `POST: /api/users/{username}/claim` attempts to claim the specified username, returns the user data after creation, including token which users must save to access private routes
- `GET: /api/plants` returns the data on every plant in the game
- `GET: /api/plants/{plantName}` returns the data on the specified plant
- `GET: /api/metrics` returns map of metrics, including global market buy/sell

**Secure Routes**
- `GET: /api/my/account` returns the private user data (includes token)
- `GET: /api/my/assistants` returns a list of the player's assistants
- `GET: /api/my/assistants/{id}` returns the assistant specified by `id` (use numeric ID, i.e. for the uuid: `Greenitthe|Assistant-0` call this with 0 e.g. `GET: /api/my/assistants/0`)
- `GET: /api/my/farms` returns a list of the player's farms
- `GET: /api/my/farms/{location-symbol}` returns the farm specified by `location-symbol` (for the uuid: `Greenitthe|Farm-TS-PR-HF` call this with TS-PR-HF e.g. `GET: /api/my/farms/TS-PR-HF`)
- `GET: /api/my/contracts` returns a list of the player's contracts
- `GET: /api/my/contracts/{id}` returns the contract specified by `id` (use numeric ID, i.e. for the uuid: `Greenitthe|Contract-0` call this with 0 e.g. `GET: /api/my/contracts/0`)
- `GET: /api/my/warehouses` returns a list of the player's warehouses
- `GET: /api/my/warehouses/{location-symbol}` returns the warehouse specified by `location-symbol` (for the uuid: `Greenitthe|Warehouse-TS-PR-HF` call this with 0 e.g. `GET: /api/my/warehouses/TS-PR-HF`)
- `GET: /api/my/locations` returns the details for any location with an assistant as well as any owned farms
- `GET: /api/my/locations/{location-symbol}` returns the details of the location specified by `location-symbol` IF the location is an owned farm or holds an assistant
- `GET: /api/my/nearby-locations` returns a list of the names of every nearby location (all locations of every island with atleast one assistant), for navigational purposes
- `GET: /api/my/plots` returns a list of the player's plots
- `GET: /api/my/plots/{plot-id}` returns the plot specified by `plot-id`  (for the uuid: `Greenitthe|Farm-TS-PR-HF|Plot-1` call this with 0 e.g. `GET: /api/my/warehouses/TS-PR-HF_Plot-1`)
- `PUT: /api/my/plots/{plot-id}/clear` returns the plot if successful in attempt to clear plot, no request body expected
- `POST: /api/my/plots/{plot-id}/plant` returns the updated warehouse and plot data if successful in attempt to plant specified plant in plot, as well as the info on the next growth stage of the plant
- - **Request Body** Expects `name` of seed, `quantity` of seed, `size` of plant. Example:
```json
{
    "name": "Cabbage Seeds",
    "quantity": 10,
    "size": "Miniature"
}
```
- `PATCH: /api/my/plots/{plot-id}/interact` returns the updated warehouse and plot data if successful in attempt to interact with specified plant in plot, as well as the info on the next growth stage of the plant
- - **Request Body** Expects `action` desired action, `consumable` name of good to be consumed (optional depending on step, will be ignored if passed to step that has no consumable options). Example:
```json
{
    "action": "Water",
    "consumable": "Enchanted Water"
}
```
- `GET: /api/my/markets` returns a list of the markets the player can see
- `GET: /api/my/markets/{location-symbol}` returns the market specified by `location-symbol`
- `PATCH: /api/my/markets/{location-symbol}/order` returns the user ledger and local warehouse data if successful in placing the market order specified by the request body
- - **Request Body** Expects `order_type` (`MARKET` only available currently, filled instantly at current market price), `transaction_type` from [`BUY`, `SELL`], `item_type` from [`PRODUCE`, `GOODS`, `SEEDS`, `TOOLS`], `item_name` of desired transactable item from local market listing (for `PRODUCE` include the `Size` in the name after a pipe e.g. `Potato|Miniature`), and `quantity` of good to be transacted. Example:
```json
{
    "order_type": "MARKET",
    "transaction_type": "SELL",
    "item_type": "PRODUCE",
    "item_name": "Potato|Tiny",
    "quantity": 1000
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

### Complete: **[v0.4]**

- ~~Add `GET` endpoints for regions, `GET` select island endpoint~~
- - ~~Sectors are regions now~~
- ~~Helper functions for getting entries from route_vars/mux.Vars(r) so that always returned in correct format (e.g. Title(ToLower), ToUpper, etc.)~~
- - ~~URI params are case insensitive~~
- - ~~Usernames are case insensitive, though saved and displayed casefully~~
- ~~Look through log entries to ensure all going to correct namespace (debug, important, error, etc.)~~
- ~~Look through responses and ensure all are using correct response code~~
- ~~Look through response codes and ensure all are using correct http response code~~
- ~~Update my user auto-creation for Greenitthe with everything in the game for testing~~
- ~~Add user auto-creation for Viridis (basic setup), so I still lock down the username ;)~~
- ~~Placeholder market at the farm itself with buy/sell `market` orders and set prices (probably some ledger currency helper funcs necessary)~~
- - ~~`GET` endpoints~~
- - ~~multiply good base value by the integer value of the Size to get total value~~
- - ~~metric tracking number of each item bought and sold~~
- - ~~Troubleshoot/test buy/sell, figure out produce specifics, maybe remove "Location" from market order request as that should be obvious from the `/my/markets/TS-PR-HF/order` endpoint~~
- ~~Update `my/.../{selector}` endpoints to not need the Username or type (e.g. Warehouse- or Assistant-) parts in selector (just `/my/markets/TS-PR-HF` for example), if using numeric id, add plain id to struct~~
- ~~Get for items bought and sold metric~~
- ~~Add user coins metric~~

---

### Started: **[v0.5]** First Public Alpha

- ~~Update starting user template with appropriate tools, seeds, goods, produce, currencies~~
- Update documentation
- Deploy alpha server and separate dev server once complete, host timed pre-alpha test via discord (no point till can sell produce and buy new seeds)
- Uncomment cooldown section for plants
- ~~Fix server crash when not specifying produce size in item name of market order~~

---

### Planned: **[v0.6]**

- Respond to feedback from v0.5 alpha
- Add YAML-defined contract/quest paths from NPCs
- NPC endpoints for story/lore and contracts
- - `/talk`
- Cheat codes, maybe entered by talking to the NPC on your farm
- Rename Contracts to QuestContracts to fit the theme better (yes, its a capitalist society, but its an oldtimey fantasy world, gotta have quests)
- NPCs defined in YAML
- Use data field for request error responses to convey programmatically what failed validation
- Request validation functions
- Validate 1:1 mapping for every seed and plant after loading both

---

### Planned: **[v0.7]**

- Respond to feedback from v0.6 alpha
- Ratelimiting
- Boldor, Yoggoth, Tyldia YAML defined
- Assistants can transfer things between warehouses
- Market uses 4 types of market order, simulates dynamic NPC supply/demand/pricing that evolves over time and based on all player investment in market
- Markets initial state defined in YAML
- At least 10 different plants excluding Wild Seeds, buy and sell at the markets of at least 3 towns, all starting town NPCs have quests
- Atomize functions, write tests

---

### Planned: **[v0.8]**

- Respond to feedback from v0.7 alpha
- At least 20 plants excluding Wild Seeds, add at least 2 additional tools for growing some of the new plants, add randomized contracts, consider adding additional markets, all NPCs on starting map have quests
- Leaderboards (basically top-10 of ranked metrics?)

---

### Planned: **[v0.9]**

- Respond to feedback from v0.8 alpha
- Add refining/crafting with at least 8 recipes, add at least 4 new tools to support crafting, add at least 2 new buildings to support crafting
- Add researching plants (with associated building) to reveal full information
- Fog of War and hide unresearched plant information

---

### Planned: **[v1.0]**

- Respond to feedback from v0.9 alpha
- Meta account and progression, leaderboards, full documentation, separate dev partition that won't affect live
- Live server that is persistent, only wiped on update day. Updates are pushed every (other?) week when available
- Wild Seeds implemented?

---

### Planned: Post-1.0

- Simplified routing (pass a full route, even over oceans, and server will calculate total fare and travel time, rather than requiring manual travel between each intermediate location) - meta progression unlock?
- Add Gigantic, Colossal, and Titanic planting. Harvesting these don't give Goods, but rather Trophies. Trophies are given unique IDs when harvested, yield contributes to Quality rather than Quantity, and hold Grower Username, Grown Date, Size, numeric Quality, Plant Rarity, and optionally (once sold) Sold Price. These are not taken to any old market, but are transported to special Auctioneer locations where they are auctioned for a set amount of time. Size, Quality, and Rarity are combined, +/- 10% from RNG to give a Sold Price. All trophies are stored in the DB for the season, and there are leaderboards for who can get the highest quality and highest sale price.
- - Add a type of plant that can be indefinitely yield boosted to allow more active competing over Trophy leaderboards
- Add a plant that grows tools

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
