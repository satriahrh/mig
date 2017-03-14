# mig

mig is a database migration tool. Manage your database's evolution by creating incremental SQL files.

[![GoDoc Widget]][GoDoc]

### Goals of this fork

This is a heavily restructured and rewritten fork of https://github.com/pressly/goose
which includes many necessary bug fixes and general improvements to
better suit the Go abcweb framework https://github.com/nullbio/abcweb --
although feel free to use this migration tool as a standalone or in your own 
projects even if you do not use abcweb, it by no means requires it.

# Install

    $ go get -u github.com/nullbio/mig/...

# Usage

```
mig is a database migration tool for Postgres, MySQL and SQLite3.

Usage:
  mig [command]

Examples:
mig postgres "user=postgres dbname=postgres sslmode=disable" up
mig mysql "user:password@/dbname" down
mig sqlite3 ./foo.db status
mig create add_users

Available Commands:
  create      Create a blank migration template
  down        Roll back the version by one
  help        Help about any command
  redo        Re-run the latest migration
  status      Dump the migration status for the database
  up          Migrate the database to the most recent version available
  up          Migrate the database to the most recent version available

Flags:
      --version   Print the mig tool version

Use "mig [command] --help" for more information about a command.
```
## create

Create a new SQL migration.

    $ mig create add_users
    $ Created db/migrations/20130106093224_add_users.sql

Edit the newly created script to define the behavior of your migration. Your
SQL statements should go below the Up and Down comments.

## up

Apply all available migrations.

    $ goose up
    $ goose: migrating db environment 'development', current version: 0, target: 3
    $ OK    001_basics.sql
    $ OK    002_next.sql
    $ OK    003_and_again.go

## down

Roll back a single migration from the current version.

    $ goose down
    $ goose: migrating db environment 'development', current version: 3, target: 2
    $ OK    003_and_again.go

## redo

Roll back the most recently applied migration, then run it again.

    $ goose redo
    $ goose: migrating db environment 'development', current version: 3, target: 2
    $ OK    003_and_again.go
    $ goose: migrating db environment 'development', current version: 2, target: 3
    $ OK    003_and_again.go

## status

Print the status of all migrations:

    $ goose status
    $ goose: status for environment 'development'
    $   Applied At                  Migration
    $   =======================================
    $   Sun Jan  6 11:25:03 2013 -- 001_basics.sql
    $   Sun Jan  6 11:25:03 2013 -- 002_next.sql
    $   Pending                  -- 003_and_again.go

## dbversion

Print the current version of the database:

    $ goose dbversion
    $ goose: dbversion 002

# Migrations

goose supports migrations written in SQL or in Go.

## SQL Migrations

A sample SQL migration looks like:

```sql
-- +goose Up
CREATE TABLE post (
    id int NOT NULL,
    title text,
    body text,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE post;
```

Notice the annotations in the comments. Any statements following `-- +goose Up` will be executed as part of a forward migration, and any statements following `-- +goose Down` will be executed as part of a rollback.

By default, SQL statements are delimited by semicolons - in fact, query statements must end with a semicolon to be properly recognized by goose.

More complex statements (PL/pgSQL) that have semicolons within them must be annotated with `-- +goose StatementBegin` and `-- +goose StatementEnd` to be properly recognized. For example:

```sql
-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION histories_partition_creation( DATE, DATE )
returns void AS $$
DECLARE
  create_query text;
BEGIN
  FOR create_query IN SELECT
      'CREATE TABLE IF NOT EXISTS histories_'
      || TO_CHAR( d, 'YYYY_MM' )
      || ' ( CHECK( created_at >= timestamp '''
      || TO_CHAR( d, 'YYYY-MM-DD 00:00:00' )
      || ''' AND created_at < timestamp '''
      || TO_CHAR( d + INTERVAL '1 month', 'YYYY-MM-DD 00:00:00' )
      || ''' ) ) inherits ( histories );'
    FROM generate_series( $1, $2, '1 month' ) AS d
  LOOP
    EXECUTE create_query;
  END LOOP;  -- LOOP END
END;         -- FUNCTION END
$$
language plpgsql;
-- +goose StatementEnd
```

## Go Migrations

Import `github.com/pressly/goose` from your own project (see [example](./example/migrations-go/cmd/main.go)), register migration functions and run goose command (ie. `goose.Up(db *sql.DB, dir string)`).

A [sample Go migration 00002_users_add_email.go file](./example/migrations-go/00002_users_add_email.go) looks like:

```go
package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up, Down)
}

func Up(tx *sql.Tx) error {
	_, err := tx.Query("ALTER TABLE users ADD COLUMN email text DEFAULT '' NOT NULL;")
	if err != nil {
		return err
	}
	return nil
}

func Down(tx *sql.Tx) error {
	_, err := tx.Query("ALTER TABLE users DROP COLUMN email;")
	if err != nil {
		return err
	}
	return nil
}
```

## License

Licensed under [MIT License](./LICENSE)

[GoDoc]: https://godoc.org/github.com/pressly/goose
[GoDoc Widget]: https://godoc.org/github.com/pressly/goose?status.svg
