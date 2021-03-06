package cmd

import (
	"github.com/spf13/cobra"
)

const (
	SUCCESS = "Success"
	FAILED  = "Failed"
)

var rootCmd = &cobra.Command{
	Use:   "argo-agent",
	Short: "Codefresh argocd agent",
	Long:  `Codefresh argocd agent`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
