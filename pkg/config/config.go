package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	KubeConfig = field.StringField(
		"kube-config",
		field.WithRequired(false),
		field.WithDescription("Path to kubeconfig file"),
		field.WithDisplayName("Kube Config"),
	)
	Namespace = field.StringField(
		"namespace",
		field.WithDefaultValue("default"),
		field.WithDescription("Kubernetes namespace"),
		field.WithDisplayName("Namespace"),
	)

	// FieldRelationships defines relationships between the fields.
	FieldRelationships = []field.SchemaFieldRelationship{}
)

//go:generate go run ./gen
var Configuration = field.NewConfiguration([]field.SchemaField{
	KubeConfig,
	Namespace,
}, field.WithConstraints(FieldRelationships...))

// ValidateConfig is run after the configuration is loaded.
func ValidateConfig(cfg *Openshift) error {
	kubeConfigPath := cfg.KubeConfig
	if kubeConfigPath == "" {
		return nil
	}

	if _, err := os.Stat(kubeConfigPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("kubeconfig file does not exist: %s", kubeConfigPath)
		}
		return fmt.Errorf("unable to stat kubeconfig file (%s): %w", kubeConfigPath, err)
	}
	return nil
}
