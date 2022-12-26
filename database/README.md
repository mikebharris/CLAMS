# HacktionLab Workshops Database #

An SSH tunnel is created to the relevant database and the scripts are deployed using a [FlywayDB](https://flywaydb.org/) docker image. 

Please read the following documentation to become familiar with the concept of migrations:

* [Flyway Migrations Documentation](https://flywaydb.org/documentation/migrations)

## Getting Started

The following bullet points describe the way of working with this project.

* Make required changes on a `feature/<jira_number>` branch.
* Add SQL files to `sql/` with `*.sql` suffix and following flyway migrations naming conventions.
* Run **PGSanity Check** and local testing.
* If the changes look viable, commit and push to your `feature` branch.
* Manually run matching Jenkins multi-branch pipeline to deploy into `nonprod` environment.
* Deploy changes and verify them in the `nonprod` environment.
* Create merge request into `main` branch.

## PGSanity Check

As part of the automated pipeline, `pgsanity` is used to check the validity of SQL for PostgreSQL. This can also be called locally via `bash` or `gitbash`:

```sh
./pgsanity.sh
```

Note - `pgsanity` requires the following pre-requisities installed:

* Docker

### Example `pgsanity` Output

```s
Running script [/sql/V1__Initial.sql]
Running script [/sql/error.sql]
line 1: ERROR: unrecognized data type name "FISH"
```

## Creating Local Docker Testing Database

* Export local DB mode environment variables:

```sh
export LOCAL_DB_CREATION_MODE=true
export FLYWAY_COMMAND=migrate
export DOCKER_CONTAINER_NAME=reviewer-hub-local-db
```

* Run script to deploy migrations.

```sh
./flyway.sh
```