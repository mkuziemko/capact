package manifestgen

import (
	"capact.io/capact/cmd/cli/cmd/alpha/manifest-gen/common"
	"capact.io/capact/cmd/cli/cmd/alpha/manifest-gen/implementation"
	_interface "capact.io/capact/cmd/cli/cmd/alpha/manifest-gen/interface"
	"capact.io/capact/internal/cli/alpha/manifestgen"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/utils/strings/slices"
)

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

	typeFuncMap := map[string]getManifestFun{
		common.InterfaceType:      generateInterface,
		common.InterfaceGroupType: generateGroupInterface,
		common.TypeType:           generateType,
		common.AttributeType:      generateAttribute,
		common.ImplementationType: generateImplementation,
	}
	var mergeFiles map[string]string

	for k, v := range typeFuncMap {
		if slices.Contains(opts.ManifestsType, k) {
			files, err := v(opts)
			if err != nil {
				return errors.Wrap(err, "when generating the file")
			}
			mergeFiles = MergeMaps(mergeFiles, files)
		}
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
