package cmd

import (
	"embed"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/fn/framework/command"
	"sigs.k8s.io/kustomize/kyaml/fn/framework/parser"
)

//go:embed templates/*
var tempalteFS embed.FS

type v1alpha1JavaAppGen struct {
	Metadata Metadata `yaml:"metadata"`
	Spec     Spec     `yaml:"spec,omitempty"`
}
type Metadata struct {
	Name string `yaml:"name"`
}
type Spec struct {
	Stateful  bool       `yaml:"stateful"`
	Replicas  int        `yaml:"replicas"`
	Image     string     `yaml:"image"`
	Ports     []PortItem `yaml:"ports"`
	JVM       JVM        `yaml:"jvm"`
	Env       []EnvItem  `yaml:"env"`
	Upstreams []string   `yaml:"upstreams"`
}
type PortItem struct {
	Name string `yaml:"name"`
	Port int    `yaml:"port"`
}
type JVM struct {
	Heap  int    `yaml:"heap"`
	Limit int    `yaml:"limit"`
	Extra string `yaml:"extra"`
}
type EnvItem struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func buildProcessor(app *v1alpha1JavaAppGen) framework.ResourceListProcessor {
	return framework.TemplateProcessor{
		ResourceTemplates: []framework.ResourceTemplate{{
			Templates: parser.TemplateFiles("templates").FromFS(tempalteFS),
		}},
		TemplateData: app,
	}
}

func createEmbeddedCommand() *cobra.Command {
	cmd := command.Build(processor, command.StandaloneDisabled, false)
	cmd.Use = "javagen"
	cmd.Short = "Generate Java App Kubernetes resources"
	cmd.Long = `Generate Kubernetes resources as as a kustomize transfomer.
						Implemented as a containerized KRM function.`
	return cmd
}

var fnConfig = &v1alpha1JavaAppGen{}
var processor = buildProcessor(fnConfig)
var rootCmd = createEmbeddedCommand()

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
