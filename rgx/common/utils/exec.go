package utils

import (
	"fmt"
	"os/exec"
	"rgx/common/log"
	"strings"
)

func RunScript(scriptCmd, scriptDir string, envMap map[string]string) {
	var command string
	var args []string
	if strings.HasSuffix(scriptCmd, ".cmd") {
		command = "cmd"
		args = []string{"/c", scriptCmd}
	} else {
		command = "sh"
		args = []string{scriptCmd}
	}
	cmd := exec.Command(command, args...)
	log.Trace("setting script directory to %s", scriptDir)
	cmd.Dir = scriptDir
	m := cmd.Environ()
	for k, v := range envMap {
		m = append(m, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = m
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("could not run command '%s': %s", scriptCmd, err.Error())
	}
	fmt.Println(string(output))
}
