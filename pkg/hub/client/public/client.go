package public

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	gqlpublicapi "capact.io/capact/pkg/hub/api/graphql/public"
	"github.com/avast/retry-go"

	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
)

const retryAttempts = 1

// Client used to communicate with the Capact Public Hub GraphQL APIs
type Client struct {
	client *graphql.Client
}

// NewClient creates a public client with a given GraphQL custom client instance.
func NewClient(cli *graphql.Client) *Client {
	return &Client{client: cli}
}

// FindInterfaceRevision returns the InterfaceRevision for the given InterfaceReference.
// It will return nil, if the InterfaceRevision is not found.
func (c *Client) FindInterfaceRevision(ctx context.Context, ref gqlpublicapi.InterfaceReference, opts ...InterfaceRevisionOption) (*gqlpublicapi.InterfaceRevision, error) {
	findOpts := &InterfaceRevisionOptions{}
	findOpts.Apply(opts...)

	query, params := c.interfaceQueryForRef(findOpts.fields, ref)
	req := graphql.NewRequest(fmt.Sprintf(`query FindInterfaceRevision($interfacePath: NodePath!, %s) {
		  interface(path: $interfacePath) {
				%s
		  }
		}`, params.Query(), query))

	req.Var("interfacePath", ref.Path)
	params.PopulateVars(req)

	var resp struct {
		Interface struct {
			Revision *gqlpublicapi.InterfaceRevision `json:"rev"`
		} `json:"interface"`
	}
	err := retry.Do(func() error {
		return c.client.Run(ctx, req, &resp)
	}, retry.Attempts(retryAttempts))

	if err != nil {
		return nil, errors.Wrap(err, "while executing query to fetch Hub Interface Revision")
	}

	return resp.Interface.Revision, nil
}

// ListTypeRefRevisionsJSONSchemas returns the list of requested Types.
// Only a few fields are populated. Check the query fields for more information.
func (c *Client) ListTypeRefRevisionsJSONSchemas(ctx context.Context, filter gqlpublicapi.TypeFilter) ([]*gqlpublicapi.TypeRevision, error) {
	req := graphql.NewRequest(`query ListTypeRefsJSONSchemas($typeFilter: TypeFilter!)  {
		  types(filter: $typeFilter) {
			revisions {
			  revision
			  metadata {
				path
			  }
			  spec {
				jsonSchema
			  }
			}
		  }
		}`)

	req.Var("typeFilter", filter)

	var resp struct {
		Types []*gqlpublicapi.Type `json:"types"`
	}
	err := retry.Do(func() error {
		return c.client.Run(ctx, req, &resp)
	}, retry.Attempts(retryAttempts))

	if err != nil {
		return nil, errors.Wrap(err, "while executing query to list Types")
	}

	var out []*gqlpublicapi.TypeRevision
	for _, t := range resp.Types {
		out = append(out, t.Revisions...)
	}

	return out, nil
}

// ListInterfaces returns all Interfaces. By default only root fields are populated. Use options to add
// latestRevision fields or apply additional filtering.
func (c *Client) ListInterfaces(ctx context.Context, opts ...InterfaceOption) ([]*gqlpublicapi.Interface, error) {
	ifaceOpts := &InterfaceOptions{}
	ifaceOpts.Apply(opts...)

	req := graphql.NewRequest(fmt.Sprintf(`query ListInterfaces($interfaceFilter: InterfaceFilter!)  {
		  interfaces(filter: $interfaceFilter) {
			path
			name
			prefix
			%s
		  }
		}`, ifaceOpts.additionalFields))

	req.Var("interfaceFilter", ifaceOpts.filter)

	var resp struct {
		Interfaces []*gqlpublicapi.Interface `json:"interfaces"`
	}
	err := retry.Do(func() error {
		return c.client.Run(ctx, req, &resp)
	}, retry.Attempts(retryAttempts))

	if err != nil {
		return nil, errors.Wrap(err, "while executing query to list Hub Interfaces")
	}

	return resp.Interfaces, nil
}

// GetInterfaceLatestRevisionString returns the latest revision of the available Interfaces.
// Semantic versioning is used to determine the latest revision.
func (c *Client) GetInterfaceLatestRevisionString(ctx context.Context, ref gqlpublicapi.InterfaceReference) (string, error) {
	req := graphql.NewRequest(`query GetInterfaceLatestRevisionString($interfacePath: NodePath!) {
		interface(path: $interfacePath) {
			latestRevision {
				revision
			}
		}		
	}`)

	req.Var("interfacePath", ref.Path)

	var resp struct {
		Interface struct {
			LatestRevision *struct {
				Revision string `json:"revision"`
			} `json:"latestRevision"`
		} `json:"interface"`
	}
	err := retry.Do(func() error {
		return c.client.Run(ctx, req, &resp)
	}, retry.Attempts(retryAttempts))
	if err != nil {
		return "", errors.Wrap(err, "while executing query to fetch Interface latest revision string")
	}

	if resp.Interface.LatestRevision == nil {
		return "", fmt.Errorf("cannot find latest revision for Interface %q", ref.Path)
	}

	return resp.Interface.LatestRevision.Revision, nil
}

