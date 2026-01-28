# ok

A unified task runner that discovers and executes tasks from multiple build tools.

## Overview

`ok` automatically discovers tasks from various build tool configuration files in your project (like `Makefile` and `package.json`) and provides a single interface to run them all. No more remembering whether to use `make build`, `npm run build`, or another tool-specific command.

## Features

- **Automatic task discovery** - Scans your project for supported build tool files
- **Unified interface** - Run tasks from any tool using the same command
- **Multiple tool support** - Currently supports:
  - Just (Justfile)
  - Make (Makefile)
  - NPM (package.json scripts)
  - Yarn (yarn.lock + package.json scripts)
  - Rake (Rakefile)
  - Task (Taskfile.yml / Taskfile.yaml)

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

Run `ok` without arguments to see all discovered tasks:

```bash
ok
```

This will show tasks from all supported tools in your project.

### Run a task

```bash
ok <task-name> [args...]
```

For example:
```bash
ok build
ok test --verbose
ok start
```

`ok` will automatically find the task and run it with the appropriate tool, passing along any additional arguments.

### Just
- Discovers recipes from `Justfile`
- Runs tasks using `just <recipe>`

### Make
- Discovers targets from `Makefile`
- Runs tasks using `make <target>`

### NPM
- Discovers scripts from `package.json`
- Runs tasks using `npm run <script>`

### Yarn
- Discovers scripts from `package.json` when `yarn.lock` is present
- Runs tasks using `yarn run <script> -- [args...]`

### Rake
- Discovers tasks from `Rakefile`
- Runs tasks using `rake <task>`

### Task
- Discovers tasks from `Taskfile.yml` / `Taskfile.yaml`
- Runs tasks using `task <task> -- [args...]` (extra args become `CLI_ARGS`)

## Development

### Requirements

- Go 1.25.6 or later

### Build

```bash
go build
```

### Test

```bash
go test ./...
```
