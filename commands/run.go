package commands

import (
	"errors"
	"strconv"

	core_utils "github.com/jfrog/jfrog-cli-core/v2/artifactory/utils"
	core_commands "github.com/jfrog/jfrog-cli-core/v2/common/commands"
	core_components "github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	core_config "github.com/jfrog/jfrog-cli-core/v2/utils/config"
	client_artifactory "github.com/jfrog/jfrog-client-go/artifactory"
	client_utils "github.com/jfrog/jfrog-client-go/utils"
	client_log "github.com/jfrog/jfrog-client-go/utils/log"
)

type RunConfiguration struct {
	fileSpecsPath string
	dryRun        bool
	recursive     bool
	verbose       bool
}

func GetRunCommand() core_components.Command {
	return core_components.Command{
		Name:        "run",
		Description: "Runs retention",
		Aliases:     []string{},
		Arguments:   GetRunArguments(),
		Flags:       GetRunFlags(),
		EnvVars:     GetRunEnvVar(),
		Action: func(c *core_components.Context) error {
			return RunCmd(c)
		},
	}
}

func GetRunArguments() []core_components.Argument {
	return []core_components.Argument{
		{
			Name:        "filespecs",
			Description: "Path to the filespecs file/dir",
		},
	}
}

func GetRunFlags() []core_components.Flag {
	return []core_components.Flag{
		core_components.BoolFlag{
			Name:         "dry-run",
			Description:  "disable communication with Artifactory",
			DefaultValue: true,
		},
		core_components.BoolFlag{
			Name:         "verbose",
			Description:  "output verbose logging",
			DefaultValue: false,
		},
		core_components.BoolFlag{
			Name:         "recursive",
			Description:  "recursively find filespecs files in the given dir",
			DefaultValue: false,
		},
	}
}

func GetRunEnvVar() []core_components.EnvVar {
	return []core_components.EnvVar{}
}

func RunCmd(context *core_components.Context) error {
	runConfig, err := ParseRunConfig(context)
	if err != nil {
		return err
	}

	if runConfig.verbose {
		client_log.Info("runConfig:")
		client_log.Info("    fileSpecsPath:", runConfig.fileSpecsPath)
		client_log.Info("    dryRun:", runConfig.dryRun)
		client_log.Info("    recursive:", runConfig.recursive)
		client_log.Info("    verbose:", runConfig.verbose)
	}

	client_log.Info("Fetching Artifactory details")
	artifactoryDetails, err := GetArtifactoryDetails(context)
	if err != nil {
		return err
	}

	client_log.Info("Configuring Artifactory manager")
	artifactoryManager, err := core_utils.CreateServiceManager(artifactoryDetails, 3, 5000, runConfig.dryRun)
	if err != nil {
		return err
	}

	client_log.Info("Parsing retention configuration")
	fileSpecsFiles, err := FindFiles(runConfig.fileSpecsPath, runConfig.recursive)

	client_log.Info("Running", len(fileSpecsFiles), "policies")
	if runConfig.verbose {
		for _, file := range fileSpecsFiles {
			client_log.Info("    " + file)
		}
	}

	client_log.Info("Running", len(fileSpecsFiles), "policies")
	if err = RunArtifactRetention(artifactoryManager, fileSpecsFiles); err != nil {
		return err
	}

	client_log.Info("Done")
	return nil
}

func ParseRunConfig(context *core_components.Context) (*RunConfiguration, error) {
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

func GetArtifactoryDetails(c *core_components.Context) (*core_config.ServerDetails, error) {
	details, err := core_commands.GetConfig("", false)
	if err != nil {
		return nil, err
	}

	if details.Url == "" {
		return nil, errors.New("no server-id was found, or the server-id has no url")
	}

	details.Url = client_utils.AddTrailingSlashIfNeeded(details.Url)
	err = core_config.CreateInitialRefreshableTokensIfNeeded(details)
	if err != nil {
		return nil, err
	}

	return details, nil
}

func RunArtifactRetention(artifactoryManager client_artifactory.ArtifactoryServicesManager, fileSpecsFiles []string) error {
	totalFiles := len(fileSpecsFiles)
	for i, file := range fileSpecsFiles {
		client_log.Info(i+1, "/", totalFiles, ":", file)

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
