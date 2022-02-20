package golang

import (
	"fmt"
	"os"
)

func SetEnv(arch, arm string) (err error) {
	err = os.Setenv("GOARCH", arch)
	if err != nil {
		return fmt.Errorf("cannot set GOARCH=%s: %w", arch, err)
	}

	err = os.Setenv("GOARM", arm)
	if err != nil {
		return fmt.Errorf("cannot set GOARM=%s: %w", arm, err)
	}

	return nil
}
