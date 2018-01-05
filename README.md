# k8s-chaos

Introduce a little chaos into your Kubernetes namespace. Randomly kill pods.
Will your application survive?

<!-- TOC depthFrom:2 depthTo:6 withLinks:1 updateOnSave:1 orderedList:0 -->

- [Install](#install)
- [Usage](#usage)
- [Disclaimer](#disclaimer)

<!-- /TOC -->

## Install

```
dep ensure
go install ./cmd/k8s-chaos
```

## Usage

k8s-chaos --namespace my-namespace --max-kill 2 --period 60 --regex my-thing

Options: **Bold** options are required.

* **namespace** - The kubernetes namespace to target.
* max-kill - The maximum number of pods to kill during each cycle.
* period - The amount of time between each cycle.
* regex - A regex that will be used to filter the list of pods in a namespace.
  If you don't supply a value, all pods in the namespace will be considered for
  deletion.

## Extras

In the [docker](./docker) folder, you will find an assortment of images that you
can use to test different failure scenarios:

* crash - An image that exits with non-zero exit status after `CRASH_DELAY` seconds.

## Disclaimer

This repo is pretty immature. Expect the possibility of breaking changes.
