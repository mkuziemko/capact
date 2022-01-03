package common

import (
	"capact.io/capact/internal/cli/alpha/manifestgen"
)

// Metadata is a alias for MetaDataInfo struct
type Metadata = manifestgen.MetaDataInfo

// Maintainers is a alias for Maintainer struct
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
	// InterfaceManifest hold a name for Interface Manifest
	InterfaceManifest = "Interface"
	// InterfaceGroupManifest hold a name for InterfaceGroup Manifest
	InterfaceGroupManifest = "InterfaceGroup"
	// ImplementationManifest hold a name for Implementation Manifest
	ImplementationManifest = "Implementation"
	// TypeManifest hold a name for Type Manifest
	TypeManifest = "Type"
	// AttributeManifest hold a name for Attribute Manifest
	AttributeManifest = "Attribute"
	// GCPProvider hold a name for GCP Provider
	GCPProvider = "GCP"
	// AWSProvider hold a name for AWS Provider
	AWSProvider = "AWS"
)
