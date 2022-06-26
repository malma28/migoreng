# How to use?
## Installation
Install this module by running the following command
```
go get github.com/malma28/migoreng
```
Sometimes you already have installed this module instead you want to get with the latest version, run the following command
```
go get -u github.com/malma28/migoreng
```
## Tutorial
### 1. Create Migrator
First of all we have to make a migrator, we already have functions to create migrator.
#### - SQL
`NewSQL` to create Migrator for the database that use SQL (e.g `PostgreSQL`, `MySQL` / `MariaDB`). This function requires 3 parameters, the first is the database that you want to use, the second is the `*sql.DB` instance that was previously opened by you, and the third is the options for your migrator, use `nil` to use the default options.
### 2. Set Source
Currently the migrator doesn't know what to do, so we have to set the source to the migrator, the source here is `[]Source`, where `Source` itself is a struct which has 3 fields
- `Id` which is useful as a name for the Source.
- `Up` which is a function that will be executed by Migrator when doing up migration.
- `Down` which is a function that will be executed by Migrator when doing down migration and when `Up` makes an error.

To set the source to the migrator, use the `SetSource` method, where the parameter is the Sources you want to set (`[]Source`).
### 3. Doing Migration
Now migrator knows what to run, but migrator doesn't know how to run it.

There are 2 types of migration:
- Up, when doing this migration, migrator will see the last version we migrated, then migrator will run the latest Source that it hasn't run if it exists (if it doesn't exist, then migrator will not run anything and migration is done), if it's run successfully then the migrator will increase the migration version, if it doesn't work then the migrator will migrate down.
- Down, there are 2 conditions where this migration is carried out, the first is called by ourselves, if this migration is done by ourselves, then the migrator will see the version last time we migrated and will run the Source that we last ran, and other conditions for this migration is run is that when the Up migration fails to run, it will run the same Source that the Up migration ran.

To migrate Up, use the `Up` method, this method requires one parameter, namely `n`, where `n` is how many total sources the migrator will run. For example when we migrate we have 5 sources in the migrator, then the current version from our migration is version 1, then we run `Up` method where `n` is 3, then migrator will run Source 2, Source 3, and Source 4.

And then to migrate Down, use the `Down` method, this method requires one parameter too, namely `n`, where `n` is how many total sources the migrator will run. For example when we migrate we have 5 sources in the migrator, then the current version from our migration is version 4, then we run `Down` method where `n` is 2, then migrator will run Source 3 and Source 2.

What happens if we set `n` to exceed the number of Sources?

Then the migrator will migrate until the last source, for example we migrate up with the current migration version is version 1, the total of our Sources is 5, we set the `n` to 10, then the migrator will only run until Source 5.
For Down migration, the migrator will run until the first Source.

Use `n` < 1 to migrate until the latest source (for Up migration) or the first source (for Down migration).

# Supported Database
## SQL
- PostgreSQL
- MySQL / MariaDB