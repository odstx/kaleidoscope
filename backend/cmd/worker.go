package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"kaleidoscope/config"
	"kaleidoscope/worker"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start the Asynq worker",
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
			os.Exit(1)
		}

		// Initialize logger
		logger, err := config.InitLogger(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}
		defer logger.Sync()

		// Create Asynq worker
		asynqWorker := worker.NewWorker(
			fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
			cfg.Redis.Password,
			cfg.Redis.DB,
			&cfg.Email,
			logger,
		)

		// Start worker
		logger.Info("Starting Asynq worker...")
		if err := asynqWorker.Start(); err != nil {
			logger.Fatal("Failed to start Asynq worker", zap.Error(err))
		}
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
