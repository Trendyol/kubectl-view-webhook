[![License][license-img]][license]
[![Go Report Card][report-card-img]][report-card]

# kubectl-view-webhook

Visualize your webhook configurations in Kubernetes.

![Output](https://raw.githubusercontent.com/Trendyol/kubectl-view-webhook/master/.res/output.png)

## Installation

### Source

Option 1 (if you have a Go compiler and want to tweak the code):
```bash
$ git clone https://github.com/Trendyol/kubectl-view-webhook
$ cd kubectl-view-webhook
$ go build .
```

## Usage

```bash
$ kubectl view-webhook [flags]
$ kubectl view-webhook NAME [flags]
```

## License

This repository is available under the [Apache License 2.0](https://github.com/Trendyol/kubectl-view-webhook/blob/master/LICENSE).

![goreleaser](https://github.com/Trendyol/kubectl-view-webhook/workflows/goreleaser/badge.svg) 
[report-card-img]: https://goreportcard.com/badge/github.com/Trendyol/kubectl-view-webhook?style=flat-square
[report-card]: https://goreportcard.com/report/github.com/Trendyol/kubectl-view-webhook

[license-img]: https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square
[license]: https://github.com/Trendyol/kubectl-view-webhook/blob/master/LICENSE
