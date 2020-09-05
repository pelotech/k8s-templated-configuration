# Mutating Admission Webhook for Templated Configuration

## Overview

Often when working with teams new to Kubernetes we find their applications are not factored in a Kubernetes friendly way. One of the most immediate ways this manifests itself is configuration and secret data being mixed together. While we can address this by storing all configuration containing private data in secrets we've found this causes a number of additional challenges. Especially for teams new to Kubernetes.

Ideally we'd be able to extract just the secret content from our configuration, and at pod startup inject it into our ConfigMaps using a standard template format.

This Kubernetes controller adds that capability. By injected an init container into newly created pods based on annotations we can split configuration and secrets in a way transparent to the underlying workloads.

At present we use [envtemplate](https://github.com/orls/envtemplate) to evaluate go templated files. In the future we will aim to support go templates with sprig to somewhat mirror Helm as well as extensible template engines to allow teams to use existing template files of their preferred format.

This repository is based on [k8s-webhook-example] a production ready [Kubernetes admission webhook][k8s-admission-webhooks] example using [Kubewebhook].

## Structure

The application is mainly structured in 3 parts:

- `main`: This is where everything is created, wired, configured and set up, [cmd/k8s-webhook](cmd/k8s-webhook/main.go).
- `http`: This is the package that configures the HTTP server, wires the routes and the webhook handlers.  [internal/http/webhook](internal/http/webhook).
- Application services: These services have the domain logic of the validators and mutators:
  - [`mutation/template`](internal/mutation/template): Logic for `template.pelo.tech` webhook.

Apart from the webhook referring stuff we have other parts like:

- [Decoupled metrics](internal/metrics)
- [Decoupled logger](internal/log)
- [Application command line flags](cmd/k8s-webhook/config.go)

And finally there is an example of how we could deploy our webhooks on a production server:

- [Deploy](deploy)

## Webhooks

### `template.pelo.tech`

- Webhook type: Mutating.
- Resources affected: `pods`

This webhook adds an `initContainer` to parse secrets into volumes as template files.

TODO: Put more details here

[k8s-webhook-example]: https://github.com/pelotech/k8s-templated-configuration
[k8s-admission-webhooks]: https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/
[Kubewebhook]: https://github.com/slok/kubewebhook
