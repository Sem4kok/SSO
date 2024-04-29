package migrator

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	ErrorMigrationTableEmpty = errors.New("migrations table must be specified")
)

const (
	StoragePath     = 1
	MigrationsPath  = 2
	MigrationsTable = 3
	GooseDriver     = 4
)

// SetupMigrator returns MigratorSetup
func SetupMigrator() (*MigratorSetup, error) {
	storagePath := flag.String("storage_path", "", "Path to storage")
	migrationsPath := flag.String("migrations_path", "", "Path to migrations")
	migrationsTable := flag.String("migrations_table", "", "Migrations table name")
	gooseDriver := flag.String("goose_driver", "", "Goose database driver")
	flag.Parse()

	migratorStp := &MigratorSetup{
		StoragePath:     getPath(*storagePath, StoragePath),
		MigrationsPath:  getPath(*migrationsPath, MigrationsPath),
		MigrationsTable: getPath(*migrationsTable, MigrationsTable),
		GooseDriver:     getPath(*gooseDriver, GooseDriver),
	}

	switch "" {
	case migratorStp.MigrationsPath:
		return nil, fmt.Errorf("migration path must be specified")
	case migratorStp.StoragePath:
		return nil, fmt.Errorf("storage path must be specified")
	case migratorStp.MigrationsTable:
		return nil, ErrorMigrationTableEmpty
	case migratorStp.GooseDriver:
		return nil, fmt.Errorf("goose driver must be specified")

	}

	return migratorStp, nil
}

// getPath can parse flag > env > default
func getPath(argument string, argType int) string {
	var env string

	switch argType {
	case StoragePath:
		env = "STORAGE_PATH"
	case MigrationsPath:
		env = "MIGRATIONS_PATH"
	case MigrationsTable:
		env = "MIGRATIONS_TABLE"
	case GooseDriver:
		env = "GOOSE_DRIVER"
	}

	if argument == "" {
		argument = os.Getenv(env)
	}

	return argument
}
