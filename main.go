package main

import (
	"bufio"
	"fmt"
	"os"

	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type Spec struct {
	Value string `yaml:"value,omitempty"`
}
type Example struct {
	Spec Spec `yaml:"spec,omitempty"`
}

func runFunction(rlSource *kio.ByteReadWriter) error {
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

	p := framework.SimpleProcessor{Config: functionConfig, Filter: kio.FilterFunc(fn)}
	err := framework.Execute(p, rlSource)
	return fmt.Errorf("runFn: %w", err)
}

func main() {
	f, err := os.Open("input.yaml")
	if err != nil {
		panic(err)
	}

	input := bufio.NewReader(f)

	runFunction(&kio.ByteReadWriter{Reader: input})
}
