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

func isDeployment(rn *yaml.RNode) bool {
	kind := rn.GetKind()
	return kind == "Deployment" || kind == "StatefulSet"
}

func (trans *JavaOptsTransformerv1alpha1) Filter(items []*yaml.RNode) ([]*yaml.RNode, error) {
	for _, item := range items {
		if !isDeployment(item) || item.GetName() != trans.Metadata.Name {
			continue
		}
		err := item.PipeE(injectJavaOpts(trans.Spec))
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

type JavaOptsSetter struct {
	MemLimit string
	Value    string
	Index    int
}

func injectJavaOpts(spec Spec) JavaOptsSetter {
	memLimitNum := int(float64(spec.Heap) / spec.Ratio)
	memLimit := fmt.Sprintf("%dMi", memLimitNum)
	jvmOptsValue := fmt.Sprintf("-Xms%dm -Xmx%dm %s", spec.Heap, spec.Heap, spec.Extra)
	return JavaOptsSetter{memLimit, jvmOptsValue, spec.Index}
}

func (s JavaOptsSetter) Filter(item *yaml.RNode) (*yaml.RNode, error) {
	jvmOpts, err := yaml.Parse(`
name: JAVA_OPTIONS
value: ` + s.Value)
	if err != nil {
		return item, err
	}

	// remove existing JAVA_OPTIONS var
	err = item.PipeE(
		yaml.LookupCreate(yaml.SequenceNode,
			"spec", "template", "spec", "containers"),
		yaml.GetElementByIndex(s.Index),
		yaml.LookupCreate(yaml.SequenceNode, "env"),
		yaml.ElementSetter{
			Keys:    []string{"name"},
			Values:  []string{"JAVA_OPTIONS"},
			Element: nil,
		},
	)
	if err != nil {
		return item, err
	}
	// append JAVA_OPTIONS to the env spec end
	err = item.PipeE(
		yaml.LookupCreate(yaml.SequenceNode,
			"spec", "template", "spec", "containers"),
		yaml.GetElementByIndex(s.Index),
		yaml.LookupCreate(yaml.SequenceNode, "env"),
		yaml.ElementAppender{Elements: []*yaml.Node{jvmOpts.YNode()}},
	)
	if err != nil {
		return item, err
	}

	// requests
	err = item.PipeE(
		yaml.LookupCreate(yaml.SequenceNode,
			"spec", "template", "spec", "containers"),
		yaml.GetElementByIndex(s.Index),
		yaml.LookupCreate(yaml.MappingNode, "resources", "requests"),
		yaml.SetField("memory", yaml.NewStringRNode(s.MemLimit)),
	)
	if err != nil {
		return item, err
	}

	// limits
	err = item.PipeE(
		yaml.LookupCreate(yaml.SequenceNode,
			"spec", "template", "spec", "containers"),
		yaml.GetElementByIndex(s.Index),
		yaml.LookupCreate(yaml.MappingNode, "resources", "limits"),
		yaml.SetField("memory", yaml.NewStringRNode(s.MemLimit)),
	)
	if err != nil {
		return item, err
	}

	return item, nil
}
