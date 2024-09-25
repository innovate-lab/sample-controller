package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"os"
)

var rootCmd = &cobra.Command{
	Use:          "controller [command]",
	SilenceUsage: true,
	Short:        "",
	Long:         ``,
}

func Execute() {
	rootCmd.AddCommand(sampleController())
	if err := rootCmd.Execute(); err != nil {
		klog.Errorf("rootCmd execute failed: %v", err)
		os.Exit(-1)
	}
}
