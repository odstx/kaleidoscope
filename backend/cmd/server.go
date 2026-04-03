package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"kaleidoscope/config"
	"kaleidoscope/server"
)

var serverPort string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
			os.Exit(1)
		}

		// Override port if specified via command line flag
		if serverPort != "" {
			cfg.Server.Port = serverPort
		}

		// Initialize logger
		logger, err := config.InitLogger(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}
		defer logger.Sync()

		// Create server instance
		s := server.NewServer(logger, cfg)

		// Start server
		if err := s.Start(); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}

		// Wait for shutdown signal
		s.WaitForShutdown()
	},
}

func init() {
	serverCmd.Flags().StringVarP(&serverPort, "port", "p", "", "Server port (default: 9000)")
	rootCmd.AddCommand(serverCmd)
}
