package main

import (
	"code.google.com/p/goauth2/oauth"
	"flag"
	"fmt"
	"github.com/Clever/clever-cli/clevertable"
	csvSink "github.com/azylman/optimus/sinks/csv"
	"github.com/azylman/optimus/transformer"
	clevergo "gopkg.in/Clever/clever-go.v1"
	"log"
	"net/url"
	"os"
	"strings"
)

var acceptedEndpoints = []string{"students", "schools", "sections", "teachers"}

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
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [options] endpoint action [query]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "endpoints: %s\n\n", strings.Join(acceptedEndpoints, " "))
		fmt.Fprintf(os.Stderr, "actions: list\n\n")
		fmt.Fprintf(os.Stderr, "query:\n  - for list, a JSON-stringified where query parameter\n\n")
		fmt.Fprintf(os.Stderr, "options:\n")
		flag.PrintDefaults()
	}
	host := flag.String("host", "https://api.clever.com", "base URL of Clever API")
	token := flag.String("token", "", "API token to use for authentication (required)")
	output := flag.String("output", "csv", "output method. supported options: csv")
	help := flag.Bool("help", false, "if true, display help and exit")
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

	if *output != "csv" {
		exitWithArgError(fmt.Sprintf("'%s' is not a valid output", *output))
	}

	if len(flag.Args()) < 2 {
		exitWithArgError("need at least two arguments: endpoint and action")
	}

	endpoint := flag.Args()[0]
	if !validEndpoint(endpoint) {
		exitWithArgError(fmt.Sprintf("'%s' is not a valid endpoint", endpoint))
	}
	action := flag.Args()[1]
	if action != "list" {
		exitWithArgError(fmt.Sprintf("'%s' is not a valid action", action))
	}

	query := ""
	if len(flag.Args()) > 2 {
		query = flag.Args()[2]
	}

	transport := &oauth.Transport{
		Token: &oauth.Token{AccessToken: *token},
	}
	client := transport.Client()
	clever := clevergo.New(client, *host)

	var params url.Values
	if query != "" {
		params = url.Values{"where": []string{query}}
	}

	t := transformer.New(clevertable.New(endpoint, params, clever)).
		Map(clevertable.FlattenRow).
		Map(clevertable.StringifyArrayVals).
		Table()
	if err := csvSink.New(t, endpoint+".csv"); err != nil {
		log.Fatal(err)
	}
}
