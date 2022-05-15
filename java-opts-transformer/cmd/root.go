package cmd

import (
	"os"

	"github.com/denis-zakharov/kz-multiapp/java-opts-transformer/pkg"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/kyaml/fn/framework/command"
)

func Execute(mode command.CLIMode) {
	rootCmd := buildRootCmd(mode)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func buildRootCmd(mode command.CLIMode) *cobra.Command {
	config := &pkg.JavaOptsTransformerv1alpha1{}
	p := config.BuildProcessor()

	cmd := command.Build(p, mode, false)

	cmd.Use = "fn"
	cmd.Short = "Inject JAVA_OPTIONS env var"
	cmd.Long = `Compute and inject JAVA_OPTIONS env var and resources
according to input parameters`

	// gen Dockerfile
	command.AddGenerateDockerfile(cmd)

	return cmd
}
