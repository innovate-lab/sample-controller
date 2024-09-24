# sample-controller
This repository implements a simple controller for watching Foo resources as defined with a CustomResourceDefinition (CRD).

Note: go-get or vendor this package as k8s.io/sample-controller.

This particular example demonstrates how to perform basic operations such as:

- How step by step to write a controller with the controller-runtime as well as code-generator
- How to register a new custom resource (custom resource type) of type Foo using a CustomResourceDefinition.
- How to create/get/list instances of your new resource type Foo.
- How to setup a controller on resource handling create/update/delete events.

## Prerequisites
- golang 1.23.0
- code-client v0.31.1

## Steps
