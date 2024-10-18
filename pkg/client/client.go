package client

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	userv1 "github.com/openshift/client-go/user/clientset/versioned/typed/user/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	config      *rest.Config
	usersClient *userv1.UserV1Client
	k8sClient   *kubernetes.Clientset
}

func New(c *rest.Config) (*Client, error) {
	usrc, err := userv1.NewForConfig(c)
	if err != nil {
		return nil, fmt.Errorf("unable to create UserV1Client client, error: %w", err)
	}
	k8sc, err := kubernetes.NewForConfig(c)
	if err != nil {
		return nil, fmt.Errorf("unable to create k8s client, error: %w", err)
	}

	return &Client{config: c, usersClient: usrc, k8sClient: k8sc}, nil
}

// ListUsers list the users of the Openshift cluster
func (c *Client) ListUsers(ctx context.Context) ([]*v2.Resource, error) {
	list, err := c.usersClient.Users().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	users, err := convertV1Users2Resources(list.Items)
	if err != nil {
		return nil, fmt.Errorf("unable to convert []v1.User to []*v2.Resource, error: %w", err)
	}

	return users, nil
}

// ListRoles list the available (roles) entitlements
func (c *Client) ListRoles(ctx context.Context, namespace string) ([]*v2.Resource, error) {
	list, err := c.k8sClient.RbacV1().Roles(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to list entitlements, error: %w", err)
	}

	roles, err := convertV1RoleLists2Resources(list.Items)
	if err != nil {
		return nil, fmt.Errorf("unable to convert []v1.Role to []*v2.Resource, error: %w", err)
	}

	return roles, nil
}

func (c *Client) ListRoleBindings(ctx context.Context, namespace string, entitlement *v2.Resource, users []*v2.Resource) ([]*v2.Grant, error) {
	list, err := c.k8sClient.RbacV1().RoleBindings(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to list grants, error: %w", err)
	}

	grants, err := convertV1RoleBindings2Resources(list.Items, entitlement, users)
	if err != nil {
		return nil, fmt.Errorf("unable to convert []v1.RoleBinding to []*v2.Grant, error: %w", err)
	}

	return grants, nil
}

func (c *Client) ListGroups(ctx context.Context) ([]*v2.Resource, error) {
	list, err := c.usersClient.Groups().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	groups, err := convertV1Groups2Resources(list.Items)
	if err != nil {
		return nil, fmt.Errorf("unable to convert []v1.Group to []*v2.Resource, error: %w", err)
	}

	return groups, nil
}

func (c *Client) MatchUsersToGroup(ctx context.Context, entitlement *v2.Resource, users []*v2.Resource) ([]*v2.Grant, error) {
	var gnts []*v2.Grant

	list, err := c.usersClient.Groups().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, group := range list.Items {
		// match a group with the entitlement
		if entitlement.Id.Resource == string(group.UID) {
			for _, user := range users {
				for _, member := range group.Users {
					// NOTE(shackra): 3 levels of for-loops isn't that performing! right?
					if user.DisplayName == member {
						gnts = append(gnts, grant.NewGrant(entitlement, "member", user.Id))
					}
				}
			}
		}
	}

	return gnts, nil
}
