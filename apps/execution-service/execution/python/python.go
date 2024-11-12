package python

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func RunPythonCode(code string, input string) (string, string, error) {
	cmd := exec.Command(
		"docker", "run", "--rm",
		"-i", // allows for standard input to be passed in
		"apps_python-sandbox",
		"python", "-c", code,
	)

	// Pass in any input data to the Python script
	cmd.Stdin = bytes.NewBufferString(input)

	// Capture output and error
	var output bytes.Buffer
	var errorOutput bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &errorOutput

	// Run the command
	if err := cmd.Run(); err != nil {
		return "", fmt.Sprintf("Command execution failed: %s: %v", errorOutput.String(), err), nil
	}

	return strings.TrimSuffix(output.String(), "\n"), strings.TrimSuffix(errorOutput.String(), "\n"), nil
}
