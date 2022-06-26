package migoreng

import (
	"database/sql"
	"errors"
	"fmt"
)

type mysqlMigrator struct {
	db      *sql.DB
	options *MigratorOptions
	sources []Source
}

func (migrator *mysqlMigrator) init() (int, error) {
	// create the table migration if not exist
	if _, err := migrator.db.Exec(
		fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS %v (id INTEGER NOT NULL DEFAULT 1 PRIMARY KEY, version INTEGER NOT NULL);",
			defaultMigratorOptions.TableMigrationName,
		),
	); err != nil {
		return 0, err
	}

	// Check if row exist
	row := migrator.db.QueryRow(
		fmt.Sprintf(
			"SELECT EXISTS(SELECT * FROM %v WHERE id = 1);",
			defaultMigratorOptions.TableMigrationName,
		),
	)
	if err := row.Err(); err != nil {
		return 0, err
	}

	var rowExist bool
	if err := row.Scan(&rowExist); err != nil {
		return 0, err
	}

	// Insert the row if not exist
	if !rowExist {
		tx, err := migrator.db.Begin()
		if err != nil {
			return 0, err
		}

		if _, err := tx.Exec(
			fmt.Sprintf(
				"INSERT INTO %v (id, version) VALUES (1, 0);",
				defaultMigratorOptions.TableMigrationName,
			),
		); err != nil {
			return 0, err
		}

		if err := tx.Commit(); err != nil {
			return 0, err
		}

		return 0, nil
	}

	row = migrator.db.QueryRow(
		fmt.Sprintf(
			"SELECT version FROM %v WHERE id = 1;",
			defaultMigratorOptions.TableMigrationName,
		),
	)
	if err := row.Err(); err != nil {
		return 0, err
	}

	var version int
	if err := row.Scan(&version); err != nil {
		return 0, err
	}

	return version, nil
}

func (migrator *mysqlMigrator) Up(n int) error {
	version, err := migrator.init()
	if err != nil {
		return err
	}

	if n < 0 {
		n = len(migrator.sources)
	}

	currentVersion := 0
	sourcesLength := len(migrator.sources)
	for i := 0; i < n && version+i < sourcesLength; i++ {
		currentVersion = version + i
		if err := migrator.sources[currentVersion].Up(migrator.db); err != nil {
			return err
		}
	}

	if _, err := migrator.db.Exec(
		fmt.Sprintf(
			"UPDATE %v SET version = $1 WHERE id = 1;",
			migrator.options.TableMigrationName,
		),
		currentVersion+1,
	); err != nil {
		return err
	}

	return nil
}

func (migrator *mysqlMigrator) Down(n int) error {
	version, err := migrator.init()
	if err != nil {
		return err
	}

	if n < 0 {
		n = len(migrator.sources)
	}

	currentVersion := 0
	for i := 0; i < n && version-(i+1) >= 0; i++ {
		currentVersion = version - (i + 1)
		if err := migrator.sources[currentVersion].Down(migrator.db); err != nil {
			return err
		}
	}

	if _, err := migrator.db.Exec(
		fmt.Sprintf(
			"UPDATE %v SET version = $1 WHERE id = 1;",
			migrator.options.TableMigrationName,
		),
		currentVersion,
	); err != nil {
		return err
	}

	return nil
}

func (migrator *mysqlMigrator) SetSource(sources []Source) error {
	if sources == nil {
		return errors.New("sources is nil")
	}
	for _, source := range sources {
		if source.Up == nil || source.Down == nil {
			return fmt.Errorf("source.Up / source.Down field is nil in source that have id \"%v\"", source.Id)
		}
	}
	migrator.sources = sources
	return nil
}
