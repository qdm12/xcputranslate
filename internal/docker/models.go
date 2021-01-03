package docker

type Arch string

type ArchVersion string

type Platform struct {
	Arch    Arch
	Version ArchVersion
}
