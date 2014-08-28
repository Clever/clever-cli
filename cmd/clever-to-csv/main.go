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

func main() {
	host := flag.String("host", "https://api.clever.com", "base URL of Clever API")
	endpoint := flag.String("endpoint", "", "endpoint to download data from")
	where := flag.String("where", "", "optional where query to limit results")
	token := flag.String("token", "", "API token to use for authentication")
	flag.Parse()

	for _, required := range []*string{host, token, endpoint} {
		if len(*required) == 0 {
			flag.Usage()
			os.Exit(1)
		}
	}

	transport := &oauth.Transport{
		Token: &oauth.Token{AccessToken: *token},
	}
	client := transport.Client()
	clever := clevergo.New(client, *host)

	var params url.Values
	if *where != "" {
		params = url.Values{"where": []string{*where}}
	}

	t := transformer.New(clevertable.New(*endpoint, params, clever)).
		Map(clevertable.FlattenRow).
		Map(clevertable.StringifyArrayVals).
		Table()
	if err := csvSink.New(t, *endpoint+".csv"); err != nil {
		log.Fatal(err)
	}
}
