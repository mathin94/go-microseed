package main

import (
	"context"
	"fmt"
	"os"

	"microseed/internal/app"
	"microseed/internal/config"
	"microseed/internal/db"
	"microseed/internal/migrate"
	"microseed/internal/seed"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	root := &cobra.Command{
		Use:   "microseed",
		Short: "Microseed — Go microservice skeleton",
	}

	// serve
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Run HTTP API server",
		Run: func(cmd *cobra.Command, args []string) {
			fx.New(app.Module).Run()
		},
	}

	// migrate
	var steps int
	migrateCmd := &cobra.Command{Use: "migrate", Short: "Database migrations"}
	upCmd := &cobra.Command{
		Use: "up", Short: "Apply all up migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := config.New()
			log, _ := zap.NewProduction()
			defer log.Sync()
			return migrate.Up(cmd.Context(), cfg, log)
		},
	}
	downCmd := &cobra.Command{
		Use: "down", Short: "Rollback N migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := config.New()
			log, _ := zap.NewProduction()
			defer log.Sync()
			if steps <= 0 {
				steps = 1
			}
			return migrate.Down(cmd.Context(), cfg, log, steps)
		},
	}
	downCmd.Flags().IntVar(&steps, "step", 1, "steps to rollback")
	resetCmd := &cobra.Command{
		Use: "reset", Short: "Migrate down to version 0",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := config.New()
			log, _ := zap.NewProduction()
			defer log.Sync()
			return migrate.Reset(cmd.Context(), cfg, log)
		},
	}
	migrateCmd.AddCommand(upCmd, downCmd, resetCmd)

	// seed
	seedCmd := &cobra.Command{
		Use: "seed", Short: "Run data seeders",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			cfg, _ := config.New()
			log, _ := zap.NewProduction()
			defer log.Sync()
			gdb, err := db.NewGorm(cfg, log)
			if err != nil {
				return err
			}
			defer db.Close(gdb)
			return seed.SeedAll(ctx, gdb, log)
		},
	}

	root.AddCommand(serveCmd, migrateCmd, seedCmd)

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
