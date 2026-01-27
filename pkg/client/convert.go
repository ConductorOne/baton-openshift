package client

// convert.go have helper functions that translate types from
// Openshift/Kubernetes to those that Baton SDK can comprehend. At
// some places DRY principles were not followed for the sake of
// clarity and code locality.
//
// All these helpers are used on client.go.

import (
	"errors"
	"fmt"
	"strings"
	"time"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	grant "github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	v1 "github.com/openshift/api/user/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

// convertV1Users2Resources (plural) convert users of Openshift to resources of Baton SDK.
func convertV1Users2Resources(users []v1.User) ([]*v2.Resource, error) {
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

// convertV1User2Resource (singular) convert a user to a resource, use by `convertV1Users2Resources`.
func convertV1User2Resource(user v1.User) (*v2.Resource, error) {
	annos := annotations.Annotations{}
	annos.Update(&v2.SkipEntitlementsAndGrants{})

	profile := map[string]interface{}{
		"name":          user.Name,
		"generate_name": user.GenerateName,
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
		string(user.UID),
		traits,
	)
}

// convertV1RoleLists2Resources (plural) convert a list of roles of Openshift to resources of Baton SDK.
func convertV1RoleLists2Resources(roleLists []rbacv1.Role) ([]*v2.Resource, error) {
	var rsc []*v2.Resource
	for _, role := range roleLists {
		result, err := convertV1RoleList2Resource(role)
		if err != nil {
			return nil, fmt.Errorf("resource %s, error: %w", role.UID, err)
		}
		rsc = append(rsc, result)
	}

	return rsc, nil
}

// convertV1RoleList2Resource (singular) convert a role to a resource, use by `convertV1RoleLists2Resources`.
func convertV1RoleList2Resource(roleList rbacv1.Role) (*v2.Resource, error) {
	annos := annotations.Annotations{}
	annos.Update(&v2.SkipEntitlementsAndGrants{})

	profile := map[string]interface{}{
		"name":          roleList.Name,
		"generate_name": roleList.GenerateName,
	}

	traits := []rs.RoleTraitOption{
		rs.WithRoleProfile(profile),
		// FIXME(shackra): add creation time
	}

	return rs.NewRoleResource(
		roleList.Name,
		&v2.ResourceType{
			Id:          "role",
			DisplayName: "Role",
			Traits: []v2.ResourceType_Trait{
				v2.ResourceType_TRAIT_ROLE,
			},
			Annotations: annos,
		},
		string(roleList.UID),
		traits,
	)
}

var errRoleNotGranted = errors.New("role not granted to this resource")

// convertV1RoleBindings2Resources (plural) convert role bindings of Openshift to grants of Baton SDK.
// for a given entitlement and a  list of user resources.
func convertV1RoleBindings2Resources(roleBindings []rbacv1.RoleBinding, entitlement *v2.Resource, users []*v2.Resource) ([]*v2.Grant, error) {
	var grts []*v2.Grant
	for _, binding := range roleBindings {
		result, err := convertV1RoleBinding2Resource(binding, entitlement, users)
		if err != nil {
			if errors.Is(err, errRoleNotGranted) {
				// skip
				continue
			}
			return nil, fmt.Errorf("binding %s - resource %s, error: %w", binding.RoleRef.Name, entitlement.DisplayName, err)
		}
		grts = append(grts, result)
	}
	return grts, nil
}

// convertV1RoleBinding2Resource (singular) convert a role binding, for a given entitlement and a list of user resources.
// use by `convertV1RoleBindings2Resources`.
func convertV1RoleBinding2Resource(roleBinding rbacv1.RoleBinding, entitlement *v2.Resource, users []*v2.Resource) (*v2.Grant, error) {
	if len(roleBinding.Subjects) == 0 {
		return nil, errRoleNotGranted
	}

	splittedRoleBinding := strings.Split(roleBinding.Name, "-")
	if len(splittedRoleBinding) <= 1 {
		return nil, errRoleNotGranted
	}
	if strings.Contains(entitlement.DisplayName, splittedRoleBinding[1]) {
		for _, user := range users {
			if user.DisplayName == roleBinding.Subjects[0].Name {
				return grant.NewGrant(
					entitlement,
					"member",
					user.Id,
				), nil
			}
		}
	}

	return nil, errRoleNotGranted
}

// convertV1Groups2Resources (plural) convert a list of groups of Openshift to resources of Baton SDK.
func convertV1Groups2Resources(groups []v1.Group) ([]*v2.Resource, error) {
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

// convertV1Group2Resource (singular) convert a group into a resource, use by `convertV1Groups2Resources`.
func convertV1Group2Resource(group v1.Group) (*v2.Resource, error) {
	profile := map[string]any{
		"name":          group.GetName(),
		"generate_name": group.GetGenerateName(),
		"created_at":    group.CreationTimestamp.Format(time.RFC3339),
	}

	traits := []rs.GroupTraitOption{
		rs.WithGroupProfile(profile),
	}

	return rs.NewGroupResource(
		group.GetName(),
		&v2.ResourceType{
			Id:          "group",
			DisplayName: "Team",
			Traits: []v2.ResourceType_Trait{
				v2.ResourceType_TRAIT_GROUP,
			},
		},
		string(group.UID),
		traits,
	)
}
