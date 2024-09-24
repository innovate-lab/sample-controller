package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"os"
)

var rootCmd = &cobra.Command{
	Use:          "kubewg [command]",
	SilenceUsage: true,
	Short:        "kubewg is cni plugin for wireguard",
	Long:         `kubewg is cni plugin for wireguard`,
}

func Execute() {
	rootCmd.AddCommand(sampleController())
	if err := rootCmd.Execute(); err != nil {
		klog.Errorf("rootCmd execute failed: %v", err)
		os.Exit(-1)
	}
}
