package connector

import (
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// The user resource type is for all user objects from the database.
var userResourceType = &v2.ResourceType{
	Id:          "user",
	DisplayName: "User",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_USER},
}

var groupResourceType = &v2.ResourceType{
	Id:          "group",
	DisplayName: "Group",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_GROUP},
}

var roleResourceType = &v2.ResourceType{
	Id:          "role",
	DisplayName: "Role",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_ROLE},
}

// Openshift related
var groupOpenshiftGVR = schema.GroupVersionResource{
	Group:    "user.openshift.io",
	Version:  "v1",
	Resource: "groups",
}

var roleOpenshiftGVR = schema.GroupVersionResource{
	Group:    "rbac.authorization.k8s.io",
	Version:  "v1",
	Resource: "roles",
}
