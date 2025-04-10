# krepe

Kubernetes configuration management for engineers who just want to write raw YAML without any tools getting in their way.

## Installation

```bash
go install -C krepe .
```

## Example

A krepe package is a directory containing a `krepe.yaml` file and any number of Kubernetes resource files. The `krepe.yaml` file defines the imported resources, imported package, and its pipelines.

A pipeline is a series of functions to execute across all imported resources.

To edit a resource directly without using the pipelines mechanism, this can be done in-place with your favourite text editor.

```yaml
# krepe.yaml
kind: Krepe
apiVersion: krepe.io/v1
metadata:
  name: krepe-example
imports:
  files:
    - deployment.yaml
    - service.yaml
    - ingress.yaml
  packages:
    - relativePath: ../some-other-pkg
      name: some-other-pkg
pipelines:
  - name: default
    steps:
      - function: set-labels
        configMap:
          app: my-app
```

From this basic `krepe.yaml`, we can run the `default` pipeline using `krepe run default` (`default` can be omitted). In the `default` pipeline above, we set a common label across all resources.

## CLI Usage

### `krepe -h`

```bash
Kubernetes configuration management tooling.

Usage:
  krepe [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  install     Import a package
  run         Run a pipeline on a package

Flags:
  -h, --help             help for krepe
  -C, --pkgPath string   path to the package to run (default ".")

Use "krepe [command] --help" for more information about a command.
```

### `krepe run -h`

```bash
Run a pipeline on a package.

Usage:
  krepe run [pipelineName]

Arguments:
  pipelineName  The name of the pipeline to run. This argument is optional. Default is 'default'.

Example:
  krepe run default

Usage:
  krepe run [flags]

Flags:
  -h, --help   help for run

Global Flags:
  -C, --pkgPath string   path to the package to run (default ".")
```

### `krepe install -h`

```bash
Import a pacakge.

Usage:
  krepe install [path]

Arguments:
  path  The path to the package to install. Relative paths must be relative to the package being operated on. This argument is required.

Example:
  krepe install ../some_package

Usage:
  krepe install [flags]

Flags:
  -h, --help          help for install
  -n, --name string   name override of the package being installed

Global Flags:
  -C, --pkgPath string   path to the package to run (default ".")
```
