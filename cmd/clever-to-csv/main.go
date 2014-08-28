package main

import (
	"code.google.com/p/goauth2/oauth"
	"flag"
	"fmt"
	"github.com/Clever/clever-to-csv"
	csvSink "github.com/azylman/optimus/sinks/csv"
	"github.com/azylman/optimus/transformer"
	clevergo "gopkg.in/Clever/clever-go.v1"
	"log"
	"net/url"
	"os"
	"sync"
)

type endpointConf struct {
	Params url.Values
	Name   string
}

func (e *endpointConf) Set(val string) error {
	if val == "true" {
		return nil
	}
	e.Params = url.Values{"where": []string{val}}
	return nil
}

func (e *endpointConf) String() string {
	return ""
}

func (e *endpointConf) IsBoolFlag() bool {
	return true
}

var acceptedEndpoints = []string{"schools", "sections", "students", "teachers"}

func main() {
	host := flag.String("host", "https://api.clever.com", "base URL of Clever API")
	for _, endpoint := range acceptedEndpoints {
		end := &endpointConf{Name: endpoint}
		flag.Var(end, endpoint, fmt.Sprintf("if included, dumps the %s endpoint. can optionally include JSON-stringified map of 'where' query parameters", endpoint))
	}
	token := flag.String("token", "", "API token to use for authentication")
	flag.Parse()

	endpoints := []flag.Value{}
	flag.Visit(func(f *flag.Flag) {
		if f.Name != "host" && f.Name != "token" {
			endpoints = append(endpoints, f.Value)
		}
	})

	for _, required := range []*string{host, token} {
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

	wg := sync.WaitGroup{}
	errors := make(chan error)
	go func() {
		wg.Wait()
		close(errors)
	}()
	for _, end := range endpoints {
		wg.Add(1)
		go func(endp flag.Value) {
			defer wg.Done()
			endpoint, ok := endp.(*endpointConf)
			if !ok {
				errors <- fmt.Errorf("invalid flag value %#v", end)
				return
			}

			t := transformer.New(clevertable.New(endpoint.Name, endpoint.Params, clever)).
				Map(clevertable.FlattenRow).
				Map(clevertable.StringifyArrayVals).
				Table()
			if err := csvSink.New(t, endpoint.Name+".csv"); err != nil {
				errors <- err
			}
		}(end)
	}
	for err := range errors {
		log.Fatalf("got err %#v\n", err)
	}
}
