package node

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func RunJavaScriptCode(code string, input string) (string, string, error) {
	cmd := exec.Command(
		"docker", "run", "--rm",
		"-i",                // allows standard input to be passed in
		"apps-node-sandbox", // Docker image with Node.js environment
		"node", "-e", code,  // Runs JavaScript code with Node.js
	)

	// Pass input to the JavaScript script
	cmd.Stdin = bytes.NewBufferString(input)

	// Capture standard output and error output
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
