package client

import (
	"context"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	userv1 "github.com/openshift/client-go/user/clientset/versioned/typed/user/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type Client struct {
	usersClient *userv1.UserV1Client
}

func New(c *rest.Config) (*Client, error) {
	usrc, err := userv1.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	return &Client{usersClient: usrc}, nil
}

// ListUsers list the users of the Openshift cluster
func (c *Client) ListUsers(ctx context.Context) ([]*v2.Resource, error) {
	list, err := c.usersClient.Users().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	users, err := convertV1Users2Resources(list.Items)
	if err != nil {
		return nil, err
	}

	return users, nil
}
