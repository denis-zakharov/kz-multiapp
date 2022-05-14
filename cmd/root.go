package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/fn/framework/command"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type Spec struct {
	Value string `yaml:"value,omitempty"`
}
type Example struct {
	Spec Spec `yaml:"spec,omitempty"`
}

func createResourceListProcessor() *framework.SimpleProcessor {
	functionConfig := &Example{}

	fn := func(items []*yaml.RNode) ([]*yaml.RNode, error) {
		for i := range items {
			err := items[i].PipeE(yaml.SetAnnotation("custom.io/the-value", functionConfig.Spec.Value))
			if err != nil {
				return nil, err
			}
		}
		return items, nil
	}

	return &framework.SimpleProcessor{Config: functionConfig, Filter: kio.FilterFunc(fn)}
}

func createEmbeddedCommand() *cobra.Command {
	cmd := command.Build(processor, command.StandaloneDisabled, false)
	cmd.Use = "javagen"
	cmd.Short = "Generate Java App Kubernetes resources"
	cmd.Long = `Generate Kubernetes resources as as a kustomize transfomer.
						Implemented as a containerized KRM function.`
	return cmd
}

var processor = createResourceListProcessor()
var rootCmd = createEmbeddedCommand()

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
