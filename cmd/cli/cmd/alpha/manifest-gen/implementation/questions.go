package implementation

import (
	"capact.io/capact/internal/cli/alpha/manifestgen"
	"github.com/AlecAivazis/survey/v2"
)

type helmChart struct {
	// URL is address to helm repository
	URL string
	// Version defines a helm chart version
	Version string
}

func askForImplementationTool() (string, error) {
	var selectedTool string
	availableTool := []string{helmTool, terraformTool, emptyManifest}
	prompt := &survey.Select{
		Message: "Based on which tool do you want to generate implementation:",
		Options: availableTool,
	}
	err := survey.AskOne(prompt, &selectedTool)
	return selectedTool, err
}

func askForProvider() (manifestgen.Provider, error) {
	var selectedProvider string
	availableProviders := []string{string(manifestgen.ProviderAWS), string(manifestgen.ProviderGCP)}
	prompt := &survey.Select{
		Message: "Create a provider-specific workflow:",
		Options: availableProviders,
	}
	err := survey.AskOne(prompt, &selectedProvider)
	return manifestgen.Provider(selectedProvider), err
}

func askForSource() (string, error) {
	var source string
	prompt := []*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Path to the Terraform module, such as URL to Tarball or Git repository",
				Default: "",
			},
		},
	}
	err := survey.Ask(prompt, &source)
	return source, err
}

func askForHelmChartDetails() (helmChart, error) {
	var helmChartInfo helmChart
	var qs = []*survey.Question{
		{
			Name: "URL",
			Prompt: &survey.Input{
				Message: "URL of the Helm repository",
				Default: "",
			},
		},
		{
			Name: "Version",
			Prompt: &survey.Input{
				Message: "Version of the Helm chart",
				Default: "",
			},
		},
	}
	err := survey.Ask(qs, &helmChartInfo)
	return helmChartInfo, err
}
