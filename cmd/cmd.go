package cmd

import "github.com/spf13/cobra"

func New() *cobra.Command {
	c := &cobra.Command{
		Use:   "invasion",
		Short: "Easy tool for simulating aliens invasion.",
		Long: `We constantly live in danger of an aliens invasion. This tool will help you to be more prepared.
Using invasion you can simulate any type of aliens invasion scenario against any earth area. Be ready for an invasion!`,

		SilenceUsage:  true,
		SilenceErrors: true,
	}

	c.AddCommand(NewSimulate())

	return c
}
