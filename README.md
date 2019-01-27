# mig

mig by [@satriahrh](https://github.com/satriahrh) is a database migration tool forked from [volatiletech/mig](https://github.com/volatiletech/mig). Manage your database's evolution by creating incremental SQL files.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/satriahrh/mig/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/satriahrh/mig)](http://goreportcard.com/report/satriahrh/mig)

### Goals of this fork

To provide a highly effective MySQL migration according to my needs. You can watch my development progress on the [project tab](https://github.com/satriahrh/mig/projects) and please do help me too if you think it is worth enough to your life.

# Install

    $ go get -u github.com/satriahrh/mig/...

# Usage

```
mig is a database migration tool for MySQL.

Usage:
  mig [command]

Examples:
$ mig up "user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true"
$ mig down "user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true"
$ mig create add_users

Available Commands:
  create      Create a blank migration template
  down        Roll back the version by one
  downall     Roll back all migrations
  help        Help about any command
  redo        Down then up the latest migration
  redoall     Down then up all migrations
  status      Dump the migration status for the database
  up          Migrate the database to the most recent version available
  upone       Migrate the database by one version
  version     Print the current version of the database

Flags:
      --version   Print the mig tool version

Use "mig [command] --help" for more information about a command.
```

## Supported Databases

mig supports MySQL. The drivers used are:

https://github.com/go-sql-driver/mysql

See these drivers for details on the format of their connection strings.

## Couple of example runs

### create

Create a new SQL migration.

    $ mig create add_users
    $ Created db/migrations/20130106093224_add_users.sql

Edit the newly created script to define the behavior of your migration. Your
SQL statements should go below the Up and Down comments.

An example command run:

### up

Apply all available migrations.

    $ mig up "user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true"
	$ Success   20170314220650_add_dogs.sql
	$ Success   20170314221501_add_cats.sql
	$ Success   2 migrations

## Migrations

A sample SQL migration looks like:

```sql
-- +mig Up
CREATE TABLE post (
    id int NOT NULL,
    title text,
    body text,
    PRIMARY KEY(id)
);

-- +mig Down
DROP TABLE post;
```

Notice the annotations in the comments. Any statements following `-- +mig Up` will be executed as part of a forward migration, and any statements following `-- +mig Down` will be executed as part of a rollback.

By default, SQL statements are delimited by semicolons - in fact, query statements must end with a semicolon to be properly recognized by mig.

More complex statements (PL/pgSQL) that have semicolons within them must be annotated with `-- +mig StatementBegin` and `-- +mig StatementEnd` to be properly recognized. For example:

```sql
-- +mig Up
-- +mig StatementBegin
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
-- +mig StatementEnd
```

## Library functions


```go
// Global io.Writer variable that can be changed to get incremental success 
// messages from function calls that process more than one migration,
// for example Up and DownAll. Defaults to ioutil.Discard.
var mig.Log

// Create a templated migration file in dir
mig.Create(name, dir string) (name string, err error)

// Down rolls back the version by one
mig.Down(driver, conn, dir string) (name string, err error)

// DownAll rolls back all migrations.
// Logs success messages to global writer variable Log.
mig.DownAll(driver, conn, dir string) (count int, err error)

// Up migrates to the highest version available
mig.Up(driver, conn, dir string) (count int, err error)

// UpOne migrates one version
mig.UpOne(driver, conn, dir string) (name string, err error)

// Redo re-runs the latest migration.
mig.Redo(driver, conn, dir string) (name string, err error)

// Return the status of each migration
mig.Status(driver, conn, dir string) (status, error)

// Return the current migration version
mig.Version(driver, conn string) (version int64, err error)
```
