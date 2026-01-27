package main

import (
	cfg "github.com/conductorone/baton-openshift/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("openshift", cfg.Configuration)
}
