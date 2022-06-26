package migoreng

import (
	"os"
	"strconv"
	"strings"
)

// Now you can enter the --migration={up/down}{step} argument when running your application.
//
// Example: ./your-application --migration=up4
func UseArgs(migrator Migrator) error {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "--migrate=") {
			migrateAction := strings.TrimPrefix(arg, "--migrate=")
			if strings.HasPrefix(migrateAction, "up") {
				stepStr := strings.TrimPrefix(migrateAction, "up")
				step, err := strconv.Atoi(stepStr)
				if err != nil {
					return err
				}
				if err := migrator.Up(step); err != nil {
					return err
				}
			} else if strings.HasPrefix(migrateAction, "down") {
				stepStr := strings.TrimPrefix(migrateAction, "down")
				step, err := strconv.Atoi(stepStr)
				if err != nil {
					return err
				}
				if err := migrator.Up(step); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
