package client

import (
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	v1 "github.com/openshift/api/user/v1"
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
		string(user.UID),
		traits,
	)
}
