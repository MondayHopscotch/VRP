package commands

import (
	"fmt"
	"os"
	"vrp/internal"
	"vrp/internal/parsing"
	"vrp/internal/routing"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "vrp",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expected single argument for load file path")
		}

		loads, err := parsing.ParseAllLoads(args[0])
		if err != nil {
			return err
		}

		solver := routing.NewNearestNeighborSolver(loads)
		routes := solver.PlanRoutes()

		for _, r := range routes {
			r.PrintLoadNumbers()
		}
		return nil
	},
}

// Execute runs the root command of the vrp cli tool
func Execute() {
	rootCmd.PersistentFlags().BoolVar(&internal.Debug, "debug", false, "enable debug output")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
