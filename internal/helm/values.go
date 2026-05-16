package helm

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

func ShowValues(chart string) ([]byte, error) {
	cmd := exec.Command("helm", "show", "values", chart)
	return runHelmValuesCommand(cmd, "helm show values")
}

func GetReleaseValues(release string, namespace string) ([]byte, error) {
	args := []string{"get", "values", release, "--all"}
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	cmd := exec.Command("helm", args...)
	return runHelmValuesCommand(cmd, "helm get values")
}

func runHelmValuesCommand(cmd *exec.Cmd, label string) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return nil, fmt.Errorf("helm is not installed or not found in PATH")
		}
		msg := stderr.String()
		if msg == "" {
			msg = err.Error()
		}
		return nil, fmt.Errorf("%s failed: %s", label, msg)
	}

	return stdout.Bytes(), nil
}
