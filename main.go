package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/Clever/clever-cli/clevertable"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	clevergo "gopkg.in/Clever/clever-go.v1"
	"gopkg.in/Clever/optimus.v3"
	csvSink "gopkg.in/Clever/optimus.v3/sinks/csv"
	jsonSink "gopkg.in/Clever/optimus.v3/sinks/json"
	"gopkg.in/Clever/optimus.v3/transformer"
)

var acceptedEndpoints = []string{"students", "schools", "sections", "teachers"}
var acceptedActions = []string{"list", "get"}

func validEndpoint(endpoint string) bool {
	for _, accepted := range acceptedEndpoints {
		if endpoint == accepted {
			return true
		}
	}
	return false
}

func exitWithArgError(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n\n", msg)
	flag.Usage()
	os.Exit(1)
}

func main() {
	baseUsage := func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [options] endpoint action [action options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nendpoints: %s\n\n", strings.Join(acceptedEndpoints, " "))
		fmt.Fprintf(os.Stderr, "actions: %s\n\n", strings.Join(acceptedActions, " "))
	}
	host := flag.String("host", "https://api.clever.com", "base URL of Clever API")
	token := flag.String("token", "", "API token to use for authentication (required)")
	output := flag.String("output", "csv", "output method. supported options: csv, json")
	help := flag.Bool("help", false, "if true, display help and exit")

	// Action-specifc flags
	listFlags := flag.NewFlagSet("list options", flag.ExitOnError)
	listUsage := func() {
		fmt.Fprintln(os.Stderr, "list:")
		listFlags.PrintDefaults()
	}
	where := listFlags.String("where", "", "a JSON-stringified where query parameter")

	flag.Usage = func() {
		baseUsage()
		fmt.Fprint(os.Stderr, "action options:\n\n")
		listUsage()
	}
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	for _, required := range [][2]string{{"host", *host}, {"token", *token}, {"output", *output}} {
		name, value := required[0], required[1]
		if len(value) == 0 {
			exitWithArgError(fmt.Sprintf("must provide '%s'", name))
		}
	}

	var sink optimus.Sink
	switch *output {
	case "csv":
		sink = csvSink.New(os.Stdout)
	case "json":
		sink = jsonSink.New(os.Stdout)
	default:
		exitWithArgError(fmt.Sprintf("'%s' is not a valid output", *output))
	}

	if len(flag.Args()) < 2 {
		exitWithArgError("need at least two arguments: endpoint and action")
	}

	endpoint := flag.Args()[0]
	if !validEndpoint(endpoint) {
		exitWithArgError(fmt.Sprintf("'%s' is not a valid endpoint", endpoint))
	}

	client := oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: *token,
	}))
	clever := clevergo.New(client, *host)

	action := flag.Args()[1]

	var table optimus.Table

	switch action {
	case "list":
		if len(flag.Args()) > 2 {
			listFlags.Parse(flag.Args()[2:])
		} else {
			listFlags.Parse([]string{})
		}
		var params url.Values
		if *where != "" {
			params = url.Values{"where": []string{*where}}
		}
		table = clevertable.NewList(endpoint, params, clever)
	case "get":
		if len(flag.Args()) != 3 {
			exitWithArgError(fmt.Sprintf("get action requires an <id> argument"))
		}
		id := flag.Args()[2]
		table = clevertable.NewGet(endpoint, id, clever)
	default:
		exitWithArgError(fmt.Sprintf("'%s' is not a valid action", action))
	}

	if err := transformer.New(table).
		Map(clevertable.FlattenRow).
		Map(clevertable.StringifyArrayVals).
		Sink(sink); err != nil {
		log.Fatal(err)
	}
}
