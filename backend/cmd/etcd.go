package cmd

import (
	"context"
	"fmt"
	"os"

	"kaleidoscope/etcd"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var etcdEndpoints []string
var etcdUsername string
var etcdPassword string

var registerCmd = &cobra.Command{
	Use:   "register <app_name> <version> <instance_id> <endpoint>",
	Short: "Register a microservice instance with etcd",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewProduction()
		defer logger.Sync()

		client, err := etcd.New(&etcd.Config{
			Endpoints: etcdEndpoints,
			Username:  etcdUsername,
			Password:  etcdPassword,
		}, logger)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to connect to etcd: %v\n", err)
			os.Exit(1)
		}
		defer client.Close()

		appName := args[0]
		version := args[1]
		instanceID := args[2]
		endpoint := args[3]

		if err := client.RegisterService(context.Background(), appName, version, instanceID, endpoint, nil); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to register service: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Registered: %s/%s/%s -> %s\n", appName, version, instanceID, endpoint)
	},
}

var deregisterCmd = &cobra.Command{
	Use:   "deregister <app_name> <version> <instance_id>",
	Short: "Deregister a microservice instance from etcd",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewProduction()
		defer logger.Sync()

		client, err := etcd.New(&etcd.Config{
			Endpoints: etcdEndpoints,
			Username:  etcdUsername,
			Password:  etcdPassword,
		}, logger)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to connect to etcd: %v\n", err)
			os.Exit(1)
		}
		defer client.Close()

		appName := args[0]
		version := args[1]
		instanceID := args[2]

		if err := client.DeregisterService(context.Background(), appName, version, instanceID); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to deregister service: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Deregistered: %s/%s/%s\n", appName, version, instanceID)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered service instances",
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewProduction()
		defer logger.Sync()

		client, err := etcd.New(&etcd.Config{
			Endpoints: etcdEndpoints,
			Username:  etcdUsername,
			Password:  etcdPassword,
		}, logger)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to connect to etcd: %v\n", err)
			os.Exit(1)
		}
		defer client.Close()

		services, err := client.ListServices(context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to list services: %v\n", err)
			os.Exit(1)
		}

		if len(services) == 0 {
			fmt.Println("No services registered")
			return
		}

		for appName, info := range services {
			fmt.Printf("\n%s:\n", appName)
			for _, instance := range info.Instances {
				fmt.Printf("  - %s/%s -> %s (%s)\n",
					instance.Version, instance.InstanceID, instance.Endpoint, instance.Status)
			}
		}
	},
}

var etcdCmd = &cobra.Command{
	Use:   "etcd",
	Short: "Manage etcd service registry",
}

func init() {
	etcdCmd.PersistentFlags().StringSliceVarP(&etcdEndpoints, "endpoints", "e", []string{"localhost:2379"}, "etcd endpoints")
	etcdCmd.PersistentFlags().StringVarP(&etcdUsername, "username", "u", "", "etcd username")
	etcdCmd.PersistentFlags().StringVarP(&etcdPassword, "password", "p", "", "etcd password")

	etcdCmd.AddCommand(registerCmd)
	etcdCmd.AddCommand(deregisterCmd)
	etcdCmd.AddCommand(listCmd)
	rootCmd.AddCommand(etcdCmd)
}
