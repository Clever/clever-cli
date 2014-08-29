# clever-cli

`clever-cli` is a command line tool to access the Clever API.

## Installing

```shell
go get github.com/Clever/clever-cli
```

## Examples

```shell
$ clever-cli --token=DEMO_TOKEN teachers list
$ clever-cli --token=DEMO_TOKEN sections list
$ clever-cli --token=DEMO_TOKEN sections list --where='{"subject":"math"}'
```

## Usage

`clever-cli [options] endpoint action [action options]`

### Options

There is one required command line flag:
  - token: API token to use for authentication

And three optional ones:
  - help=false: if true, display help and exit
  - host="https://api.clever.com": base URL of Clever API
  - output="csv": output method. currently CSV is the only option

### Endpoint

Which endpoint to query in the Clever API.
Valid options are `students`, `schools`, `sections`, or `teachers`.

### Action

What you want to do with that endpoint. Valid options are list, which returns all the results.

### Action options

A set of optional command line flags that modify the request to the Clever API.
Varies based on action type.

#### List

  - where="": a JSON-stringified where query parameter
