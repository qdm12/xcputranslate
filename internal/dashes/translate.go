package dashes

import "github.com/qdm12/xcputranslate/internal/docker"

func Translate(dockerPlatform docker.Platform) string {
	switch dockerPlatform.Arch {
	case docker.AMD64:
		return "amd64"
	case docker.C386:
		return "x86_64"
	case docker.ARM64:
		return "arm64"
	case docker.ARM:
		if dockerPlatform.Version == docker.V6 {
			return "arm-v6"
		}
		return "arm-v7"
	case docker.S390X:
		return "s390x"
	case docker.PPC64LE:
		return "ppc64le"
	case docker.RISCV64:
		return "riscv64"
	}
	return ""
}
