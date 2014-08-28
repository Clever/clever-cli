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
$ clever-cli --token=DEMO_TOKEN sections list '{"subject":"math"}'
```

## Usage

`clever-cli [options] endpoint action [query]`

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

### Query

A string that modifies the request to Clever. Meaning varies based on `action`.
For `list`, it is a JSON-string that will be sent as a `where` query.
