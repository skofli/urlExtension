package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/skofli/migrator"
	"os"
	"urlExtension/backend/api"
	"urlExtension/backend/store"
)

type Options struct {
	DBPath        string `long:"db" env:"DB" required:"true" description:"path to database with user:pass"`
	MigrationPath string `long:"migration-path" env:"MIGRATION_PATH" required:"true" description:"path to migration dir"`
	store.Options
}

func main() {
	options := parseOpts()
	err := migrator.Migrate(options.DBPath, options.MigrationPath)
	if err != nil {
		os.Exit(0)
	}
	db, err := store.New(options.DBPath)
	if err != nil {
		os.Exit(0)
	}
	server := api.New(db)
	if err := server.Run(options.Options); err != nil {
		os.Exit(0)
	}
}

func parseOpts() Options {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}
	return opts
}
