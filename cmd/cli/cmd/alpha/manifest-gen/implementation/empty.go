package implementation

import (
	"strings"

	"capact.io/capact/cmd/cli/cmd/alpha/manifest-gen/common"
	"capact.io/capact/internal/cli/alpha/manifestgen"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/utils/strings/slices"
)

// NewEmpty returns a cobra.Command to bootstrap empty manifests.
func NewEmpty() *cobra.Command {
	var emptyCfg manifestgen.EmptyImplementationConfig

	cmd := &cobra.Command{
		Use:   "empty [MANIFEST_PATH] [HELM_CHART_NAME]",
		Short: "Generate empty implementation manifests",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("accepts two arguments: [MANIFEST_PATH]")
			}

			path := args[0]
			if !strings.HasPrefix(path, "cap.implementation.") || len(strings.Split(path, ".")) < 4 {
				return errors.New(`manifest path must be in format "cap.implementation.[PREFIX].[NAME]"`)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			emptyCfg.ManifestPath = args[0]
			emptyCfg.ManifestMetadata = common.GetDefaultMetadata()

			files, err := manifestgen.GenerateEmptyManifests(&emptyCfg)
			if err != nil {
				return errors.Wrap(err, "while generating Helm manifests")
			}

			outputDir, err := cmd.Flags().GetString("output")
			if err != nil {
				return errors.Wrap(err, "while reading output flag")
			}

			overrideManifests, err := cmd.Flags().GetBool("overwrite")
			if err != nil {
				return errors.Wrap(err, "while reading overwrite flag")
			}

			if err := manifestgen.WriteManifestFiles(outputDir, files, overrideManifests); err != nil {
				return errors.Wrap(err, "while writing manifest files")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&emptyCfg.InterfacePathWithRevision, "interface", "i", "", "Path with revision of the Interface, which is implemented by this Implementation")
	cmd.Flags().StringVarP(&emptyCfg.ManifestRevision, "revision", "r", "0.1.0", "Revision of the Implementation manifest")

	return cmd
}

func generateEmptyManifests(opts common.ManifestGenOptions) (map[string]string, error) {
	var emptyManifestCfg manifestgen.EmptyImplementationConfig
	emptyManifestCfg.ManifestPath = common.CreateManifestPath(common.ImplementationManifest, opts.ManifestPath)
	emptyManifestCfg.ManifestMetadata = opts.Metadata
	emptyManifestCfg.ManifestRevision = opts.Revision
	if slices.Contains(opts.ManifestsType, common.InterfaceManifest) {
		emptyManifestCfg.InterfacePathWithRevision = "cap.interface." + opts.ManifestPath + ":" + opts.Revision
	}
	files, err := manifestgen.GenerateEmptyManifests(&emptyManifestCfg)
	if err != nil {
		return nil, errors.Wrap(err, "while generating Empty manifests")
	}
	return files, nil
}
