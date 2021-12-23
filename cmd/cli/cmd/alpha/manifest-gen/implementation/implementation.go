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
	var helmCfg manifestgen.HelmConfig
	helmCfg.ManifestPath = "cap.implementation." + opts.ManifestPath
	helmCfg.ChartName = basedToolDir
	helmCfg.ManifestMetadata = opts.Metadata
	if slices.Contains(opts.ManifestsType, common.InterfaceType) {
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
	var tfContentCfg manifestgen.TerraformConfig
	tfContentCfg.ManifestPath = "cap.implementation." + opts.ManifestPath
	tfContentCfg.ModulePath = basedToolDir
	tfContentCfg.ManifestMetadata = opts.Metadata
	if slices.Contains(opts.ManifestsType, common.InterfaceType) {
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
	emptyManifestCfg.ManifestPath = "cap.implementation." + opts.ManifestPath
	emptyManifestCfg.ManifestMetadata = opts.Metadata
	if slices.Contains(opts.ManifestsType, common.InterfaceType) {
		emptyManifestCfg.InterfacePathWithRevision = "cap.interface." + opts.ManifestPath + ":0.1.0"
	}
	files, err := manifestgen.GenerateEmptyManifests(&emptyManifestCfg)
	if err != nil {
		return nil, errors.Wrap(err, "while generating Empty manifests")
	}
	return files, nil
}

func askForImplementationTool() (string, error) {
	var selectTool string
	availableManifestsTool := []string{helmTool, terraformTool, emptyManifest}
	toolPrompt := &survey.Select{
		Message: "Based on which tool do you want to generate implementation:",
		Options: availableManifestsTool,
	}
	err := survey.AskOne(toolPrompt, &selectTool)
	return selectTool, err
}