// ListImplementationRevisions returns ImplementationRevisions. Use options to apply additional filtering.
func (c *Client) ListImplementationRevisions(ctx context.Context, opts ...ListImplementationRevisionsOption) ([]*gqlpublicapi.ImplementationRevision, error) {
	getOpts := &ListImplementationRevisionsOptions{}
	getOpts.Apply(opts...)

	req := graphql.NewRequest(fmt.Sprintf(`query ListImplementationRevisions{
		implementations {
			revisions {
				%s
			}
		}
	}`, getOpts.fields))

	var resp struct {
		Implementations []gqlpublicapi.Implementation `json:"implementations"`
	}

	err := retry.Do(func() error {
		return c.client.Run(ctx, req, &resp)
	}, retry.Attempts(retryAttempts))

	if err != nil {
		return nil, errors.Wrap(err, "while executing query to fetch Hub Implementations")
	}

	var revs []*gqlpublicapi.ImplementationRevision

	for _, impl := range resp.Implementations {
		revs = append(revs, impl.Revisions...)
	}

	return revs, nil
}

// ListImplementationRevisionsForInterface returns ImplementationRevisions for the given Interface.
func (c *Client) ListImplementationRevisionsForInterface(ctx context.Context, ref gqlpublicapi.InterfaceReference, opts ...ListImplementationRevisionsForInterfaceOption) ([]gqlpublicapi.ImplementationRevision, error) {
	getOpts := &ListImplementationRevisionsForInterfaceOptions{}
	getOpts.Apply(opts...)

	query, params := c.interfaceQueryForRef(ifaceRevisionAllFields, ref)
	req := graphql.NewRequest(fmt.Sprintf(`query ListImplementationRevisionsForInterface($interfacePath: NodePath!, %s) {
		  interface(path: $interfacePath) {
				%s
		  }
		}`, params.Query(), query))

	req.Var("interfacePath", ref.Path)
	params.PopulateVars(req)

	var resp struct {
		Interface struct {
			LatestRevision struct {
				ImplementationRevisions []gqlpublicapi.ImplementationRevision `json:"implementationRevisions"`
			} `json:"rev"`
		} `json:"interface"`
	}
	err := retry.Do(func() error {
		return c.client.Run(ctx, req, &resp)
	}, retry.Attempts(retryAttempts))

	if err != nil {
		return nil, errors.Wrap(err, "while executing query to fetch Hub Implementation")
	}

	result := FilterImplementationRevisions(resp.Interface.LatestRevision.ImplementationRevisions, getOpts)

	result = SortImplementationRevisions(result, getOpts)

	return result, nil
}

// CheckManifestRevisionsExist checks if manifests with provided manifest references exist.
func (c *Client) CheckManifestRevisionsExist(ctx context.Context, manifestRefs []gqlpublicapi.ManifestReference) (map[gqlpublicapi.ManifestReference]bool, error) {
	if len(manifestRefs) == 0 {
		return map[gqlpublicapi.ManifestReference]bool{}, nil
	}

	getAlias := func(i int) string {
		return fmt.Sprintf("partial%d", i)
	}

	strBuilder := strings.Builder{}
	for i, manifestRef := range manifestRefs {
		alias := getAlias(i)
		queryName, err := manifestRef.GQLQueryName()
		if err != nil {
			return nil, errors.Wrap(err, "while getting GraphQL query name for a given manifest")
		}

		partialQuery := fmt.Sprintf(`
			%s: %s(path:"%s") {
				revision(revision:"%s") {
					revision
				}
			}
		`, alias, queryName, manifestRef.Path, manifestRef.Revision)
		strBuilder.WriteString(partialQuery)
	}

	req := graphql.NewRequest(fmt.Sprintf(`
		query CheckManifestRevisionsExist {
			%s
		}`,
		strBuilder.String(),
	))

	var resp map[string]struct {
		Revision struct {
			Revision *string `json:"revision"`
		} `json:"revision"`
	}

	err := retry.Do(func() error {
		return c.client.Run(ctx, req, &resp)
	}, retry.Attempts(retryAttempts))

	if err != nil {
		return nil, errors.Wrap(err, "while executing query to check Type Revisions exist")
	}

	result := map[gqlpublicapi.ManifestReference]bool{}
	for i, manifestRef := range manifestRefs {
		alias := getAlias(i)
		result[manifestRef] = resp[alias].Revision.Revision != nil
	}

	return result, nil
}

var key = regexp.MustCompile(`\$(\w+):`)

// Args is used to store arguments to GraphQL queries.
type Args map[string]interface{}

// Query returns the definition for the arguments
// stored in this Args, which has to be put in the
// GraphQL query.
func (a Args) Query() string {
	var out []string
	for k := range a {
		out = append(out, k)
	}
	return strings.Join(out, ",")
}

// PopulateVars fills the variables stores in this Args
// in the provided *graphql.Request.
func (a Args) PopulateVars(req *graphql.Request) {
	for k, v := range a {
		name := key.FindStringSubmatch(k)
		req.Var(name[1], v)
	}
}

func (c *Client) interfaceQueryForRef(fields string, ref gqlpublicapi.InterfaceReference) (string, Args) {
	if ref.Revision == "" {
		return c.latestInterfaceRevision(fields)
	}

	return c.specificInterfaceRevision(fields, ref.Revision)
}

func (c *Client) latestInterfaceRevision(fields string) (string, Args) {
	latestRevision := fmt.Sprintf(`
			rev: latestRevision {
				%s
			}`, fields)

	return latestRevision, Args{}
}

func (c *Client) specificInterfaceRevision(fields string, rev string) (string, Args) {
	specificRevision := fmt.Sprintf(`
			rev: revision(revision: $interfaceRev) {
				%s
			}`, fields)

	return specificRevision, Args{
		"$interfaceRev: Version!": rev,
	}
}
