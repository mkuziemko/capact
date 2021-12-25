package attribute

import (
	"capact.io/capact/cmd/cli/cmd/alpha/manifest-gen/common"
	"capact.io/capact/internal/cli/alpha/manifestgen"
	"github.com/pkg/errors"
)

func GenerateAttributeFile(opts common.ManifestGenOptions) (map[string]string, error) {
	var attributeCfg manifestgen.AttributeConfig
	attributeCfg.ManifestPath = common.CreateManifestPath(common.AttributeManifest, opts.ManifestPath)
	attributeCfg.ManifestMetadata = opts.Metadata
	attributeCfg.ManifestRevision = opts.Revision
	files, err := manifestgen.GenerateAttributeTemplatingConfig(&attributeCfg)
	if err != nil {
		return nil, errors.Wrap(err, "while generating content files")
	}
	return files, nil
}
