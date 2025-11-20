# Mini-Processor Takehome Project

> NOTE: The go.mod file was left as github.com/censys/scan-takehome intentionally to avoid adding the dependency.
> Additionally, the original README.md was renamed to OLD-README.md for reference.

## Assumptions

1. It was ok to modify the `cmd/scanner/Dockerfile` to allow it to build with a newer version of go (the original README.md only mentioned not modifying `cmd/scanner/main.go`).
2. It was ok to copy the original repo to this new repo for ease of testing / use.
3. For the purposes of the testing assignment it was ok to mix integration & unit tests, whereas normally they would be separated, for time savings.

## Running The Scanner &amp; Processor

To run the scanner and processor together, run: `./project start`.

This will start the pubsub emulator, topic creation, subscription creation, scanner, database, flyway and processor docker images.

The database can be connected to via `localhost:5432` (assuming the default `DATABASE_TYPE` of `postgres` is used) with your favourite SQL tool (e.g. JetBrains DataGrip).

## Adding New Databases

The `internal/database` packages contains a RegisterDB function which allows for new database implementations to be added.

To add a new database implementation:

1. create a new package under `internal/database/<new-db-type>` and implement the `dal.Scan` interface.
1. In the `init()` function of the new package, call `database.RegisterDB("<new-db-type>", <DB Creation Func>)` where `<DB Creation Func>` is a function that returns a new instance of the database implementation.
    * See the `internal/database/noop` or `internal/database/psql` packages for examples.
1. If necessary, add new migration files to `third_party/flyway/<new-db-type>` and update the `third_party/flyway/Dockerfile` to include the new migration files instead of the `psql` ones.
1. Update the `docker-compose.yml` file to change the `DATABASE_TYPE` environment variable to `<new-db-type>` and update any additional environment variables needed to connect to the new database.

> WARNING: The `<new-db-type>` string must match the database connection string prefix (e.g. `postgres` to start a `postgres://` connection string).

The database schemas under `third_party/flyway/psql` do not currently contain any postgres specific SQL and should be able to be re-used with most SQL databases with little to no modification.
If the database you are adding requires different SQL syntax, you will need to create new migration files under a new folder in `third_party/flyway/<new-db-type>` and update the `third_party/flyway/Dockerfile` to include the new migration files instead of the psql ones.

> NOTE: The db system may not be perfect for every type of database; it is relatively simplistic to allow for the time constraints of this takehome project.

For the purposes of this takehome project, I used the `ON CONFLICT` syntax of Postgres inserts to handle the upsert logic rather than adding additional go code to handle retrieval, comparison and update of the existing record.

This approach may not be viable for all databases, as not all databases support this syntax.

### No-Op Database

The no-op database was implemented as an aid to testing and to initially verify that reads, inserts and ACKs were working correctly without needing to set up a full database.

## Testing

The project script provides the ability to run all testing (including integration tests) via: `./project test`

This will start the pubsub emulator, topic creation and subscription creation docker images.
Once these have started, the unit tests and integration tests will run together.

The coverage data file will be generated in the `<repository_root>/.reports/coverage.out` file.

> NOTE: Typically I would separate unit tests and integration tests, but for the purposes of this takehome project I have combined them to save time.

### Manual Testing

Additional manual testing was done to verify the processor was working correctly with the scanner and that the `ON CONFLICT` logic for postgres was working correctly.
To manually test the `ON CONFLICT` logic I:

1. Started the database with `docker compose up -f docker-compose.yml -d flyway` (flyway depends on the `database` service so it will start the database as well)
1. Connected to the database via the Goland Database tool window
1. Inserted a sample record into the `scans` table
1. Attempted to insert a conflicting record with an older `last_scanned` timestamp and verified that the record was not updated
1. Attempted to insert a conflicting record with a newer `last_scanned` timestamp and verified that the record was updated

Additionally, for the full system validation:

1. Started the full system with `./project start`
1. Verified that the scanner was publishing messages to the Pub/Sub emulator by checking the logs of the `scanner` service
1. Verified that the processor was receiving messages from the Pub/Sub emulator and inserting/updating records in the database by:
    1. Connecting to the database via the Goland Database tool window
    1. Querying the `scans` table to see if records were being inserted/updated

### Automated Testing

The project contains both unit tests and integration tests (currently mixed together; see note(s) above).

The tests are automatically run via a GitHub Actions workflow on every push to a pull request in the repository.

After the tests are complete, and automated coverage check is run.

Normally the threshold would be set to 80% across the board, but for the purposes of this takehome project some file and package thresholds are set to 70% to allow for time constraints.

On a successful set of test runs, a text coverage report is generated and posted as a step summary to the workflow.

Additionally, the `docker-compose.yml` file is validated by attempting to bring the system up and then shut it down on each run.

## HTML Coverage Report

To check the coverage for the code via the standard golang HTML coverage report (again including integration tests), run: `./project coverage`

This will start the pubsub emulator, topic creation and subscription creation docker images.
Once these have started, the unit tests and integration tests will run together and present the HTML coverage on completion.
