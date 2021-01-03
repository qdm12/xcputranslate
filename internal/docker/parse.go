package docker

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMalformed          = errors.New("string is malformed")
	ErrNotEnoughFields    = errors.New("not enough fields")
	ErrArchVersionMissing = errors.New("arch version is missing")
	ErrInvalidArch        = errors.New("invalid architecture")
	ErrTooManyFields      = errors.New("too many fields")
	ErrBadArchVersion     = errors.New("invalid arch version")
)

func Parse(s string) (platform Platform, err error) {
	if !strings.HasPrefix(s, "linux/") {
		return platform, fmt.Errorf("%w: %q", ErrMalformed, s)
	}
	elements := strings.Split(s, "/")
	switch len(elements) {
	case 2:
		platform.Arch = Arch(elements[1])
		switch platform.Arch {
		case AMD64, C386, ARM64, S390X, PPC64LE:
		case ARM:
			return platform, fmt.Errorf("%w from %q", ErrArchVersionMissing, s)
		default:
			return platform, fmt.Errorf("%w %q in %q", ErrInvalidArch, elements[1], s)
		}
	case 3:
		platform.Arch = Arch(elements[1])
		switch platform.Arch {
		case AMD64, C386, ARM64, S390X, PPC64LE:
			return platform, fmt.Errorf("%w in %q", ErrTooManyFields, s)
		case ARM:
			platform.Version = ArchVersion(elements[2])
			switch platform.Version {
			case V6, V7:
			default:
				return platform, fmt.Errorf("%w %q in %q", ErrBadArchVersion, elements[2], s)
			}
		default:
			return platform, fmt.Errorf("%w %q in %q", ErrInvalidArch, elements[1], s)
		}
	}
	return platform, nil
}
