package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path"
	"strconv"

	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type ExpandConfiguration struct {
	subscriptionsPath string
	templatesPath     string
	outputPath        string
	recursive         bool
	verbose           bool
}

func GetExpandCommand() components.Command {
	return components.Command{
		Name:        "expand",
		Description: "Expands retention templates",
		Aliases:     []string{},
		Arguments:   GetExpandArguments(),
		Flags:       GetExpandFlags(),
		EnvVars:     GetExpandEnvVar(),
		Action: func(c *components.Context) error {
			return ExpandCmd(c)
		},
	}
}

func GetExpandArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "subscriptions-path",
			Description: "Path to the subscriptions JSON file",
		},
		{
			Name:        "templates-path",
			Description: "Path to the templates dir",
		},
		{
			Name:        "output-path",
			Description: "Path to output the generated filespecs",
		},
	}
}

func GetExpandFlags() []components.Flag {
	return []components.Flag{
		components.BoolFlag{
			Name:         "verbose",
			Description:  "output verbose logging",
			DefaultValue: false,
		},
		components.BoolFlag{
			Name:         "recursive",
			Description:  "recursively find templates in the given dir",
			DefaultValue: false,
		},
	}
}

func GetExpandEnvVar() []components.EnvVar {
	return []components.EnvVar{}
}

func ExpandCmd(context *components.Context) error {
	expandConfig, argErr := ParseExpandConfig(context)
	if argErr != nil {
		return argErr
	}

	if expandConfig.verbose {
		log.Info("expandConfig:")
		log.Info("    subscriptionsPath:", expandConfig.subscriptionsPath)
		log.Info("    templatesPath:", expandConfig.templatesPath)
		log.Info("    outputPath:", expandConfig.outputPath)
		log.Info("    recursive:", expandConfig.recursive)
		log.Info("    verbose:", expandConfig.verbose)
	}

	log.Info("Parsing subscriptions file")
	subscriptionsFile, readErr := os.ReadFile(expandConfig.subscriptionsPath)
	if readErr != nil {
		return readErr
	}

	log.Info("Collecting template files")
	templateFiles, findErr := FindFiles(expandConfig.templatesPath, expandConfig.recursive)
	if findErr != nil {
		return findErr
	}

	if len(templateFiles) == 0 {
		log.Warn("Found no template files")
	} else {
		log.Info("Found", len(templateFiles), "template files")
	}

	if expandConfig.verbose {
		for _, file := range templateFiles {
			log.Info("    " + file)
		}
	}

	var subscriptions map[string][]interface{}
	if jsonErr := json.Unmarshal(subscriptionsFile, &subscriptions); jsonErr != nil {
		return jsonErr
	}

	for subscription, entries := range subscriptions {
		log.Info("Expanding", subscription)

		templateText, readErr := os.ReadFile(path.Join(expandConfig.templatesPath, subscription+".json"))
		if readErr != nil {
			return readErr
		}

		template, parseErr := template.New(subscription).Parse(string(templateText))
		if parseErr != nil {
			return parseErr
		}

		if dirErr := os.MkdirAll(path.Join(expandConfig.outputPath, subscription), 0755); dirErr != nil {
			return dirErr
		}

		for index, entry := range entries {

			resultFile, fileErr := os.Create(path.Join(expandConfig.outputPath, subscription, fmt.Sprint(index, ".json")))
			if fileErr != nil {
				return fileErr
			}

			if templatingErr := template.Execute(resultFile, entry); templatingErr != nil {
				return templatingErr
			}

			if expandConfig.verbose {
				log.Info("    ", resultFile.Name())
			}
		}
	}

	log.Info("Done")
	return nil
}

func ParseExpandConfig(context *components.Context) (*ExpandConfiguration, error) {
	if len(context.Arguments) != 3 {
		return nil, errors.New("Expected 3 argument, received " + strconv.Itoa(len(context.Arguments)))
	}

	var expandConfig = new(ExpandConfiguration)
	expandConfig.subscriptionsPath = context.Arguments[0]
	expandConfig.templatesPath = context.Arguments[1]
	expandConfig.outputPath = context.Arguments[2]
	expandConfig.recursive = context.GetBoolFlagValue("recursive")
	expandConfig.verbose = context.GetBoolFlagValue("verbose")

	return expandConfig, nil
}
