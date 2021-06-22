package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/qdm12/xcputranslate/internal/dashes"
	"github.com/qdm12/xcputranslate/internal/docker"
	"github.com/qdm12/xcputranslate/internal/golang"
	"github.com/qdm12/xcputranslate/internal/models"
	"github.com/qdm12/xcputranslate/internal/sleep"
	"github.com/qdm12/xcputranslate/internal/uname"
)

//nolint:gochecknoglobals
var (
	version   = "unknown"
	commit    = "unknown"
	buildDate = "an unknown date"
)

var (
	errNoCommand       = errors.New("no command specified")
	errInvalidCommand  = errors.New("invalid command")
	errInvalidField    = errors.New("invalid field requested")
	errInvalidLanguage = errors.New("invalid language requested")
	errInvalidPlatform = errors.New("invalid platform")
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	buildInfo := models.BuildInfo{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
	}

	errorCh := make(chan error)
	go func() {
		errorCh <- _main(ctx, os.Args, buildInfo)
	}()

	signalsCh := make(chan os.Signal, 1)
	signal.Notify(signalsCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
	)

	exitCode := 0
	select {
	case signal := <-signalsCh:
		fmt.Println("\nShutting down: signal", signal)
		exitCode = 1
		cancel()
		timer := time.NewTimer(time.Second)
		select {
		case <-errorCh:
			if !timer.Stop() {
				<-timer.C
			}
		case <-timer.C:
			fmt.Println("Shutdown timed out")
		}
	case err := <-errorCh:
		if err != nil {
			fmt.Println("Fatal error:", err)
			exitCode = 1
		}
		cancel()
	}
	os.Exit(exitCode)
}

func _main(ctx context.Context, args []string, buildInfo models.BuildInfo) error {
	if len(args) == 1 {
		return fmt.Errorf("%w: can be one of: version, translate", errNoCommand)
	}

	switch args[1] {
	case "version":
		fmt.Printf("🤖 Version %s (commit %s built on %s)\n",
			buildInfo.Version, buildInfo.Commit, buildInfo.BuildDate)
		return nil
	case "translate":
		return translate(args)
	case "sleep":
		return clisleep(ctx, args)
	default:
		return fmt.Errorf("%w: %s", errInvalidCommand, args[1])
	}
}

func translate(args []string) (err error) {
	flagSet := flag.NewFlagSet(args[1], flag.ExitOnError)
	languagePtr := flagSet.String("language", "golang", "can be one of: golang, uname, dashes")
	fieldPtr := flagSet.String("field", "", "required for golang and can be one of: arch, arm")
	targetPlatformPtr := flagSet.String("targetplatform", "", "can be for example linux/arm64")
	if err := flagSet.Parse(args[2:]); err != nil {
		return err
	}

	language, field, targetPlatform := *languagePtr, *fieldPtr, *targetPlatformPtr

	platform, err := docker.Parse(targetPlatform)
	if err != nil {
		return fmt.Errorf("%w: %s", errInvalidPlatform, err)
	}

	var output string
	switch language {
	case "golang":
		arch, arm := golang.Translate(platform)
		switch field {
		case "arch":
			output = arch
		case "arm":
			output = arm
		default:
			return fmt.Errorf("%w: %q", errInvalidField, field)
		}
	case "uname":
		output = uname.Translate(platform)
	case "dashes":
		output = dashes.Translate(platform)
	default:
		return fmt.Errorf("%w: %q", errInvalidLanguage, language)
	}

	fmt.Println(output)

	return nil
}

func clisleep(ctx context.Context, args []string) (err error) {
	flagSet := flag.NewFlagSet(args[1], flag.ExitOnError)
	targetPlatformPtr := flagSet.String("targetplatform", "", "can be for example linux/arm64")
	const defaultOrder = "linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6,linux/386,linux/ppc64le,linux/s390x,linux/riscv64" //nolint:lll
	orderPtr := flagSet.String("order", defaultOrder,
		"order of CPU architectures to build. Use this to reduce the sleeping times.")
	const defaultBuiltTime = 3 * time.Second
	buildTimePtr := flagSet.Duration("buildtime", defaultBuiltTime, "approximate build time")
	if err := flagSet.Parse(args[2:]); err != nil {
		return err
	}

	buildTime, orderString, targetPlatformString := *buildTimePtr, *orderPtr, *targetPlatformPtr

	targetPlatform, err := docker.Parse(targetPlatformString)
	if err != nil {
		return fmt.Errorf("%w: %s", errInvalidPlatform, err)
	}

	orderPlatforms := strings.Split(orderString, ",")
	order := make([]docker.Platform, len(orderPlatforms))
	for i, s := range orderPlatforms {
		platform, err := docker.Parse(s)
		if err != nil {
			return fmt.Errorf("%w: in order at position %d: %s", errInvalidPlatform, i, err)
		}
		order[i] = platform
	}

	return sleep.Sleep(ctx, targetPlatform, order, buildTime)
}
