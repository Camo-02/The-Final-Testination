# The Final Testination

## Run locally

# sudo sysctl -w kernel.apparmor_restrict_unprivileged_userns=0

A `.env` file is required in the root directory of the project. A `.env.sample` is available in the same directory with some sensible default values.

```sh
#### Run this first in one terminal and wait for the DB to start...
docker compose -f development.yml up # DB

#### ...then run this in a second terminal...
cd backend && ../scripts/addenv go run main.go # Backend

#### ...and then this in the third and last terminal
cd frontend && ../scripts/addenv npm run dev # Frontend
```

### Stop the compose

This command will stop the database

```sh
docker compose -f development.yml down
```

## Run tests locally

The following commands assume you are currently in the root of the project:

```sh
#### In one terminal
docker compose -f development.yml up

#### In another terminal
cd backend && ../scripts/addenv go test ./...

#### In a third terminal
cd frontend
npm i
../scripts/addenv npm run dev

#### In a fourth terminal
cd frontend
../scripts/addenv npm run build:test && ../scripts/addenv npm run host
# If you want to run all the tests from the cli run this in another terminal:
../scripts/addenv npm run cypress # add ` -- -b <browser-name>` (or path) to run on a specific broser (e.g. `-- -b chrome`). There's a specific shortcut for chrome: `npm run cypress:chrome`
# If you only want to run a specific type of tests, run:
../scripts/addenv npm run cypress:e2e # or cypress:component
# If you want to open cypress in a browser, then run
../scripts/addenv npm run cypress:open
```

**Note**: you have to `npm i` and `npm run dev` in order to have the dependencies and the `.svelte-kit` directory before running cypress tests, otherwise they will not work.

# psql

To run queries as strings against the local db, run:

```sh
psql -U postgres -p 5432 -h 127.0.0.1 -d final_testination -c "[query]"
```
