package docker

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		s        string
		platform Platform
		err      error
	}{
		{
			s:   "",
			err: errors.New(`string is malformed: ""`),
		},
		{
			s:   "linux",
			err: errors.New(`string is malformed: "linux"`),
		},
		{
			s:   "linux/",
			err: errors.New(`invalid architecture "" in "linux/"`),
		},
		{
			s:        "linux/amd64",
			platform: Platform{Arch: AMD64},
		},
		{
			s:        "linux/386",
			platform: Platform{Arch: C_386},
		},
		{
			s:        "linux/arm64",
			platform: Platform{Arch: ARM64},
		},
		{
			s:        "linux/s390x",
			platform: Platform{Arch: S390X},
		},
		{
			s:        "linux/ppc64le",
			platform: Platform{Arch: PPC64LE},
		},
		{
			s:   "linux/arm",
			err: errors.New(`arch version is missing from "linux/arm"`),
		},
		{
			s:   "linux/arm/v1",
			err: errors.New(`invalid arch version "v1" in "linux/arm/v1"`),
		},
		{
			s:   "linux/amd64/v1",
			err: errors.New(`too many fields in "linux/amd64/v1"`),
		},
		{
			s:   "linux/bla/v1",
			err: errors.New(`invalid architecture "bla" in "linux/bla/v1"`),
		},
		{
			s:        "linux/arm/v6",
			platform: Platform{Arch: ARM, Version: V6},
		},
		{
			s:        "linux/arm/v7",
			platform: Platform{Arch: ARM, Version: V7},
		},
	}
	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.s, func(t *testing.T) {
			t.Parallel()

			platform, err := Parse(testCase.s)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.platform, platform)
			}
		})
	}
}
