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

type groupBuilder struct {
	namespace string
	groups    *dynamic.DynamicClient
}

func (o *groupBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return groupResourceType
}

// List returns all the groups from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *groupBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	groups, err := o.groups.Resource(groupOpenshiftGVR).Namespace(o.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, "", nil, err
	}
	list, err := convertV1Groups2Resources(groups.Items)
	if err != nil {
		return nil, "", nil, err
	}
	return list, "", nil, nil
}

// TODO(shackra): figure what should I do here for groups
func (o *groupBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// TODO(shackra): figure what should I do here for groups
func (o *groupBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newGroupBuilder(namespace string, config *rest.Config) (*groupBuilder, error) {
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &groupBuilder{
		namespace: namespace,
		groups:    client,
	}, nil
}

func convertV1Groups2Resources(groups []unstructured.Unstructured) ([]*v2.Resource, error) {
	var rsc []*v2.Resource
	for _, group := range groups {
		result, err := convertV1Group2Resource(group)
		if err != nil {
			return nil, fmt.Errorf("resource %s, error: %w", group.GetUID(), err)
		}
		rsc = append(rsc, result)
	}

	return rsc, nil
}

func convertV1Group2Resource(group unstructured.Unstructured) (*v2.Resource, error) {
	annos := annotations.Annotations{}
	// NOTE(shackra): Maybe this is not wanted for this case?
	annos.Update(&v2.SkipEntitlementsAndGrants{})

	profile := map[string]any{
		"name":          group.GetName(),
		"generate_name": group.GetGenerateName(),
		"annotations":   group.GetAnnotations(),
	}

	traits := []rs.GroupTraitOption{
		rs.WithGroupProfile(profile),
	}

	return rs.NewGroupResource(
		group.GetName(),
		&v2.ResourceType{
			Id:          "group",
			DisplayName: "Group",
			Traits: []v2.ResourceType_Trait{
				v2.ResourceType_TRAIT_GROUP,
			},
			Annotations: annos,
		},
		group.GetUID(),
		traits,
	)
}
