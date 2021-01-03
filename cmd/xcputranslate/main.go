package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/qdm12/xcputranslate/internal/docker"
	"github.com/qdm12/xcputranslate/internal/golang"
)

func main() {
	os.Exit(_main(os.Args))
}

func _main(args []string) int {
	flagSet := flag.NewFlagSet(args[0], flag.ExitOnError)
	languagePtr := flagSet.String("language", "golang", "can be one of: golang")
	fieldPtr := flagSet.String("field", "arch", "can be one of: arch, arm")
	if err := flagSet.Parse(args[1:]); err != nil {
		fmt.Println(err)
		return 1
	}

	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return 1
	}

	s = strings.TrimSuffix(s, "\n")

	output, err := getOutput(s, *languagePtr, *fieldPtr)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	fmt.Println(output)

	return 0
}

var (
	errInvalidField    = errors.New("invalid field requested")
	errInvalidLanguage = errors.New("invalid language requested")
)

func getOutput(s, language, field string) (output string, err error) {
	platform, err := docker.Parse(s)
	if err != nil {
		return "", err
	}

	switch language {
	case "golang":
		arch, arm := golang.Translate(platform)
		switch field {
		case "arch":
			return arch, nil
		case "arm":
			return arm, nil
		default:
			return "", fmt.Errorf("%w: %q", errInvalidField, field)
		}
	default:
		return "", fmt.Errorf("%w: %q", errInvalidLanguage, language)
	}
}
