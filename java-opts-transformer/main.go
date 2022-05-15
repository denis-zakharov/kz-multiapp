package main

import (
	"os"
	"strconv"

	"github.com/denis-zakharov/kz-multiapp/java-opts-transformer/cmd"
	"sigs.k8s.io/kustomize/kyaml/fn/framework/command"
)

func main() {
	debugVar := os.Getenv("FN_STANDALONE")
	mode := map[bool]command.CLIMode{
		false: command.StandaloneDisabled,
		true:  command.StandaloneEnabled,
	}

	debug, err := strconv.ParseBool(debugVar)
	if err != nil {
		cmd.Execute(command.StandaloneDisabled)
	} else {
		cmd.Execute(mode[debug])
	}
}
