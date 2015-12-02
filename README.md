# clever-cli

`clever-cli` is a command line tool to access the Clever API.

## Installing

You can install an official release from the [Github releases page](https://github.com/Clever/clever-cli/releases) or install from source via:

```shell
go get github.com/Clever/clever-cli
```

## Examples

```shell
$ clever-cli --token=DEMO_TOKEN teachers list
$ clever-cli --token=DEMO_TOKEN sections list
$ clever-cli --token=DEMO_TOKEN sections list --where='{"subject":"math"}'
$ clever-cli --token=DEMO_TOKEN teachers get EXAMPLEID
```

## Usage

`clever-cli [options] endpoint action [action options]`

### Options

There is one required command line flag:
  - token: API token to use for authentication

And three optional ones:
  - help=false: if true, display help and exit
  - host="https://api.clever.com": base URL of Clever API
  - output="csv": output method. supported options: csv, json

### Endpoint

Which endpoint to query in the Clever API.
Valid options are `students`, `schools`, `sections`, or `teachers`.

### Action

What you want to do with that endpoint. Valid options are list (which returns all the results), and get (which returns a specific object by Clever ID).

### Action options

A set of optional command line flags that modify the request to the Clever API.
Varies based on action type.

#### List

  - where="": a JSON-stringified where query parameter

#### Get

  - Get takes a single positional argument, which is the Clever ID of the object you wish to get.

## Local Development

Set this repository up in the [standard location](https://golang.org/doc/code.html) in your `GOPATH`, i.e. `$GOPATH/src/github.com/Clever/clever-cli`.
Once this is done, `make test` runs the tests.

The release process requires a cross-compilation toolchain.
[`gox`](https://github.com/mitchellh/gox) can install the toolchain with one command: `gox -build-toolchain`.
From there you can build release tarballs for different OS and architecture combinations with `make release`.

### Rolling an official release

Official releases are listed on the [releases](https://github.com/Clever/clever-cli/releases) page.
To create an official release:

1. On `master`, bump the version in the `VERSION` file in accordance with [semver](http://semver.org/).
You can do this with [`gitsem`](https://github.com/clever/gitsem), but make sure not to create the tag, e.g. `gitsem -tag=false patch`.

2. Push the change to Github. Drone will automatically create a release for you.


## Vendoring

Please view the [dev-handbook for instructions](https://github.com/Clever/dev-handbook/blob/master/golang/godep.md).
