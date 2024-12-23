package tools

import (
	"fmt"
	"os/exec"
)

func CheckCleverAuth() error {
	cmd := exec.Command("clever", "profile")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("you should be logged in Clever Tools to convert your access logs")
	}
	return nil
}
