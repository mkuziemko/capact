package common

import (
	"capact.io/capact/internal/cli/alpha/manifestgen"
)

type Metadata = manifestgen.MetaDataInfo
type Maintainers = manifestgen.Maintainer

// ManifestGenOptions is a struct based on which manifests are generated
type ManifestGenOptions struct {
	ManifestsType  []string
	ManifestPath   string
	Directory      string
	Overwrite      bool
	Metadata       Metadata
	InterfacePath  string
	TypeInputPath  string
	TypeOutputPath string
	Revision       string
}

var (
	InterfaceManifest      = "Interface"
	InterfaceGroupManifest = "InterfaceGroup"
	ImplementationManifest = "Implementation"
	TypeManifest           = "Type"
	AttributeManifest      = "Attribute"
	GCPProvider            = "GCP"
	AWSProvider            = "AWS"
)
