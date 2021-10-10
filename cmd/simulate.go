package cmd

import (
	"fmt"
	"os"

	"github.com/ivanovpetr/invasion/services/simulator"
	"github.com/spf13/cobra"
)

const flagAliensNumber = "n"

func NewSimulate() *cobra.Command {
	c := &cobra.Command{
		Use:   "simulate [path/to/map]",
		Short: "simulates invasion of aliens",
		Long: `We constantly live in danger of an aliens invasion. This tool will help you to be more prepared.
Using invasion you can simulate any type of aliens invasion scenario against any earth area. Be ready for an invasion!`,
		Args: cobra.ExactArgs(1),
		RunE: simulateHandler,
	}

	c.Flags().Int(flagAliensNumber, 15, "Number of aliens during the simulation")

	return c
}

func simulateHandler(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	numberOfAliens, _ := cmd.Flags().GetInt(flagAliensNumber)
	// parse the provided map file
	simulation, err := simulator.CreateSimulationFromPath(filePath)
	if err != nil {
		return err
	}
	result := simulation.Run(int64(numberOfAliens))
	for _, log := range result.Logs {
		fmt.Println(log)
	}
	_ = result.PrintResultMap(os.Stdout)
	return nil
}
