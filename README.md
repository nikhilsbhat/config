# config


[![Go Report Card](https://goreportcard.com/badge/github.com/nikhilsbhat/config)](https://goreportcard.com/report/github.com/nikhilsbhat/config) [![shields](https://img.shields.io/badge/license-apache%20v2-blue)](https://github.com/nikhilsbhat/config/blob/master/LICENSE)


A utility which help to switch between multiple [Kubernetes](https://kubernetes.io/) of [GKE](https://cloud.google.com/kubernetes-engine/) clusters.

## Introduction

It is difficult to switch context of different kubernets clusters hosted in GCP projects.
If one has to connect cluster using gcloud, will end up runnig multiple gcloud commands and is painful task.

Yeah GCP has a option of cloud shell, where one can connect to the cluster hassle-free. Its little hard if we have to connect locally from our machines.

Config solves exactly the same thing, by letting one to switch the cluster in one command. At a stage it's interactive shell helps one in selection of the cluster they want to switch.

## Requires

This isn't a standalone tool it still depends on few things like [`gcloud`](https://cloud.google.com/sdk/gcloud/). But makes life lot easier handling it.
* [`gcloud`](https://cloud.google.com/sdk/install) version 253.0.0 or higher (tested)

## Installation

```golang
go get -u github.com/nikhilsbhat/config
go build
```
Use the executable just like any other go-cli application.

If incase few to use this in your piece of code import package in your code.
```golang
import (
    "github.com/nikhilsbhat/config"
)
```

### `knife mediawiki stack create`

```bash
    knife mediawiki stack create (options)
```
