package pipetimer

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/vitraum/golang-pipedrive"
)

// PipeWriter is used to write deals in CSV format
type PipeWriter struct {
	csv    *csv.Writer
	stages pipedrive.Stages
}

// DataProvider represents the data to be writtten in CSV format
type DataProvider interface {
	ID() string
	Age() time.Duration
	Status() string
	DealUpdates() pipedrive.DealFlowUpdates
	Value() float64
	DecisionTime() time.Time
}

// NewPipeWriter constructs a new PipeWriter object
func NewPipeWriter(out io.Writer, stages pipedrive.Stages) *PipeWriter {
	s := PipeWriter{
		csv:    csv.NewWriter(out),
		stages: stages,
	}

	return &s
}

// WriteHeader generates the CSV header
func (w *PipeWriter) WriteHeader() error {
	columnNames := []string{
		"Lead ID",
		"Status",
		"Endscheidungsdatum",
		"Lead Alter",
		"Wert",
	}
	for _, stage := range w.stages {
		columnNames = append(columnNames, stage.Name)
		columnNames = append(columnNames, stage.Name+" Dauer")
	}
	return w.csv.Write(columnNames)
}

// Write a new line in CSV format
func (w *PipeWriter) Write(d DataProvider) error {
	data := []string{
		d.ID(),
		d.Status(),
		d.DecisionTime().Local().Format("2006-01-02 15:04:05"),
		fmt.Sprintf("%v", int(d.Age().Seconds()/86400)),
		fmt.Sprintf("%0.2f", d.Value()),
	}
	for _, stage := range w.stages {
		pit := ""
		duration := ""
		for _, update := range d.DealUpdates() {
			if update.Phase == stage.Name {
				pit = update.PiT.String()
				duration = fmt.Sprintf("%d", int(update.Duration/86400))
			}
		}
		data = append(data, pit, duration)
	}
	return w.csv.Write(data)
}

// Flush the csv file
func (w *PipeWriter) Flush() error {
	w.csv.Flush()
	return nil
}