package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conductorone/baton-openshift/pkg/connector"
	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var version = "dev"

func main() {
	ctx := context.Background()

	_, cmd, err := config.DefineConfiguration(
		ctx,
		"baton-openshift",
		getConnector,
		field.Configuration{
			Fields: ConfigurationFields,
		},
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

func getConnector(ctx context.Context, v *viper.Viper) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)

	if err := ValidateConfig(v); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	var config *rest.Config
	var err error
	kubeConfigPath := v.GetString(kubeConfig.FieldName)
	if kubeConfigPath == "" {
		l.Debug("no kubeconfig file specified. trying in-cluster config")
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to build configuration from in-cluster config, error: %w", err)
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, fmt.Errorf("unable to build configuration from kubeconfig file, error: %w", err)
		}
	}

	cb, err := connector.New(ctx, v.GetString(namespace.FieldName), config)
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
