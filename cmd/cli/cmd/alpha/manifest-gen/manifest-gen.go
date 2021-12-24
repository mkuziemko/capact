package manifestgen

import (
	"strings"

	"capact.io/capact/cmd/cli/cmd/alpha/manifest-gen/attribute"
	"capact.io/capact/cmd/cli/cmd/alpha/manifest-gen/common"
	"capact.io/capact/cmd/cli/cmd/alpha/manifest-gen/implementation"
	_interface "capact.io/capact/cmd/cli/cmd/alpha/manifest-gen/interface"
	"capact.io/capact/internal/cli/alpha/manifestgen"
	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/utils/strings/slices"
)

type getManifestFun func(cfg *manifestgen.InterfaceConfig) (map[string]string, error)

// NewCmd returns a cobra.Command for content generation operations.
func NewCmd() *cobra.Command {
	var opts common.ManifestGenOptions
	cmd := &cobra.Command{
		Use:   "manifest-gen",
		Short: "Manifests generation",
		Long:  "Subcommand for various manifest generation operations",
		Args:  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return askInteractivelyForParameters(opts)
		},
	}

	cmd.AddCommand(_interface.NewInterface())
	cmd.AddCommand(implementation.NewCmd())

	cmd.PersistentFlags().StringP("output", "o", "generated", "Path to the output directory for the generated manifests")
	cmd.PersistentFlags().Bool("overwrite", false, "Overwrite existing manifest files")

	return cmd
}

func askInteractivelyForParameters(opts common.ManifestGenOptions) error {
	var err error
	opts.ManifestsType, err = askForManifestType()
	if err != nil {
		return errors.Wrap(err, "while asking for manifest type")
	}

	opts.Directory, err = common.AskForDirectory("path to the output directory for the generated manifests", "generated")
	if err != nil {
		return errors.Wrap(err, "while asking for output directory")
	}

	opts.Overwrite, err = askIfOverwrite()
	if err != nil {
		return errors.Wrap(err, "while asking if overwrite existing manifest files")
	}

	opts.ManifestPath, err = askForManifestPathSuffix()
	if err != nil {
		return errors.Wrap(err, "while asking for manifest path suffix")
	}

	metadata, err := askForCommonMetadataInformation()
	if err != nil {
		return errors.Wrap(err, "while getting the common metadata information")
	}
	opts.Metadata = *metadata

	var mergeFiles map[string]string
	if slices.Contains(opts.ManifestsType, common.InterfaceType) {
		if slices.Contains(opts.ManifestsType, common.TypeType) {
			opts.TypeInputPath = common.CreateManifestPath(common.TypeType, opts.ManifestPath) + "-input:0.0.1"
			opts.TypeOutputPath = common.CreateManifestPath(common.TypeType, opts.ManifestPath) + "config:0.0.1"
		}
		files, err := GenerateInterfaceFile(opts, manifestgen.GenerateInterfaceTemplatingConfig)
		if err != nil {
			return errors.Wrap(err, "while generating interface templating config")
		}
		mergeFiles = MergeMaps(mergeFiles, files)
	}

	if slices.Contains(opts.ManifestsType, common.InterfaceGroupType) {
		files, err := GenerateInterfaceFile(opts, manifestgen.GenerateInterfaceGroupTemplatingConfig)
		if err != nil {
			return errors.Wrap(err, "while generating interface group templating config")
		}
		mergeFiles = MergeMaps(mergeFiles, files)
	}

	if slices.Contains(opts.ManifestsType, common.TypeType) {
		files, err := GenerateInterfaceFile(opts, manifestgen.GenerateTypeTemplatingConfig)
		if err != nil {
			return errors.Wrap(err, "while generating type templating config")
		}
		mergeFiles = MergeMaps(mergeFiles, files)
	}

	if slices.Contains(opts.ManifestsType, common.AttributeType) {
		files, err := attribute.GenerateAttributeFile(opts)
		if err != nil {
			return errors.Wrap(err, "while generating type templating config")
		}
		mergeFiles = MergeMaps(mergeFiles, files)
	}

	if slices.Contains(opts.ManifestsType, common.ImplementationType) {
		files, err := implementation.HandleInteractiveSession(opts)
		if err != nil {
			return errors.Wrap(err, "while generating implementation tool")
		}
		mergeFiles = MergeMaps(mergeFiles, files)
	}

	if err := manifestgen.WriteManifestFiles(opts.Directory, mergeFiles, opts.Overwrite); err != nil {
		return errors.Wrap(err, "while writing manifest files")
	}
	return nil
}

