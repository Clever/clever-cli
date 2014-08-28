package main

import (
	"code.google.com/p/goauth2/oauth"
	"flag"
	csvSink "github.com/azylman/optimus/sinks/csv"
	"github.com/azylman/optimus/transformer"
	clevergo "gopkg.in/Clever/clever-go.v1"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	host := flag.String("host", "https://api.clever.com", "base URL of Clever API")
	endpoints := flag.String("endpoints", "schools,sections,students,teachers", "comma-delimited list of endpoints to pull")
	token := flag.String("token", "", "API token to use for authentication")
	flag.Parse()

	for _, required := range []*string{host, endpoints, token} {
		if len(*required) == 0 {
			flag.Usage()
			os.Exit(1)
		}
	}
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: *token},
	}
	client := t.Client()
	clever := clevergo.New(client, *host)

	wg := sync.WaitGroup{}
	errors := make(chan error)
	go func() {
		wg.Wait()
		close(errors)
	}()
	for _, endpoint := range strings.Split(*endpoints, ",") {
		wg.Add(1)
		go func(endpoint string) {
			t := transformer.New(NewCleverTable(endpoint, nil, clever)).
				Map(FlattenRow).
				Map(StringifyArrayVals).
				Table()
			if err := csvSink.New(t, endpoint+".csv"); err != nil {
				errors <- err
			}
			wg.Done()
		}(endpoint)
	}
	for err := range errors {
		log.Fatal(err)
	}
}
