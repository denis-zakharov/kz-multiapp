package main

import (
	"fmt"
	"io"
	"os"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func loadYAML(path string) *yaml.RNode {
	f, err := os.Open(path)
	must(err)
	b, err := io.ReadAll(f)
	must(err)
	return parseYAML(string(b))
}

func parseYAML(data string) *yaml.RNode {
	r, err := yaml.Parse(data)
	must(err)
	return r
}

func processYAML(r *yaml.RNode) {

	jvmOpts := parseYAML(`
name: JAVA_OPTIONS
value: SUPER_HEAP
`)

	// remove
	err := r.PipeE(
		yaml.LookupCreate(yaml.SequenceNode,
			"spec", "template", "spec", "containers"),
		yaml.GetElementByIndex(0),
		yaml.LookupCreate(yaml.SequenceNode, "env"),
		yaml.ElementSetter{
			Keys:    []string{"name"},
			Values:  []string{"JAVA_OPTIONS"},
			Element: nil,
		},
	)
	must(err)
	// append
	err = r.PipeE(
		yaml.LookupCreate(yaml.SequenceNode,
			"spec", "template", "spec", "containers"),
		yaml.GetElementByIndex(0),
		yaml.LookupCreate(yaml.SequenceNode, "env"),
		yaml.ElementAppender{Elements: []*yaml.Node{jvmOpts.YNode()}},
	)
	must(err)

	fmt.Println(r.String())
}

func main() {
	for _, v := range os.Args[1:] {
		r := loadYAML(v)
		processYAML(r)
	}
}
