package uname

import "github.com/qdm12/xcputranslate/internal/docker"

func Translate(dockerPlatform docker.Platform) (arch string) {
	switch dockerPlatform.Arch {
	case docker.AMD64, docker.C386:
		return "x86_64"
	case docker.ARM64:
		return "aarch64"
	case docker.ARM:
		return "armv7l" // version does not matter for `uname -m`
	case docker.S390X:
		return "s390x"
	case docker.PPC64LE:
		return "ppc64le"
	case docker.RISCV64:
		return "riscv64"
	}
	return ""
}
