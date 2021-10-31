[![Github Actions](https://github.com/kick-project/maker/workflows/Go/badge.svg?branch=master)](https://github.com/kick-project/maker/actions) [![Go Report Card](https://goreportcard.com/badge/kick-project/maker)](https://goreportcard.com/report/kick-project/maker)  [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/kick-project/maker/blob/master/LICENSE)

# maker

maker a tool that augments make with Makefile menu help, dotenv support

## Example Help

This help output

```text
Usage: make <target>

### HELP
  help - Print Help
### DEVELOPMENT
  build - Build binary
  test  - Test
```

is generated using `make help` and this `Makefile`

```Makefile
.PHONY: help build install clean test

### HELP
help: ## Print Help
	@maker --menu Makefile

### DEVELOPMENT
build: ## Build binary
	@echo "Run build ..."

test: build ## Test
	@echo "Run Tests ..."
```

## Example dotenv support

Maker supports `.env` files. The default is to load `$HOME/.env` then the
projects `.env` file.

Given the following `.env` files.

```dotenv
# Home directory "$HOME/.env" file
MYSECRET=MySeceret
```

```dotenv
# Projects "./.env" file
REPO=http://artifactory.mycompany.com/packages
```

Then the `Makefile` will print the environment variables listed in the
`_printvars` target.

**NOTE** Maker will automatically prepend the underscore "_" to the target
*before calling make again. E.G. `make target` will become `make _target` when
*it falls through to the catchall.

```Makefile
MAKEFLAGS += --no-print-directory
.PHONY: _printvars

### HELP
_printvars: ## Print variables
	@echo "HOME=${HOME}"
	@echo "REPO=${REPO}"
	@echo "MYSECRET=${MYSECRET}"

# Catch all target to wrap tasks with a single underscore prefix.
%:
	@maker $@
```

Output of `make printvars` (Notice no preceding underscore)

```
HOME=/home/username
REPO=http://artifactory.mycompany.com/packages
MYSECRET=MySeceret
```

Maker will not override any existing environment variables. The default order of
presedence is.

1. Global environment variables.
2. Home directory `~/.env` file.
3. The current folders `.env` file.