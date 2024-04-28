package migrator

type MigratorSetup struct {
	StoragePath     string
	MigrationsPath  string
	MigrationsTable string
	GooseDriver     string
}
