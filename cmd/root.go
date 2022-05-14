package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/fn/framework/command"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type Metadata struct {
	Name string `yaml:"name"`
}

type Spec struct {
	Value string `yaml:"value,omitempty"`
}

type v1alpha1Annotator struct {
	Metadata Metadata `yaml:"metadata"`
	Spec     Spec     `yaml:"spec,omitempty"`
}

func createResourceListProcessor() *framework.SimpleProcessor {
	functionConfig := &v1alpha1Annotator{}

	fn := func(items []*yaml.RNode) ([]*yaml.RNode, error) {
		for i := range items {
			err := items[i].PipeE(yaml.SetAnnotation("custom.io/the-value", functionConfig.Spec.Value))
			if err != nil {
				return nil, err
			}
			err = items[i].PipeE(yaml.SetAnnotation("custom.io/the-name", functionConfig.Metadata.Name))
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
