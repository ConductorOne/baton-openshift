package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conductorone/baton-openshift/pkg/config"
	"github.com/conductorone/baton-openshift/pkg/connector"
	configSchema "github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/connectorrunner"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var version = "dev"

func main() {
	ctx := context.Background()

	_, cmd, err := configSchema.DefineConfiguration(
		ctx,
		"baton-openshift",
		getConnector,
		config.Configuration,
		connectorrunner.WithDefaultCapabilitiesConnectorBuilder(&connector.Connector{}),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmd.Version = version

	err = cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getConnector(ctx context.Context, cfg *config.Openshift) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)

	if err := config.ValidateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	var restConfig *rest.Config
	var err error
	kubeConfigPath := cfg.KubeConfig
	if kubeConfigPath == "" {
		l.Debug("no kubeconfig file specified. trying in-cluster config")
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to build configuration from in-cluster config, error: %w", err)
		}
	} else {
		restConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, fmt.Errorf("unable to build configuration from kubeconfig file, error: %w", err)
		}
	}

	cb, err := connector.New(ctx, cfg.Namespace, restConfig)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	connectorServer, err := connectorbuilder.NewConnector(ctx, cb)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}
	return connectorServer, nil
}
