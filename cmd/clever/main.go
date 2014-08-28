package main

import (
	"code.google.com/p/goauth2/oauth"
	"flag"
	"github.com/Clever/clever-to-csv"
	csvSink "github.com/azylman/optimus/sinks/csv"
	"github.com/azylman/optimus/transformer"
	clevergo "gopkg.in/Clever/clever-go.v1"
	"log"
	"net/url"
	"os"
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

func main() {
	host := flag.String("host", "https://api.clever.com", "base URL of Clever API")
	token := flag.String("token", "", "API token to use for authentication")
	output := flag.String("output", "csv", "output method. supported options: csv")
	flag.Parse()

	for _, required := range []*string{host, token, output} {
		if len(*required) == 0 {
			flag.Usage()
			os.Exit(1)
		}
	}

	if *output != "csv" {
		log.Fatal("supported output methods: csv")
	}

	if len(flag.Args()) < 2 {
		log.Fatal("need at least two arguments: endpoint action options")
	}

	endpoint := flag.Args()[0]
	if !validEndpoint(endpoint) {
		log.Fatalf("unknown endoint. supported endpoints: %#v", acceptedEndpoints)
	}
	action := flag.Args()[1]
	if action != "list" {
		log.Fatal("unknown action. supported actions: list")
	}

	transport := &oauth.Transport{
		Token: &oauth.Token{AccessToken: *token},
	}
	client := transport.Client()
	clever := clevergo.New(client, *host)

	var params url.Values
	if len(flag.Args()) > 2 {
		params = url.Values{"where": []string{flag.Args()[2]}}
	}

	t := transformer.New(clevertable.New(endpoint, params, clever)).
		Map(clevertable.FlattenRow).
		Map(clevertable.StringifyArrayVals).
		Table()
	if err := csvSink.New(t, endpoint+".csv"); err != nil {
		log.Fatal(err)
	}
}
