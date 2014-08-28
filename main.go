package main

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"flag"
	"github.com/azylman/optimus"
	csvSink "github.com/azylman/optimus/sinks/csv"
	"github.com/azylman/optimus/transformer"
	clevergo "gopkg.in/Clever/clever-go.v1"
	"log"
	"os"
	"strings"
	"sync"
)

type cleverTable struct {
	err     error
	stopped bool
	rows    chan optimus.Row
}

func (t *cleverTable) start(endpoint string, clever *clevergo.Clever) {
	defer t.Stop()
	defer close(t.rows)

	paged := clever.QueryAll("/v1.1/"+endpoint, nil)
	for !t.stopped && paged.Next() {
		row := optimus.Row{}
		if err := paged.Scan(&row); err != nil {
			t.err = err
			break
		}
		t.rows <- row
	}
	if err := paged.Error(); err != nil {
		t.err = err
	}
}

func (t *cleverTable) Rows() <-chan optimus.Row {
	return t.rows
}

func (t *cleverTable) Err() error {
	return t.err
}

func (t *cleverTable) Stop() {
	if t.stopped {
		return
	}
	t.stopped = true
}

func newCleverTable(endpoint string, clever *clevergo.Clever) optimus.Table {
	t := &cleverTable{rows: make(chan optimus.Row)}
	go t.start(endpoint, clever)
	return t
}

func flattenRow(row optimus.Row) (optimus.Row, error) {
	newRow := optimus.Row{}
	for key, val := range row {
		if typed, ok := val.(map[string]interface{}); !ok {
			newRow[key] = val
		} else {
			flatRow, err := flattenRow(optimus.Row(typed))
			if err != nil {
				return nil, err
			}
			for partKey, val := range flatRow {
				newRow[key+"."+partKey] = val
			}
		}
	}
	return newRow, nil
}

func stringifyArrayVals(row optimus.Row) (optimus.Row, error) {
	newRow := optimus.Row{}
	for key, val := range row {
		if typed, ok := val.([]interface{}); !ok {
			newRow[key] = val
		} else {
			// convert
			bytes, err := json.Marshal(typed)
			if err != nil {
				return nil, err
			}
			newRow[key] = string(bytes)
		}
	}
	return newRow, nil
}

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
			t := transformer.New(newCleverTable(endpoint, clever)).
				Map(flattenRow).
				Map(stringifyArrayVals).
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
