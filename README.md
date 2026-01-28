# 👌 ok

A unified task runner that discovers and executes tasks from multiple build tools.

Currently supports:
- [Just (`Justfile`)](https://github.com/casey/just)
- [Make (`Makefile`)](https://www.gnu.org/software/make/manual/make.html)
- [NPM (`package.json` scripts)](https://docs.npmjs.com/cli/v8/using-npm/scripts)
- [Yarn (`yarn.lock` + `package.json` scripts)](https://classic.yarnpkg.com/lang/en/docs/package-json/#toc-scripts)
- [Rake (`Rakefile`)](https://github.com/ruby/rake)
- [Task (`Taskfile.yml` / `Taskfile.yaml`)](https://taskfile.dev)
- [Nx (`nx.json`)](https://nx.dev/)
- Shell (*.sh / *.bash / *.zsh)

## Installation

### Homebrew

```bash
brew install broothie/ok/ok
```

### Go Install

```bash
go install github.com/broothie/ok@latest
```

## Usage

### List all available tasks

```bash
$ ok
TASK        TOOL   FILE
test-just   Just   Justfile
test-make   Make   Makefile
test-npm    NPM    package.json
test-rake   Rake   Rakefile
test-shell  Shell  test-shell.sh
test-task   Task   Taskfile.yml
test-yarn   NPM    package.json
```
### Run a task

```bash
ok <task-name> [args...]
```

`ok` will automatically find the task and run it with the appropriate tool, passing along any additional arguments.
For example:

```bash
$ ok test-npm

> test-npm
> echo 'from npm'

from npm
```

### Help

```bash
$ ok -h
ok v0.3.2

Usage:
  ok [options] <task> [task args]

Options:
  -V  --version            Print command version.                                                 (default: false)
  -h  --help               Show command help.                                                     (default: false)
  -d  --directory          Directory to run command from.                                         (default: .)
      --timeout            Command timeout.                                                       (default: 5s)
      --filter-tools --ft  Filter tools by case-insensitive name. Use commas for multiple values  (default: )
      --list-tools         List tools.                                                            (default: false)
      --load-dot-env       Pick up local .env files.                                              (default: true)
      --debug              Output debug logs.                                                     (default: false)
```