func MergeMaps(maps ...map[string]string) (result map[string]string) {
	result = make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func GenerateInterfaceFile(opts common.ManifestGenOptions, fn getManifestFun) (map[string]string, error) {
	var interfaceCfg manifestgen.InterfaceConfig
	interfaceCfg.ManifestPath = common.CreateManifestPath(common.InterfaceType, opts.ManifestPath)
	interfaceCfg.ManifestMetadata = opts.Metadata
	interfaceCfg.InputPathWithRevision = opts.TypeInputPath
	interfaceCfg.OutputPathWithRevision = opts.TypeOutputPath
	files, err := fn(&interfaceCfg)
	if err != nil {
		return nil, errors.Wrap(err, "while generating content files")
	}
	return files, nil
}

func askForManifestType() ([]string, error) {
	var manifestTypes []string
	availableManifestsType := []string{common.InterfaceGroupType, common.InterfaceType, common.ImplementationType, common.AttributeType, common.TypeType}
	prompt := []*survey.Question{
		{
			Prompt: &survey.MultiSelect{
				Message: "Which manifests do you want to generate:",
				Options: availableManifestsType,
			},
			Validate: survey.MinItems(1),
		},
	}
	err := survey.Ask(prompt, &manifestTypes)
	return manifestTypes, err
}

func askForCommonMetadataInformation() (*common.Metadata, error) {
	var metadata common.Metadata
	var qs = []*survey.Question{
		{
			Name: "DocumentationURL",
			Prompt: &survey.Input{
				Message: "What is documentation URL?",
				Default: "https://example.com",
			},
		},
		{
			Name: "SupportURL",
			Prompt: &survey.Input{
				Message: "What is support URL?",
				Default: "https://example.com",
			},
		},
	}
	err := survey.Ask(qs, &metadata)
	if err != nil {
		return nil, errors.Wrap(err, "while asking for metadata")
	}

	maintainers, err := askForMaintainers()
	if err != nil {
		return nil, errors.Wrap(err, "while asking for maintainers")
	}
	metadata.Maintainers = maintainers
	return &metadata, nil
}

func askForMaintainers() ([]common.Maintainers, error) {
	var maintainers []common.Maintainers
	for {
		name := false
		prompt := &survey.Confirm{
			Message: "Do you want to add maintainer?",
		}
		err := survey.AskOne(prompt, &name)
		if err != nil {
			return nil, errors.Wrap(err, "while asking if add maintainers")
		}
		if !name {
			return maintainers, nil
		}

		maintainer, err := askForMaintainer()
		if err != nil {
			return nil, errors.Wrap(err, "while asking if for maintainer")
		}
		maintainers = append(maintainers, maintainer)
	}
}

func askForMaintainer() (common.Maintainers, error) {
	var maintainer common.Maintainers
	var qs = []*survey.Question{
		{
			Name: "Email",
			Prompt: &survey.Input{
				Message: "What is email",
				Default: "dev@example.com",
			},
		},
		{
			Name: "Name",
			Prompt: &survey.Input{
				Message: "What is a name?",
				Default: "Example Dev",
			},
		},
		{
			Name: "Url",
			Prompt: &survey.Input{
				Message: "What is a Url?",
				Default: "https://example.com",
			},
		},
	}
	err := survey.Ask(qs, &maintainer)
	return maintainer, err
}

func askForManifestPathSuffix() (string, error) {
	var manifestPath string
	prompt := []*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Manifest path suffix",
			},
			Validate: func(ans interface{}) error {
				if str, ok := ans.(string); !ok || len(strings.Split(str, ".")) < 2 {
					return errors.New(`manifest path suffix must be in format "[PREFIX].[NAME]"`)

				}
				return nil
			},
		},
	}
	err := survey.Ask(prompt, &manifestPath)
	return manifestPath, err
}

func askIfOverwrite() (bool, error) {
	overwrite := false
	prompt := &survey.Confirm{
		Message: "Do you want to overwrite existing manifest files?",
	}
	err := survey.AskOne(prompt, &overwrite)
	return overwrite, err
}
