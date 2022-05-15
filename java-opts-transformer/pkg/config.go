package pkg

import (
	"fmt"

	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type JavaOptsTransformerv1alpha1 struct {
	Metadata Metadata `yaml:"metadata"`
	Spec     Spec     `yaml:"spec,omitempty"`
}
type Metadata struct {
	Name string `yaml:"name"`
}
type Spec struct {
	Index int     `yaml:"index"`
	Heap  int     `yaml:"heap"`
	Ratio float64 `yaml:"heap"`
	Extra string  `yaml:"extra"`
}

func (trans *JavaOptsTransformerv1alpha1) Validate() error {
	if trans.Metadata.Name == "" {
		return fmt.Errorf("metadata.name: must be set to target Deployment/StatefulSet")
	}
	if trans.Spec.Index < 0 {
		return fmt.Errorf("container index: must be non-negative integer")
	}
	if trans.Spec.Heap <= 0 {
		return fmt.Errorf("heap size (Mb): must be positive integer")
	}
	if trans.Spec.Ratio < 0.5 && trans.Spec.Ratio > 0.9 {
		return fmt.Errorf("heap to mem limit ratio: must be positive fraction [0.5, 0.9]")
	}
	return nil
}

func (trans *JavaOptsTransformerv1alpha1) Filter(items []*yaml.RNode) ([]*yaml.RNode, error) {
	for i := range items {
		err := items[i].PipeE(yaml.SetAnnotation("custom.io/the-value", trans.Metadata.Name))
		if err != nil {
			return nil, err
		}
	}
	return items, nil
}

func (trans *JavaOptsTransformerv1alpha1) BuildProcessor() framework.ResourceListProcessor {
	return framework.SimpleProcessor{
		Config: trans,
		Filter: trans,
	}
}
