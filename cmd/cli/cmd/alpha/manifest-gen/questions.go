package manifestgen

import (
	"strings"

	"capact.io/capact/cmd/cli/cmd/alpha/manifest-gen/common"
	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
)

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
