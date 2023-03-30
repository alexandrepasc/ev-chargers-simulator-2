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

#### Git

To be easy to find the changes done to the projecto code in each iteration some rules must be in place. When commiting code changes the commit must have an objective that could be translated into text, to achive that focus the work in one objective and when that is fulfilled commit the work done.

The some rules are set in place to help understand the commits that are done into the repository. In first plance define the type of work as `feat`, `fix`, or `docs` with an `:` after it (the list defined for the types are list down). Then tag the work task from the Task Managing tool with `#` and the *ID* of the task (ex. `feat: #123`). After that we should resume the work done with objective message or title and text. The total amount of characters for the commit is set to **122** characters, including the *type* and the *ID*.

##### Commit types:
- **fix** - a bug fix
- **feat** - a new feature was introduced
- **docs** - update the documentation, code or readme
- **test** - add or correct tests
- **chore** - changes that do not afect the code or tests (updating dependencies)
- **ci** - related to the continuos integration
- **perf** - performance improvements
- **style** - changes to code formatting
- **build** - related to the build system of the application

## Build and Test
## Contribute