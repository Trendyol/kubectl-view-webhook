# kubectl-view-webhook

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Trendyol/kubectl-view-webhook)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/Trendyol/kubectl-view-webhook)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/Trendyol/kubectl-view-webhook/goreleaser)
![goreleaser](https://github.com/Trendyol/kubectl-view-webhook/workflows/goreleaser/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/Trendyol/kubectl-view-webhook)](https://goreportcard.com/report/github.com/Trendyol/kubectl-view-webhook)
![GitHub](https://img.shields.io/github/license/Trendyol/kubectl-view-webhook)

Visualize your webhook configurations in Kubernetes.

![Output](https://raw.githubusercontent.com/Trendyol/kubectl-view-webhook/master/.res/output.png)

## Installation
> Go binaries are automatically built with each release by GoReleaser. These can be accessed on the GitHub [releases page](https://github.com/Trendyol/kubectl-view-webhook/releases) for this project.

There are several ways to install view-webhook. The recommended installation method is via krew.
### Via Go
```bash
$ go get https://github.com/Trendyol/kubectl-view-webhook
```

### Via source code

Option 1 (if you have a Go compiler and want to tweak the code):
```bash
$ git clone https://github.com/Trendyol/kubectl-view-webhook
$ cd kubectl-view-webhook
$ go build .
```

### Via krew (Not available yet)
Krew is a kubectl plugin manager. If you have not yet installed krew, get it at [kubernetes-sigs/krew](https://github.com/kubernetes-sigs/krew). Then installation is as simple as :

```bash
$ kubectl krew install view-webhook
$ kubectl view-webhook --help
```

### Table details
```bash
| Kind                                      | Name                       | Webhook             | Service                    | Resources                                    | Operations                                  | Remaing Day        | Active Namespaces    |
|-------------------------------------------|----------------------------|---------------------|----------------------------|----------------------------------------------|---------------------------------------------|--------------------|----------------------|
| Type of the webhook (Mutating/Validating) | Name of the webhook config | Name of the webhook | service details of webhook | Kubernetes Resources which webhook interests | Kubernetes Operations(CREATE/UPDATE/DELETE) | Cert Remaining Day | Activated namespaces |
```

## Usage
By default, view-webhook will display all the Validating&Mutating Admission webhooks that available on your cluster.Also, you can get the detail of each one of them by giving its name.

```bash
$ kubectl view-webhook [flags]
$ kubectl view-webhook NAME [flags]
```

## License

This repository is available under the [Apache License 2.0](https://github.com/Trendyol/kubectl-view-webhook/blob/master/LICENSE).
