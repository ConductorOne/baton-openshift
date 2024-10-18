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

type groupBuilder struct {
	namespace string
	client    *client.Client
}

func (o *groupBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return groupResourceType
}

// List returns all the groups from the database as resource objects.
func (o *groupBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	groups, err := o.client.ListGroups(ctx)
	if err != nil {
		return nil, "", nil, err
	}
	return groups, "", nil, nil
}

func (o *groupBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement

	assigmentOptions := []ent.EntitlementOption{
		ent.WithGrantableTo(userResourceType),
		ent.WithDisplayName(fmt.Sprintf("%s Group member", resource.DisplayName)),
		ent.WithDescription(fmt.Sprintf("Access to %s group", resource.DisplayName)),
	}

	rv = append(rv, ent.NewAssignmentEntitlement(
		resource,
		"member",
		assigmentOptions...,
	))

	return rv, "", nil, nil
}

func (o *groupBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	users, err := o.client.ListUsers(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("unable to list users to match their group membership, error: %w", err)
	}

	grants, err := o.client.MatchUsersToGroup(ctx, resource, users)
	if err != nil {
		return nil, "", nil, fmt.Errorf("unable to match users membership to groups, error: %w", err)
	}
	return grants, "", nil, nil
}

func newGroupBuilder(namespace string, clt *client.Client) *groupBuilder {
	return &groupBuilder{
		namespace: namespace,
		client:    clt,
	}
}
