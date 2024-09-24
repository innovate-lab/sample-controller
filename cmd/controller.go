package cmd

import (
	"github.com/spf13/cobra"
	"sample-controller/pkg/controller"
)

func sampleController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "controller",
		Short: "start controller",
		Long:  "start controller",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runController()
		},
	}

	return cmd
}

func runController() error {
	return controller.Run()
}
