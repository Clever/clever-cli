package main

import (
	"encoding/json"
	"github.com/azylman/optimus"
	clevergo "gopkg.in/Clever/clever-go.v1"
	"net/url"
)

type cleverTable struct {
	err     error
	stopped bool
	rows    chan optimus.Row
}

func (t *cleverTable) start(endpoint string, params url.Values, clever *clevergo.Clever) {
	defer t.Stop()
	defer close(t.rows)

	paged := clever.QueryAll("/v1.1/"+endpoint, params)
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

func NewCleverTable(endpoint string, params url.Values, clever *clevergo.Clever) optimus.Table {
	t := &cleverTable{rows: make(chan optimus.Row)}
	go t.start(endpoint, params, clever)
	return t
}

func FlattenRow(row optimus.Row) (optimus.Row, error) {
	newRow := optimus.Row{}
	for key, val := range row {
		if typed, ok := val.(map[string]interface{}); !ok {
			newRow[key] = val
		} else {
			flatRow, err := FlattenRow(optimus.Row(typed))
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

func StringifyArrayVals(row optimus.Row) (optimus.Row, error) {
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
