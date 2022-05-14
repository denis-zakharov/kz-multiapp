package cmd

import (
	"sigs.k8s.io/kustomize/kyaml/fn/framework/command"

	"github.com/spf13/cobra"
)

var debugCmd = createStandaloneCommand()

func createStandaloneCommand() *cobra.Command {
	cmd := command.Build(processor, command.StandaloneEnabled, false)
	cmd.Use = "debug"
	cmd.Short = "Generate Java App Kubernetes resources"
	cmd.Long = `Generate Kubernetes resources accepting
						FunctionConfig and ResourceList items files.
						Implemented as a containerized KRM function.`
	return cmd
}

func init() {
	rootCmd.AddCommand(debugCmd)
}
