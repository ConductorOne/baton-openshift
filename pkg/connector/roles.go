package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type roleBuilder struct {
	namespace string
	roles     *dynamic.DynamicClient
}

func (o *roleBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return roleResourceType
}

// List returns all the roles from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *roleBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	roles, err := o.roles.Resource(roleOpenshiftGVR).Namespace(o.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, "", nil, err
	}
	list, err := convertV1Roles2Resources(roles.Items)
	if err != nil {
		return nil, "", nil, err
	}
	return list, "", nil, nil
}

// TODO(shackra): figure what should I do here for roles
func (o *roleBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// TODO(shackra): figure what should I do here for roles
func (o *roleBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newRoleBuilder(namespace string, config *rest.Config) (*roleBuilder, error) {
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &roleBuilder{
		namespace: namespace,
		roles:     client,
	}, nil
}

func convertV1Roles2Resources(roles []unstructured.Unstructured) ([]*v2.Resource, error) {
	var rsc []*v2.Resource
	for _, role := range roles {
		result, err := convertV1Role2Resource(role)
		if err != nil {
			return nil, fmt.Errorf("resource %s, error: %w", role.GetUID(), err)
		}
		rsc = append(rsc, result)
	}

	return rsc, nil
}

func convertV1Role2Resource(role unstructured.Unstructured) (*v2.Resource, error) {
	annos := annotations.Annotations{}
	// NOTE(shackra): Maybe this is not wanted for this case?
	annos.Update(&v2.SkipEntitlementsAndGrants{})

	profile := map[string]any{
		"name":          role.GetName(),
		"generate_name": role.GetGenerateName(),
		"annotations":   role.GetAnnotations(),
	}

	traits := []rs.RoleTraitOption{
		rs.WithRoleProfile(profile),
	}

	return rs.NewRoleResource(
		role.GetName(),
		&v2.ResourceType{
			Id:          "role",
			DisplayName: "Role",
			Traits: []v2.ResourceType_Trait{
				v2.ResourceType_TRAIT_ROLE,
			},
			Annotations: annos,
		},
		role.GetUID(),
		traits,
	)
}
