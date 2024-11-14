# ================= DEPRECATED =================

### Gusto has standardized on using [Lefthook](https://github.com/evilmartians/lefthook) to manage git hooks.

# ================= DEPRECATED =================

# Committer

Committer is a parallel runner for running static analysis in a git hook created. It supports:
- Feeding in only the relevant changed files to the analysis commands
- Running in "fix" mode, and automatically correcting files before commit
  - By default, the files will automatically be staged to be included in the commit
  - An environment variable can be set to always stop the pre-commit hook if there are changes to be staged manually
- Running any number of static analysis tools in parallel
- Helpful display of what was changed before the commit goes through

## Installation
You can run the install script by executing:
```bash
bash <(curl -s https://raw.githubusercontent.com/Gusto/committer/master/configure.sh)
```
This will:
- Add `committer` to `/usr/local/bin`
- Create a basic `committer.yml` in your current working directory to run Rubocop
- Setup your `.git/hooks/pre-commit` to run committer, if it is not already present

You can also build it from source by cloning this repository and running `go build committer.go`.

## Configuration

Committer is fed a YAML file that describes how it should run. There is an example of such file [here](https://github.com/Gusto/committer/blob/master/committer.yml). By default, it will look for a `committer.yml` in the current working directory.

### Structure

There is a required top level `tasks` key, containing an array of Task objects.

- `name`: Used to describe the task and provide a title for the output.
- `command`: The command to run when executing the linter. This command should be able to receive a list of files as an argument.
- `files`: A Regex describing the type of files that should be fed into this command.
- `excludefilenames`: Set this to true if the command should not receive the list of changed files (i.e. it is intended to run on the whole repo).
- `fix` (Optional)
  - `command`: The command to auto-correct violations. If this is specified, committer will attempt to fix and stage your files.
  - `output`: A Regex to pull the relevant fixed output for display before committing.

## Running

Committer is most often run as a pre-commit hook. A typical configuration would be to have the following in your `.git/hooks/pre-commit` script:

```bash
#!/usr/bin/env bash
committer --fix
```

## Releasing and deploying a new committer version

#### Release

1. Make your changes on a branch.
2. Update `VERSION` in `committer.go` to new version number.
3. Once landed in origin/main, add a new version tag to the branch:
```shell
git fetch && git checkout main && git pull
git tag --annotate "vNEW.VERSION.HERE"
git push --follow-tags
```
4. `.github/workflows/publish.yml` will generate a new release for your tag. Follow the action in said tab.

5. The `v` is important.

#### Deploy

1. Gusto engineers may [go/deploy-committer-with-cpe-chef](https://go/deploy-committer-with-cpe-chef).

## Autofix

Committer can auto-fix everything by default. set `export COMMITTER_AUTO_FIX=true # or 1, T, t` in your
`~/.profile` or equivalent dotfile.

## Opting out of automatic staging

Committer will stage auto-corrected files by default. In order to always leave auto-corrected files unstaged for manual staging, set `export COMMITTER_SKIP_STAGE=1` in your `~/.bashrc` or equivalent.


## Best-practices
### Purely Functional Pre-Commit Hooks
Custom git pre-commit hooks should, as much as possible, be "pure functions" that simply reject / accept changes without side effects. Pre-commit hooks by default should perform no
- Intermediate stashing
- `git add`'ing
- Autocorrection
- etc

The recommended workflow:
1. `git add` your code
2. `git commit`
3. See the pre-commit hook failure
4. Manually run `committer --fix`
5. Redo `git add` and `git commit`.

We originally performed autocorrection and other modification but there were a lot of random difficult-to-debug side effects that emerged. As we add more steps the possibility for interference grows n^2 (any step can potentially interfere with any other step). Moving to a "purely functional" pre-commit hook model lets us easily add additional steps.

### No Composition
Steps defined by `committer.yml` should be completely independent. Steps have no guaranteed order, and are not intended to be composable.
