package clevertable

import (
	"encoding/json"
	"net/url"

	clevergo "gopkg.in/Clever/clever-go.v1"
	"gopkg.in/Clever/optimus.v3"
)

type cleverTable struct {
	err     error
	stopped bool
	rows    chan optimus.Row
}

func (t *cleverTable) startList(endpoint string, params url.Values, clever *clevergo.Clever) {
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

func (t *cleverTable) startGet(endpoint, id string, clever *clevergo.Clever) {
	defer t.Stop()
	defer close(t.rows)

	resp := &struct {
		Data optimus.Row
	}{}
	if err := clever.Query("/v1.1/"+endpoint+"/"+id, nil, &resp); err != nil {
		t.err = err
		return
	}
	t.rows <- resp.Data
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

// NewList creates a Table that reads from Clever's Paged API, using the specified parameters.
func NewList(endpoint string, params url.Values, clever *clevergo.Clever) optimus.Table {
	t := &cleverTable{rows: make(chan optimus.Row)}
	go t.startList(endpoint, params, clever)
	return t
}

// NewGet creates a Table that reads a single object from Clever's API, using the specified parameters.
func NewGet(endpoint string, id string, clever *clevergo.Clever) optimus.Table {
	t := &cleverTable{rows: make(chan optimus.Row)}
	go t.startGet(endpoint, id, clever)
	return t
}

// FlattenRow takes a Row of nested maps and converts it into a flat Row, with the keys deepened.
// For example, {"key1": {"key2": "val"}} would become {"key1.key2": "val"}.
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

// StringifyArrayVals takes a Row and converts any arrays in the values into a JSON-marshalled
// representation of the array.
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
