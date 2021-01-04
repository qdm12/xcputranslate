package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/qdm12/xcputranslate/internal/docker"
	"github.com/qdm12/xcputranslate/internal/golang"
	"github.com/qdm12/xcputranslate/internal/models"
)

//nolint:gochecknoglobals
var (
	version   = "unknown"
	commit    = "unknown"
	buildDate = "an unknown date"
)

var (
	errInvalidField    = errors.New("invalid field requested")
	errInvalidLanguage = errors.New("invalid language requested")
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
	if len(args) > 1 && args[1] == "version" {
		fmt.Printf("🤖 Version %s (commit %s built on %s)\n",
			buildInfo.Version, buildInfo.Commit, buildInfo.BuildDate)
		return nil
	}

	flagSet := flag.NewFlagSet(args[0], flag.ExitOnError)
	languagePtr := flagSet.String("language", "golang", "can be one of: golang")
	fieldPtr := flagSet.String("field", "arch", "can be one of: arch, arm")
	if err := flagSet.Parse(args[1:]); err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)

	outputCh := make(chan string)
	errCh := make(chan error)
	go func() {
		s, err := reader.ReadString('\n')
		if err != nil {
			errCh <- err
			return
		}
		outputCh <- s
	}()

	var input string
	select {
	case <-ctx.Done():
		return ctx.Err()
	case input = <-outputCh:
	case err := <-errCh:
		close(outputCh)
		return err
	}

	input = strings.TrimSuffix(input, "\n")

	language, field := *languagePtr, *fieldPtr

	platform, err := docker.Parse(input)
	if err != nil {
		return err
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
	default:
		return fmt.Errorf("%w: %q", errInvalidLanguage, language)
	}

	fmt.Println(output)

	return nil
}
