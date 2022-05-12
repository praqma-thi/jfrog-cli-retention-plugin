package commands

import (
	"errors"
	"strconv"

	core_utils "github.com/jfrog/jfrog-cli-core/v2/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/v2/common/commands"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-client-go/artifactory"
	client_utils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type RunConfiguration struct {
	fileSpecsPath string
	dryRun        bool
	recursive     bool
	verbose       bool
}

func GetRunCommand() components.Command {
	return components.Command{
		Name:        "run",
		Description: "Runs retention",
		Aliases:     []string{},
		Arguments:   GetRunArguments(),
		Flags:       GetRunFlags(),
		EnvVars:     GetRunEnvVar(),
		Action: func(c *components.Context) error {
			return RunCmd(c)
		},
	}
}

func GetRunArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "filespecs-path",
			Description: "Path to the filespecs file/dir",
		},
	}
}

func GetRunFlags() []components.Flag {
	return []components.Flag{
		components.BoolFlag{
			Name:         "dry-run",
			Description:  "disable deletion of artifacts",
			DefaultValue: true,
		},
		components.BoolFlag{
			Name:         "verbose",
			Description:  "output verbose logging",
			DefaultValue: false,
		},
		components.BoolFlag{
			Name:         "recursive",
			Description:  "recursively find filespecs files in the given dir",
			DefaultValue: false,
		},
	}
}

func GetRunEnvVar() []components.EnvVar {
	return []components.EnvVar{}
}

func RunCmd(context *components.Context) error {
	runConfig, err := ParseRunConfig(context)
	if err != nil {
		return err
	}

	if runConfig.verbose {
		log.Info("runConfig:")
		log.Info("    fileSpecsPath:", runConfig.fileSpecsPath)
		log.Info("    dryRun:", runConfig.dryRun)
		log.Info("    recursive:", runConfig.recursive)
		log.Info("    verbose:", runConfig.verbose)
	}

	log.Info("Fetching Artifactory details")
	artifactoryDetails, err := GetArtifactoryDetails(context)
	if err != nil {
		return err
	}

	log.Info("Configuring Artifactory manager")
	artifactoryManager, err := core_utils.CreateServiceManager(artifactoryDetails, 3, 5000, runConfig.dryRun)
	if err != nil {
		return err
	}

	log.Info("Parsing retention configuration")
	fileSpecsFiles, err := FindFiles(runConfig.fileSpecsPath, runConfig.recursive)

	if len(fileSpecsFiles) == 0 {
		log.Warn("Found no FileSpec files")
	} else {
		log.Info("Found", len(fileSpecsFiles), "FileSpec files")
	}

	if runConfig.verbose {
		for _, file := range fileSpecsFiles {
			log.Info("    " + file)
		}
	}

	if err = RunArtifactRetention(artifactoryManager, fileSpecsFiles); err != nil {
		return err
	}

	log.Info("Done")
	return nil
}

func ParseRunConfig(context *components.Context) (*RunConfiguration, error) {
	if len(context.Arguments) != 1 {
		return nil, errors.New("Expected 1 argument, received " + strconv.Itoa(len(context.Arguments)))
	}

	var runConfig = new(RunConfiguration)
	runConfig.fileSpecsPath = context.Arguments[0]
	runConfig.dryRun = context.GetBoolFlagValue("dry-run")
	runConfig.recursive = context.GetBoolFlagValue("recursive")
	runConfig.verbose = context.GetBoolFlagValue("verbose")

	return runConfig, nil
}

func GetArtifactoryDetails(c *components.Context) (*config.ServerDetails, error) {
	details, err := commands.GetConfig("", false)
	if err != nil {
		return nil, err
	}

	if details.Url == "" {
		return nil, errors.New("no server-id was found, or the server-id has no url")
	}

	details.Url = client_utils.AddTrailingSlashIfNeeded(details.Url)
	err = config.CreateInitialRefreshableTokensIfNeeded(details)
	if err != nil {
		return nil, err
	}

	return details, nil
}

func RunArtifactRetention(artifactoryManager artifactory.ArtifactoryServicesManager, fileSpecsFiles []string) error {
	totalFiles := len(fileSpecsFiles)
	for i, file := range fileSpecsFiles {
		log.Info(i+1, "/", totalFiles, ":", file)

		deleteParams, err := ParseDeleteParams(file)
		if err != nil {
			return err
		}

		for _, dp := range deleteParams {
			pathsToDelete, err := artifactoryManager.GetPathsToDelete(dp)
			if err != nil {
				return err
			}
			defer pathsToDelete.Close()

			artifactoryManager.DeleteFiles(pathsToDelete)
		}
	}

	return nil
}
