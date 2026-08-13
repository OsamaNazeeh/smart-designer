package api

import (
	"fmt"
	"os/exec"
)



func applyTransform(options ...string) error {
	cmd := exec.Command("ffmpeg", options...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w\n%s", err, output)
	}
	return nil
}