package fake

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"

	hublocalgraphql "capact.io/capact/pkg/hub/api/graphql/local"
	hubpublicgraphql "capact.io/capact/pkg/hub/api/graphql/public"
	"capact.io/capact/pkg/hub/client/public"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

const manifestsExtension = ".yaml"

type FileSystemClient struct {
	loadTypeInstances bool
	TypeInstances     map[string]hublocalgraphql.TypeInstance
	Implementations   []hubpublicgraphql.ImplementationRevision
	Interfaces        []hubpublicgraphql.InterfaceRevision
}

func NewFromLocal(manifestDir string, loadTypeInstances bool) (*FileSystemClient, error) {
	cli := &FileSystemClient{
		loadTypeInstances: loadTypeInstances,
		Implementations:   []hubpublicgraphql.ImplementationRevision{},
		Interfaces:        []hubpublicgraphql.InterfaceRevision{},
		TypeInstances:     map[string]hublocalgraphql.TypeInstance{},
	}

	if err := cli.loadManifests(manifestDir); err != nil {
		return nil, errors.Wrap(err, "while loading Hub manifests")
	}

	return cli, nil
}

func (s *FileSystemClient) ListImplementationRevisionsForInterface(ctx context.Context, ref hubpublicgraphql.InterfaceReference, opts ...public.GetImplementationOption) ([]hubpublicgraphql.ImplementationRevision, error) {
	getOpts := &public.ListImplementationRevisionsOptions{}
	getOpts.Apply(opts...)

	var out []hubpublicgraphql.ImplementationRevision
	for i := range s.Implementations {
		impl := s.Implementations[i]
		for _, implements := range impl.Spec.Implements {
			if implements.Path == ref.Path {
				item := hubpublicgraphql.ImplementationRevision{}

				if err := deepCopy(&impl, &item); err != nil {
					return nil, err
				}

				out = append(out, item)
			}
		}
	}

	result := public.FilterImplementationRevisions(out, getOpts)

	return public.SortImplementationRevisions(result, getOpts), nil
}

func (s *FileSystemClient) ListTypeInstancesTypeRef(ctx context.Context) ([]hublocalgraphql.TypeInstanceTypeReference, error) {
	var typeInstanceTypeRefs []hublocalgraphql.TypeInstanceTypeReference
	for _, ti := range s.TypeInstances {
		if ti.TypeRef == nil {
			continue
		}

		typeInstanceTypeRefs = append(typeInstanceTypeRefs, *ti.TypeRef)
	}

	return typeInstanceTypeRefs, nil
}

func (s *FileSystemClient) GetInterfaceLatestRevisionString(ctx context.Context, ref hubpublicgraphql.InterfaceReference) (string, error) {
	var versions semver.Collection
	for _, impl := range s.Implementations {
		for _, implements := range impl.Spec.Implements {
			if implements.Path == ref.Path {
				v, err := semver.NewVersion(implements.Revision)
				if err != nil {
					return "", err
				}
				versions = append(versions, v)
			}
		}
	}

	if len(versions) == 0 {
		return "", errors.New("no Interface found for a given ref")
	}

	sort.Sort(versions)
	latestVersion := versions[len(versions)-1]
	return latestVersion.String(), nil
}

func (s *FileSystemClient) FindInterfaceRevision(ctx context.Context, ref hubpublicgraphql.InterfaceReference) (*hubpublicgraphql.InterfaceRevision, error) {
	for i := range s.Interfaces {
		iface := s.Interfaces[i]
		if iface.Metadata.Path != ref.Path {
			continue
		}

		item := hubpublicgraphql.InterfaceRevision{}

		if err := deepCopy(&iface, &item); err != nil {
			return nil, err
		}

		return &item, nil
	}

	return nil, nil
}

func (s *FileSystemClient) FindTypeInstance(_ context.Context, id string) (*hublocalgraphql.TypeInstance, error) {
	ti, found := s.TypeInstances[id]
	if !found {
		return nil, nil
	}

	return &ti, nil
}

func (s *FileSystemClient) loadManifests(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if ext := filepath.Ext(path); ext != manifestsExtension {
			return nil
		}

		if err := s.loadManifest(path); err != nil {
			return errors.Wrapf(err, "while loading manifest %s", path)
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "while walking through manifest dir")
	}

	return nil
}

func (s *FileSystemClient) loadManifest(filepath string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return errors.Wrap(err, "while reading file")
	}

	jsonData, err := yaml.YAMLToJSON(data)
	if err != nil {
		return errors.Wrap(err, "while converting YAML to JSON")
	}

	if strings.Contains(filepath, "implementation") {
		impl := hubpublicgraphql.ImplementationRevision{}
		if err := json.Unmarshal(jsonData, &impl); err != nil {
			return err
		}
		s.Implementations = append(s.Implementations, impl)
	}

	if strings.Contains(filepath, "interface") {
		iface := hubpublicgraphql.InterfaceRevision{}
		if err := json.Unmarshal(jsonData, &iface); err != nil {
			return err
		}
		s.Interfaces = append(s.Interfaces, iface)
	}

	if s.loadTypeInstances && strings.Contains(filepath, "typeinstance") {
		ti := hublocalgraphql.TypeInstance{}
		if err := json.Unmarshal(jsonData, &ti); err != nil {
			return err
		}
		s.TypeInstances[ti.ID] = ti
	}

	return nil
}

func deepCopy(src interface{}, dst interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &dst)
}