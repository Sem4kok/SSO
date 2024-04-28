package migrator

import (
	"flag"
	"fmt"
	"os"
)

const (
	StoragePath     = 1
	MigrationsPath  = 2
	MigrationsTable = 3
)

// SetupMigrator returns MigratorSetup
func SetupMigrator() (*MigratorSetup, error) {
	migratorStp := &MigratorSetup{
		StoragePath:     getPath(StoragePath),
		MigrationsPath:  getPath(MigrationsPath),
		MigrationsTable: getPath(MigrationsTable),
	}

	switch "" {
	case migratorStp.MigrationsPath:
		return nil, fmt.Errorf("migration path must be specified")
	case migratorStp.StoragePath:
		return nil, fmt.Errorf("storage path must be specified")
	case migratorStp.MigrationsTable:
		return nil, fmt.Errorf("migrations table must be specified")
	}

	return migratorStp, nil
}

// getPath can parse flag > env > default
func getPath(argument int) string {
	var flagName, env, path string

	switch argument {
	case StoragePath:
		flagName, env = "storage_path", "STORAGE_PATH"
	case MigrationsPath:
		flagName, env = "migrations_path", "MIGRATIONS_PATH"
	case MigrationsTable:
		flagName, env = "migrations_table", "MIGRATIONS_TABLE"
	}

	flag.StringVar(&path, flagName, "", "")
	flag.Parse()

	if path == "" {
		path = os.Getenv(env)
	}

	return path
}
