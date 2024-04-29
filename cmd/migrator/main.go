package main

import (
	"SSO/internal/migrator"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

func main() {
	// get flags/env path's
	mgt, err := migrator.SetupMigrator()
	if err != nil {
		if errors.Is(migrator.ErrorMigrationTableEmpty, err) {
			fmt.Printf("specified --migrations_table=./\n")
		} else {
			panic(err)
		}
	}

	db, err := sql.Open(mgt.GooseDriver, mgt.StoragePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = db.Close()
	}()

	_ = goose.SetDialect(mgt.GooseDriver)
	if err := goose.Up(db, mgt.MigrationsPath); err != nil {
		panic(err)
	}
}
