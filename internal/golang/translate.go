package golang

import "github.com/qdm12/xcputranslate/internal/docker"

func Translate(dockerPlatform docker.Platform) (arch, arm string) {
	arch = string(dockerPlatform.Arch)
	if dockerPlatform.Arch == docker.ARM {
		switch dockerPlatform.Version { //nolint:exhaustive
		case docker.V6:
			arm = "6"
		case docker.V7:
			arm = "7"
		}
	}
	return arch, arm
}
