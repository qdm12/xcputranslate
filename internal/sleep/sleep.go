package sleep

import (
	"context"
	"time"

	"github.com/qdm12/xcputranslate/internal/docker"
)

func Sleep(ctx context.Context, dockerPlatform docker.Platform,
	order []docker.Platform, buildTime time.Duration) error {
	duration := getSleepDuration(dockerPlatform, order, buildTime)
	if duration == 0 {
		return nil
	}
	timer := time.NewTimer(duration)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
	case <-timer.C:
	}

	return ctx.Err()
}

func getSleepDuration(dockerPlatform docker.Platform, order []docker.Platform,
	buildTime time.Duration) (sleep time.Duration) {
	if len(order) == 0 {
		order = []docker.Platform{ // sorted by popularity
			{Arch: docker.AMD64},
			{Arch: docker.C386},
			{Arch: docker.ARM64},
			{Arch: docker.ARM, Version: docker.V6},
			{Arch: docker.ARM, Version: docker.V7},
			{Arch: docker.PPC64LE},
			{Arch: docker.S390X},
			{Arch: docker.RISCV64},
		}
	}

	platformIndex := -1
	for i, platform := range order {
		if dockerPlatform.Equal(platform) {
			platformIndex = i
			break
		}
	}

	sleep = buildTime * time.Duration(platformIndex)
	return sleep
}
