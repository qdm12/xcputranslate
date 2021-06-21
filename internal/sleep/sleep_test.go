package sleep

import (
	"testing"
	"time"

	"github.com/qdm12/xcputranslate/internal/docker"
	"github.com/stretchr/testify/assert"
)

func Test_getSleepDuration(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		platform  docker.Platform
		order     []docker.Platform
		buildTime time.Duration
		sleep     time.Duration
	}{
		"amd64 default order": {
			platform:  docker.Platform{Arch: docker.AMD64},
			buildTime: time.Second,
		},
		"armv8 default order": {
			platform:  docker.Platform{Arch: docker.ARM64, Version: docker.V8},
			buildTime: time.Second,
			sleep:     2 * time.Second,
		},
		"amd64 custom order": {
			platform:  docker.Platform{Arch: docker.AMD64},
			order:     []docker.Platform{{}, {}, {Arch: docker.AMD64}, {}},
			buildTime: time.Second,
			sleep:     2 * time.Second,
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			sleep := getSleepDuration(testCase.platform, testCase.order, testCase.buildTime)

			assert.Equal(t, testCase.sleep, sleep)
		})
	}
}
