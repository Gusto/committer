# Building

This should be as simple as running `go build` at the root of the project repo.  There's a `make committer` target, if you prefer that.

# Tests

Tests are written in standard Go (without any framework), but they do use [testify](https://github.com/stretchr/testify) for assertions.

Test data is kept in the `test` directory.

Run all tests with `go test ./...` at the root of the repo.  For convenience, `make test` will do the same thing.

# Layout

This is a fairly small project!  All the code lives in `core`.

There's an example `committer.yml` at the root of the repo.

# Dependency Management

This project uses [`dep`](https://github.com/golang/dep), but you won't have to interact with it unless you're adding or changing dependencies.  Install it with `make deps`.

Note that the `vendor` directory and `Gopkg.lock` file *are* tracked in the repo, so make sure to commit your changes to them.

# Linting

Run `make lint` to run the linter suite over the whole project.  You'll be trusted to run this before submitting changes for review until we get a CI server set up (or publish v1 of committer itself!).
