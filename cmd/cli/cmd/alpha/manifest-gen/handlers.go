package manifestgen

import (
	"capact.io/capact/cmd/cli/cmd/alpha/manifest-gen/attribute"
	"capact.io/capact/cmd/cli/cmd/alpha/manifest-gen/common"
	"capact.io/capact/cmd/cli/cmd/alpha/manifest-gen/implementation"
	_interface "capact.io/capact/cmd/cli/cmd/alpha/manifest-gen/interface"
	"capact.io/capact/internal/cli/alpha/manifestgen"
	"github.com/pkg/errors"
	"k8s.io/utils/strings/slices"
)

type getManifestFun func(opts common.ManifestGenOptions) (map[string]string, error)

func generateInterface(opts common.ManifestGenOptions) (map[string]string, error) {
	if slices.Contains(opts.ManifestsType, common.TypeType) {
		opts.TypeInputPath = common.CreateManifestPath(common.TypeType, opts.ManifestPath) + "-input:0.0.1"
		opts.TypeOutputPath = common.CreateManifestPath(common.TypeType, opts.ManifestPath) + "config:0.0.1"
	}
	files, err := _interface.GenerateInterfaceFile(opts, manifestgen.GenerateInterfaceTemplatingConfig)
	if err != nil {
		return nil, errors.Wrap(err, "while generating interface templating config")
	}
	return files, nil
}

func generateGroupInterface(opts common.ManifestGenOptions) (map[string]string, error) {
	files, err := _interface.GenerateInterfaceFile(opts, manifestgen.GenerateInterfaceGroupTemplatingConfig)
	if err != nil {
		return nil, errors.Wrap(err, "while generating interface group templating config")
	}
	return files, nil
}

func generateType(opts common.ManifestGenOptions) (map[string]string, error) {
	files, err := _interface.GenerateInterfaceFile(opts, manifestgen.GenerateTypeTemplatingConfig)
	if err != nil {
		return nil, errors.Wrap(err, "while generating type templating config")
	}
	return files, nil
}

func generateAttribute(opts common.ManifestGenOptions) (map[string]string, error) {
	files, err := attribute.GenerateAttributeFile(opts)
	if err != nil {
		return nil, errors.Wrap(err, "while generating type templating config")
	}
	return files, nil
}

func generateImplementation(opts common.ManifestGenOptions) (map[string]string, error) {
	files, err := implementation.HandleInteractiveSession(opts)
	if err != nil {
		return nil, errors.Wrap(err, "while generating implementation tool")
	}
	return files, nil
}
