package main

import (
	"errors"
	"os"

	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/spf13/viper"
)

var (
	kubeConfig = field.StringField("kube-config", field.WithRequired(false))
	namespace  = field.StringField("namespace", field.WithRequired(true))
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
	kubeConfigLocation := v.GetString(kubeConfig.FieldName)
	// check if file exists
	if _, err := os.Stat(kubeConfigLocation); errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
