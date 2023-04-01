# ev-chargers-simulator-2
## Introduction
## Getting Started

### Code linter

To maintain the code homogenous between all the files and all the people envolved developing this application, some norms should be used.
For this the [golangci-lint](https://github.com/golangci/golangci-lint) package was selected, and it needs to be installed separated from the rest of the project dependecies.

Add the *go* bin folder `~/go/bin` to the `PATH`, editing the `~/.bashrc`, add the path, and apply the changes with `source ~/.bashrc`.

Run the linter:
- `./scripts/run_linter.sh`
- `golangci-lint run ./...`

The configuration for this is done in the file `.golangci.yml`.

### Git

To be easy to find the changes done to the projecto code in each iteration some rules must be in place. When commiting code changes the commit must have an objective that could be translated into text, to achive that focus the work in one objective and when that is fulfilled commit the work done.

The some rules are set in place to help understand the commits that are done into the repository. In first plance define the type of work as `feat`, `fix`, or `docs` with an `:` after it (the list defined for the types are list down). Then tag the work task from the Task Managing tool with `#` and the *ID* of the task (ex. `feat: #123`). After that we should resume the work done with objective message or title and text. The total amount of characters for the commit is set to **122** characters, including the *type* and the *ID*.

#### Commit types:
- **fix** - a bug fix
- **feat** - a new feature was introduced
- **docs** - update the documentation, code or readme
- **test** - add or correct tests
- **chore** - changes that do not afect the code or tests (updating dependencies)
- **ci** - related to the continuos integration
- **perf** - performance improvements
- **style** - changes to code formatting
- **build** - related to the build system of the application

With what was set in place for this project, the best way to do the commits is using the commad line, because there is a question that will be asked before the commit is done. So it is a good idea to review the commands needed to do this part.

Add a file to *stage*:
- `git add {path_to_file}`

Commit the staged changes:
- `git commit -m "message"`
- `git commit -m "title" -m "description"`

### Vulnerabilities

Since the project uses opensource packages and with the constant finding of new exploits, we should be more careful with what the project is using to prevent creating security issues to whom uses it.

There is an iniciative from the *Go security team* to have an database where the package maintainers report the vulnerabilities found in their code, and a tool to run against the code to identify the issues.

To install the tool:
`go install golang.org/x/vuln/cmd/govulncheck@latest`

To run it:
`govulncheck ./...`

Links:
[Go blog](https://go.dev/blog/vuln)
[Vuln documentation](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)

### Versioning

To use a automatic way to generate the change logs for each release we need to keep track of the version of the project. At the moment the version is set in the variable `Version` in the file `common/constants.go` and needs to be changed manually, and need to be updated when each task is finished before the merge request is created.

At the moment the verification of the version set it is being forced every commit.

### Setup development environment

To be easier to maintain the *rules* defined above there are some steps to do before starting to contribute to the project. The linter tool need to be installed in the contributer machine and the git hooks should be set in place, to do that we only need to execute the following script.

- `./scripts/setup_environment.sh`

It is better to execute the script in the root project folder, and the *go* bin path needs to be added to the `PATH` variable. After this your go to go.

## Build and Test
## Contribute

To contribute to the project some "rules" must be set