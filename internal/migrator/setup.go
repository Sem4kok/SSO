package migrator

import (
	"flag"
	"os"
)

const (
	StoragePath     = 1
	MigrationsPath  = 2
	MigrationsTable = 3
)

// SetupMigrator returns MigratorSetup
func SetupMigrator() MigratorSetup {
	return MigratorSetup{
		StoragePath:     getPath(StoragePath),
		MigrationsPath:  getPath(MigrationsPath),
		MigrationsTable: getPath(MigrationsTable),
	}
}

// getPath can parse flag > env
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

	return path
}
