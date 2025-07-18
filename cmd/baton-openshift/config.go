package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/spf13/viper"
)

var (
	kubeConfig = field.StringField("kube-config", field.WithRequired(false))
	namespace  = field.StringField("namespace", field.WithDefaultValue("default"))
)

var (
	// ConfigurationFields defines the external configuration required for the
	// connector to run. Note: these fields can be marked as optional or
	// required.
	ConfigurationFields = []field.SchemaField{kubeConfig, namespace}

	// FieldRelationships defines relationships between the fields listed in
	// ConfigurationFields that can be automatically validated. For example, a
	// username and password can be required together, or an access token can be
	// marked as mutually exclusive from the username password pair.
	FieldRelationships = []field.SchemaFieldRelationship{}
)

// ValidateConfig is run after the configuration is loaded, and should return an
// error if it isn't valid. Implementing this function is optional, it only
// needs to perform extra validations that cannot be encoded with configuration
// parameters.
func ValidateConfig(v *viper.Viper) error {
	kubeConfigPath := v.GetString(kubeConfig.FieldName)
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
