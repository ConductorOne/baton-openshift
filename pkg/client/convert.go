package client

import (
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
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
