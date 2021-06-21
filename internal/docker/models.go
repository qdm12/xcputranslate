package docker

type Arch string

type ArchVersion string

type Platform struct {
	Arch    Arch
	Version ArchVersion
}

func (p *Platform) Equal(other Platform) bool {
	if p.Arch == ARM64 && other.Arch == ARM64 { // special case for ARM64v8
		return true
	}
	return p.Arch == other.Arch && p.Version == other.Version
}
