package cmd

import "sigs.k8s.io/kustomize/kyaml/fn/framework/command"

func init() {
	command.AddGenerateDockerfile(rootCmd)
}
