package implementation

import (
	"capact.io/capact/cmd/cli/cmd/alpha/manifest-gen/common"
	"capact.io/capact/internal/cli/alpha/manifestgen"
	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/utils/strings/slices"
)

var (
	helmTool      = "Helm"
	terraformTool = "Terraform"
	emptyManifest = "Empty"
)

type generateFun func(opts common.ManifestGenOptions) (map[string]string, error)

type helmChart struct {
	Url     string
	Version string
}

// NewCmd returns a cobra.Command for Implementation manifest generation operations.
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "implementation",
		Aliases: []string{"impl", "implementations"},
		Short:   "Generate new Implementation manifests",
		Long:    "Generate new Implementation manifests for various tools.",
	}

	cmd.AddCommand(NewTerraform())
	cmd.AddCommand(NewHelm())

	return cmd
}

// HandleInteractiveSession is responsible for handling interactive session with user
func HandleInteractiveSession(opts common.ManifestGenOptions) (map[string]string, error) {
	tool, err := askForImplementationTool()
	if err != nil {
		return nil, errors.Wrap(err, "while asking for used implementation tool")
	}

	toolAction := map[string]generateFun{
		helmTool:      generateHelmManifests,
		terraformTool: generateTerraformManifests,
		emptyManifest: generateEmptyManifests,
	}

	return toolAction[tool](opts)
}

func generateHelmManifests(opts common.ManifestGenOptions) (map[string]string, error) {
	basedToolDir, err := common.AskForDirectory("Path to helm template", "")
	if err != nil {
		return nil, errors.Wrap(err, "while asking for path to helm template")
	}
	helmchartInfo, err := askForHelmChartDetails()
	if err != nil {
		return nil, errors.Wrap(err, "while asking for path to helm template")
	}

	var helmCfg manifestgen.HelmConfig
	helmCfg.ManifestPath = common.CreateManifestPath(common.ImplementationManifest, opts.ManifestPath)
	helmCfg.ChartName = basedToolDir
	helmCfg.ManifestMetadata = opts.Metadata
	helmCfg.ChartRepoURL = helmchartInfo.Url
	helmCfg.ChartVersion = helmchartInfo.Version
	if slices.Contains(opts.ManifestsType, common.InterfaceManifest) {
		helmCfg.InterfacePathWithRevision = "cap.interface." + opts.ManifestPath + ":0.1.0"
	}
	files, err := manifestgen.GenerateHelmManifests(&helmCfg)
	if err != nil {
		return nil, errors.Wrap(err, "while generating Helm manifests")
	}
	return files, nil
}

func generateTerraformManifests(opts common.ManifestGenOptions) (map[string]string, error) {
	basedToolDir, err := common.AskForDirectory("Path to terraform template", "")
	if err != nil {
		return nil, errors.Wrap(err, "while asking for path to terraform template")
	}

	provider, err := askForProvider()
	if err != nil {
		return nil, errors.Wrap(err, "while asking for path to terraform template")
	}

	source, err := askForSource()
	if err != nil {
		return nil, errors.Wrap(err, "while asking for path to terraform template")
	}

	var tfContentCfg manifestgen.TerraformConfig
	tfContentCfg.ManifestPath = common.CreateManifestPath(common.ImplementationManifest, opts.ManifestPath)
	tfContentCfg.ModulePath = basedToolDir
	tfContentCfg.ManifestMetadata = opts.Metadata
	tfContentCfg.Provider = provider
	tfContentCfg.ModuleSourceURL = source
	if slices.Contains(opts.ManifestsType, common.InterfaceManifest) {
		tfContentCfg.InterfacePathWithRevision = "cap.interface." + opts.ManifestPath + ":0.1.0"
	}
	files, err := manifestgen.GenerateTerraformManifests(&tfContentCfg)
	if err != nil {
		return nil, errors.Wrap(err, "while generating Terraform manifests")
	}
	return files, nil
}

func generateEmptyManifests(opts common.ManifestGenOptions) (map[string]string, error) {
	var emptyManifestCfg manifestgen.EmptyImplementationConfig
	emptyManifestCfg.ManifestPath = common.CreateManifestPath(common.ImplementationManifest, opts.ManifestPath)
	emptyManifestCfg.ManifestMetadata = opts.Metadata
	if slices.Contains(opts.ManifestsType, common.InterfaceManifest) {
		emptyManifestCfg.InterfacePathWithRevision = "cap.interface." + opts.ManifestPath + ":0.1.0"
	}
	files, err := manifestgen.GenerateEmptyManifests(&emptyManifestCfg)
	if err != nil {
		return nil, errors.Wrap(err, "while generating Empty manifests")
	}
	return files, nil
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
				Default: "https://example.com/terraform-module.tgz",
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
			Name: "Url",
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
