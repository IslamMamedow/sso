package migrator

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate"
)

func main() {
	var storagePath, migrationsPath, migrationTable string

	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationTable, "migration-table", "migrations", "name of the migrations table")
	flag.Parse()

	if storagePath == "" || migrationsPath == "" {
		panic("storage-path and migrations-path are required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table-name=%s", storagePath, migrationTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("Migrations applied")
}
