![Baton Logo](./docs/images/baton-logo.png)

# `baton-openshift` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-openshift.svg)](https://pkg.go.dev/github.com/conductorone/baton-openshift) ![main ci](https://github.com/conductorone/baton-openshift/actions/workflows/main.yaml/badge.svg)

`baton-openshift` is a connector for built using the [Baton SDK](https://github.com/conductorone/baton-sdk).

Check out [Baton](https://github.com/conductorone/baton) to learn more the project in general.

# Getting Started

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-openshift
baton-openshift
baton resources
```

## Requirements

This connector uses your kube-config file. When you use `oc` to authenticate into the cluster, this kube-config file is updated for you. Thus, run:

```
oc login -u <your-user> https://<host-for-your-cluster>:6443
```

Now you can synchronize the users of the cluster and their roles for a given namespace and groups:

```
baton-openshift --kube-config /home/example/.kube/config --namespace example-namespace
```

## docker

```
docker run --rm -v $(pwd):/out -v /home/example/.kube/config:/user/config -e BATON_DOMAIN_URL=domain_url -e BATON_API_KEY=apiKey -e BATON_USERNAME=username ghcr.io/conductorone/baton-openshift:latest -f "/out/sync.c1z --kube-config /user/config --namespace example-namespace"

docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-openshift/cmd/baton-openshift@main

baton-openshift

baton resources
```

# Data Model

`baton-openshift` will pull down information about the following resources:
- Users
- Roles
- Groups

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually
building spreadsheets. We welcome contributions, and ideas, no matter how
small&mdash;our goal is to make identity and permissions sprawl less painful for
everyone. If you have questions, problems, or ideas: Please open a GitHub Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-openshift` Command Line Usage

```
baton-openshift

Usage:
  baton-openshift [flags]
  baton-openshift [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --client-id string       The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string   The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string            The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                   help for baton-openshift
      --kube-config string     required: ($BATON_KUBE_CONFIG)
      --log-format string      The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string       The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
      --namespace string       required: ($BATON_NAMESPACE)
  -p, --provisioning           This must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
      --skip-full-sync         This must be set to skip a full sync ($BATON_SKIP_FULL_SYNC)
      --ticketing              This must be set to enable ticketing support ($BATON_TICKETING)
  -v, --version                version for baton-openshift

Use "baton-openshift [command] --help" for more information about a command.
```
