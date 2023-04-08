# ev-chargers-simulator-2
[![PreCheckLinter](https://github.com/alexandrepasc/ev-chargers-simulator-2/actions/workflows/on-push.yml/badge.svg)](https://github.com/alexandrepasc/ev-chargers-simulator-2/actions/workflows/on-push.yml)

[![PreCheckSecurity](https://github.com/alexandrepasc/ev-chargers-simulator-2/actions/workflows/on-merge-request.yml/badge.svg)](https://github.com/alexandrepasc/ev-chargers-simulator-2/actions/workflows/on-merge-request.yml)

[![TestsAndVersion](https://github.com/alexandrepasc/ev-chargers-simulator-2/actions/workflows/on-merge.yml/badge.svg)](https://github.com/alexandrepasc/ev-chargers-simulator-2/actions/workflows/on-merge.yml)

[![Release](https://github.com/alexandrepasc/ev-chargers-simulator-2/actions/workflows/release.yml/badge.svg)](https://github.com/alexandrepasc/ev-chargers-simulator-2/actions/workflows/release.yml)

## Introduction
To help the development and test some type of applications, like an electric management system or a charging station operator, giving the ability to have a simulated electric vehicle charger that could be controlled by the user.

This project is developed using *go* and at the moment using `goroutine` to generate the simulators.

This development has in mind the ability to run multiple simulators, with different communication protocols, be able to control the data that the simulators sends, to retrieve metrics, and to run in a container.
## Getting Started

### Code linter

To maintain the code homogeneous between all the files and all the people involved developing this application, some norms should be used.
For this the [golangci-lint](https://github.com/golangci/golangci-lint) package was selected, and it needs to be installed separated from the rest of the project dependencies.

Add the *go* bin folder `~/go/bin` to the `PATH`, editing the `~/.bashrc`, add the path, and apply the changes with `source ~/.bashrc`.

Run the linter:
- `./scripts/run_linter.sh`
- `golangci-lint run ./...`

The configuration for this is done in the file `.golangci.yml`.

### Git

To be easy to find the changes done to the project code in each iteration some rules must be in place. When committing code changes the commit must have an objective that could be translated into text, to achieve that focus the work in one objective and when that is fulfilled commit the work done.

The some rules are set in place to help understand the commits that are done into the repository. In first place define the type of work as `feat`, `fix`, or `docs` with an `:` after it (the list defined for the types are list down). Then tag the work task from the Task Managing tool with `#` and the *ID* of the task (ex. `feat: #123`). After that we should resume the work done with objective message or title and text. The total amount of characters for the commit is set to **122** characters, including the *type* and the *ID*.

#### Commit types:
- **fix** - a bug fix
- **feat** - a new feature was introduced
- **docs** - update the documentation, code or readme file(s)
- **test** - add or correct tests
- **chore** - changes that do not affect the code or tests (updating dependencies)
- **ci** - related to the continuous integration
- **perf** - performance improvements
- **style** - changes to code formatting
- **build** - related to the build system of the application

With what was set in place for this project, the best way to do the commits is using the command line, because there is a question that will be asked before the commit is done. So it is a good idea to review the commands needed to do this part.

Add a file to *stage*:
- `git add {path_to_file}`

Commit the staged changes:
- `git commit -m "message"`
- `git commit -m "title" -m "description"`

### Vulnerabilities

Since the project uses open source packages and with the constant finding of new exploits, we should be more careful with what the project is using to prevent creating security issues to whom uses it.

There is an initiative from the *Go security team* to have an database where the package maintainers report the vulnerabilities found in their code, and a tool to run against the code to identify the issues.

To install the tool:
`go install golang.org/x/vuln/cmd/govulncheck@latest`

To run it:
`govulncheck ./...`

Links:
[Go blog](https://go.dev/blog/vuln)
[Vuln documentation](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)

### Version

To use a automatic way to generate the change logs for each release we need to keep track of the version of the project. At the moment the version is set in the variable `Version` in the file `common/constants.go` and needs to be changed manually, and need to be updated when each task is finished before the merge request is created.

At the moment the verification of the version set it is being forced every commit.

### Setup development environment

To be easier to maintain the *rules* defined above there are some steps to do before starting to contribute to the project. The linter tool need to be installed in the machine and the git hooks should be set in place, to do that we only need to execute the following script.

- `./scripts/setup_environment.sh`

It is better to execute the script in the root project folder, and the *go* bin path needs to be added to the `PATH` variable. After this your go to go.

## Build and Test
To build and test the project the *go* needs to be installed in the machine, to do this we could use the information in this [page](https://go.dev/doc/install).

To be able to run it in containerized mode we the [Docker Engine](https://docs.docker.com/engine/install/) needs to be installed.

Install project dependencies:
- `go get .`

Run the project locally:
- `go run .`

Install a go package:
- `go install {package_name}`

Run the tests:
- `go test ./...`

To reduce the issues when committing new code, we should run the linter previously to do the commit using:
- `./scripts/run_linter.sh`

The commits need to be execute in the command line, since there are some git hooks that require feedback from the user.
## Contribute

To contribute to the project some "rules" must be set. So to add a change or a fix to the code one *issue* must be created, with a clear objective, the relevant information needed to fulfil the objective.

Who contribute to the project needs to have access to the project and the project [board](https://github.com/users/alexandrepasc/projects/4/views/1) where the issues are managed.

A new issue is created in the board and it will be added to the *New* column, after adding the description with the necessary information and added to the project, it will be moved to the *Backlog* column. In this stage the issue needs to be review and if approved moved to *Ready*.

This behaviour will be used to bug fixes to.

Every new feature or change should have in mind the objective of the project, and all the people contributing in it should be involved in the decision.

The code added to the project need to have in mind the way it was set and not doing implementations using other way or approach without a discussion with the rest of the involved. There is a linter set in the pipeline to prevent some miss behave.
