package alert

import (
	"os/exec"
)

// Say makes use of the say command to speak
// the arguments out load.
//
// This is useful when running on pc in the
// background so that one is aware when
// something happens.
func Say(args ...string) error {
	cmd := exec.Command("/usr/bin/say", args...)
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
