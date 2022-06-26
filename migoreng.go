package migoreng

import "database/sql"

type DatabaseSQL int

const (
	DatabasePostgresql DatabaseSQL = iota
	DatabaseMysql
)

type Source struct {
	Id   string
	Up   func(db *sql.DB) error
	Down func(db *sql.DB) error
}

type Migrator interface {
	// use n < 0 to go up to the latest
	Up(n int) error
	// use n < 0 to go down to the first
	Down(n int) error
	// set the sources
}

type MigratorSQL interface {
	Migrator
	SetSource(sources []Source) error
}

type MigratorOptions struct {
	TableMigrationName string
}

var defaultMigratorOptions = &MigratorOptions{
	TableMigrationName: "migoreng_table_migration",
}

// set options to nil to use default options
func NewSQL(database DatabaseSQL, db *sql.DB, options *MigratorOptions) MigratorSQL {
	if options == nil {
		options = defaultMigratorOptions
	}

	switch database {
	case DatabasePostgresql:
		return &postgresqlMigrator{
			db:      db,
			options: options,
		}
	case DatabaseMysql:
		return &postgresqlMigrator{
			db:      db,
			options: options,
		}
	}

	return nil
}
