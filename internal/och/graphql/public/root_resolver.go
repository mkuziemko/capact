package public

import (
	"projectvoltron.dev/voltron/internal/och/graphql/public/resolver/implementations"
	interfacegroups "projectvoltron.dev/voltron/internal/och/graphql/public/resolver/interface-groups"
	"projectvoltron.dev/voltron/internal/och/graphql/public/resolver/interfaces"
	repometadata "projectvoltron.dev/voltron/internal/och/graphql/public/resolver/repo-metadata"
	"projectvoltron.dev/voltron/internal/och/graphql/public/resolver/tags"
	"projectvoltron.dev/voltron/internal/och/graphql/public/resolver/types"
	gqlpublicapi "projectvoltron.dev/voltron/pkg/och/api/graphql/public"
)

type RootResolver struct {
	queryResolver queryResolver
}

func NewRootResolver() *RootResolver {
	return &RootResolver{
		queryResolver: queryResolver{
			ImplementationResolver: implementations.NewResolver(),
			InterfaceResolver:      interfaces.NewResolver(),
			InterfaceGroupResolver: interfacegroups.NewResolver(),
			RepoMetadataResolver:   repometadata.NewResolver(),
			TagResolver:            tags.NewResolver(),
			TypeResolver:           types.NewResolver(),
		},
	}
}

func (r *RootResolver) Query() gqlpublicapi.QueryResolver {
	return r.queryResolver
}

type queryResolver struct {
	*implementations.ImplementationResolver
	*interfaces.InterfaceResolver
	*interfacegroups.InterfaceGroupResolver
	*repometadata.RepoMetadataResolver
	*tags.TagResolver
	*types.TypeResolver
}