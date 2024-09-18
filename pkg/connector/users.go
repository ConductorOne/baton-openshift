package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	v1 "github.com/openshift/api/user/v1"
	userv1 "github.com/openshift/client-go/user/listers/user/v1"
)

type userBuilder struct {
	namespace string
	users     userv1.UserLister
}

func (o *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	users, err := o.users.List(nil)
	if err != nil {
		return nil, "", nil, err
	}
	list, err := convertV1Users2Resources(users)
	if err != nil {
		return nil, "", nil, err
	}
	return list, "", nil, nil
}

// Entitlements always returns an empty slice for users.
func (o *userBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *userBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newUserBuilder(namespace string) *userBuilder {
	// FIXME(shackra): provide cache.Indexer
	lister := userv1.NewUserLister(nil)
	return &userBuilder{
		namespace: namespace,
		users:     lister,
	}
}

func convertV1Users2Resources(users []*v1.User) ([]*v2.Resource, error) {
	var rsc []*v2.Resource
	for _, user := range users {
		result, err := convertV1User2Resource(user)
		if err != nil {
			return nil, fmt.Errorf("resource %s, error: %w", user.UID, err)
		}
		rsc = append(rsc, result)
	}

	return rsc, nil
}

func convertV1User2Resource(user *v1.User) (*v2.Resource, error) {
	annos := annotations.Annotations{}
	// NOTE(shackra): Maybe this is not wanted for this case?
	annos.Update(&v2.SkipEntitlementsAndGrants{})

	profile := map[string]any{
		"name":          user.Name,
		"generate_name": user.GenerateName,
		"annotations":   user.Annotations,
	}

	traits := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithCreatedAt(user.CreationTimestamp.Time),
	}

	return rs.NewUserResource(
		user.Name,
		&v2.ResourceType{
			Id:          "user",
			DisplayName: "User",
			Traits: []v2.ResourceType_Trait{
				v2.ResourceType_TRAIT_USER,
			},
			Annotations: annos,
		},
		user.UID,
		traits,
	)
}
