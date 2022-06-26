package migoreng

import (
	"database/sql"
	"errors"
	"fmt"
)

type postgresqlMigrator struct {
	db      *sql.DB
	sources []Source
	options *MigratorOptions
}

func (migrator *postgresqlMigrator) init() (int, error) {
	// Check if table migration exist
	row := migrator.db.QueryRow(
		fmt.Sprintf(
			"SELECT EXISTS(SELECT * FROM information_schema.tables WHERE table_name = '%v');",
			defaultMigratorOptions.TableMigrationName,
		),
	)
	if err := row.Err(); err != nil {
		return 0, err
	}

	// If not exist, create the table migration
	var tableExist bool
	if err := row.Scan(&tableExist); err != nil {
		return 0, err
	}
	if !tableExist {
		if _, err := migrator.db.Exec(
			fmt.Sprintf(
				"CREATE TABLE %v (id INTEGER NOT NULL DEFAULT 1 PRIMARY KEY, version INTEGER NOT NULL);",
				defaultMigratorOptions.TableMigrationName,
			),
		); err != nil {
			return 0, err
		}
	}

	// Check if row exist, but if the table not exist before, we dont need to check
	// Instead we directly create row
	rowExist := false
	if tableExist {
		row = migrator.db.QueryRow(
			fmt.Sprintf(
				"SELECT EXISTS(SELECT 1 FROM %v WHERE id = 1);",
				defaultMigratorOptions.TableMigrationName,
			),
		)
		if err := row.Err(); err != nil {
			return 0, err
		}
		if err := row.Scan(&rowExist); err != nil {
			return 0, err
		}
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

func (migrator *postgresqlMigrator) Up(n int) error {
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

func (migrator *postgresqlMigrator) Down(n int) error {
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

func (migrator *postgresqlMigrator) SetSource(sources []Source) error {
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
