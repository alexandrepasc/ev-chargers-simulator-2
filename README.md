# ev-chargers-simulator-2
## Introduction
## Getting Started

### Prepare environment

#### Code linter

To maintain the code homogenous between all the files and all the people envolved developing this application, some norms should be used.
For this the [golangci-lint](https://github.com/golangci/golangci-lint) package was selected, and it needs to be installed separated from the rest of the project dependecies.

Install the linter package:
- `./scripts/install_linter.sh`

Add the *go* bin folder `~/go/bin` to the `PATH`, editing the `~/.bashrc`, add the path, and apply the changes with `source ~/.bashrc`.

Run the linter:
- `./scripts/run_linter.sh`

The configuration for this is done in the file `.golangci.yml`.

## Build and Test
## Contribute