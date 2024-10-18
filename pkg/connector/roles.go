package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-openshift/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
)

type roleBuilder struct {
	namespace string
	client    *client.Client
}

func (o *roleBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return roleResourceType
}

func (o *roleBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	rsc, err := o.client.ListRoles(ctx, o.namespace)
	if err != nil {
		return nil, "", nil, err
	}
	return rsc, "", nil, nil
}

func (o *roleBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement

	// NOTE(shackra): I don't really know what's needed and what
	// is superflous
	assigmentOptions := []ent.EntitlementOption{
		ent.WithGrantableTo(userResourceType),
		ent.WithDisplayName(fmt.Sprintf("%s Role member", resource.DisplayName)),
		ent.WithDescription(fmt.Sprintf("Access to %s role in %s namespace", resource.DisplayName, o.namespace)),
	}

	rv = append(rv, ent.NewAssignmentEntitlement(
		resource,
		"member",
		assigmentOptions...,
	))

	return rv, "", nil, nil
}

func (o *roleBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	users, err := o.client.ListUsers(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("unable to list users to match their permissions, error: %w", err)
	}
	// NOTE(shackra): resource is a role, not a user!
	grants, err := o.client.ListRoleBindings(ctx, o.namespace, resource, users)
	if err != nil {
		return nil, "", nil, err
	}
	return grants, "", nil, nil
}

func newRoleBuilder(namespace string, clt *client.Client) *roleBuilder {
	return &roleBuilder{
		namespace: namespace,
		client:    clt,
	}
}
