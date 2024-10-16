package client

import (
	"errors"
	"fmt"
	"strings"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	grant "github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	v1 "github.com/openshift/api/user/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

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

func convertV1User2Resource(user v1.User) (*v2.Resource, error) {
	annos := annotations.Annotations{}
	// NOTE(shackra): Maybe this is not wanted for this case?
	annos.Update(&v2.SkipEntitlementsAndGrants{})

	profile := map[string]interface{}{
		"name":          user.Name,
		"generate_name": user.GenerateName,
		//"annotations":   user.Annotations, // TODO(shackra): check the docs and parse this accordingly
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

var roleNotGranted = errors.New("role not granted to this resource")

func convertV1RoleBindings2Resources(roleBindings []rbacv1.RoleBinding, entitlement *v2.Resource, users []*v2.Resource) ([]*v2.Grant, error) {
	var grts []*v2.Grant
	for _, binding := range roleBindings {
		result, err := convertV1RoleBinding2Resource(binding, entitlement, users)
		if err != nil {
			if errors.Is(err, roleNotGranted) {
				// skip
				continue
			}
			return nil, fmt.Errorf("binding %s - resource %s, error: %w", binding.RoleRef.Name, entitlement.DisplayName, err)
		}
		grts = append(grts, result)
	}
	return grts, nil
}

func convertV1RoleBinding2Resource(roleBinding rbacv1.RoleBinding, entitlement *v2.Resource, users []*v2.Resource) (*v2.Grant, error) {
	if len(roleBinding.Subjects) == 0 {
		return nil, roleNotGranted
	}

	splittedRoleBinding := strings.Split(roleBinding.Name, "-")
	if len(splittedRoleBinding) <= 1 {
		return nil, roleNotGranted
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

	return nil, roleNotGranted
}
