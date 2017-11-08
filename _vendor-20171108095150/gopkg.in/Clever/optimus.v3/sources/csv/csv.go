package csv

import (
	"encoding/csv"
	"io"

	"gopkg.in/Clever/optimus.v3"
)

type table struct {
	err     error
	rows    chan optimus.Row
	stopped bool
}

func (t *table) start(reader *csv.Reader) {
	defer t.Stop()
	defer close(t.rows)

	headers, err := reader.Read()
	if err != nil {
		t.handleErr(err)
		return
	}

	reader.FieldsPerRecord = len(headers)
	for {
		if t.stopped {
			break
		}
		line, err := reader.Read()
		if err != nil {
			t.handleErr(err)
			return
		}
		t.rows <- convertLineToRow(line, headers)
	}
}

func (t table) Rows() <-chan optimus.Row {
	return t.rows
}

func (t table) Err() error {
	return t.err
}

func (t *table) Stop() {
	t.stopped = true
}

func (t *table) handleErr(err error) {
	if err != io.EOF {
		t.err = err
	}
}

func convertLineToRow(line []string, headers []string) optimus.Row {
	row := optimus.Row{}
	for i, header := range headers {
		row[header] = line[i]
	}
	return row
}

// New returns a new Table that scans over the rows of a CSV.
func New(in io.Reader) optimus.Table {
	return NewWithCsvReader(csv.NewReader(in))
}

// NewWithCsvReader returns a new Table that scans over the rows from the csv reader.
func NewWithCsvReader(reader *csv.Reader) optimus.Table {
	table := &table{
		rows: make(chan optimus.Row),
	}
	go table.start(reader)
	return table
}
