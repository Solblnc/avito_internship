package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
	"log"
)

func Migrate(conn *pgx.Conn, ctx context.Context) {
	migrator, err := migrate.NewMigrator(ctx, conn, "schema_version")
	if err != nil {
		log.Fatalf("Unable to create a migrator: %v\n", err)
	}

	err = migrator.LoadMigrations("./migrations")
	if err != nil {
		log.Fatalf("Unable to load migrations: %v\n", err)
	}

	err = migrator.Migrate(ctx)
	if err != nil {
		log.Fatalf("Unable to migrate: %v\n", err)
	}

	ver, err := migrator.GetCurrentVersion(ctx)
	if err != nil {
		log.Fatalf("Unable to get current schema version: %v\n", err)
	}

	log.Println("Migration done. Current schema version: %v\n", ver)
}
