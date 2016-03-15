package pipetimer

import (
	"testing"
	"time"

	"github.com/vitraum/golang-pipedrive"
)

func TestZeroShortPhases(t *testing.T) {
	pt := Pipetimer{}
	then := time.Now().Round(time.Hour).Add(-2 * time.Hour)
	changes := []pipedrive.PipelineChangeResult{{
		Deal: pipedrive.Deal{},
		Updates: []pipedrive.DealFlowUpdate{
			{Phase: "A", PiT: pipedrive.NewTime(then), Duration: 11},
			{Phase: "B", PiT: pipedrive.NewTime(then.Add(1 * time.Minute)), Duration: 55},
			{Phase: "C", PiT: pipedrive.NewTime(then.Add(10 * time.Minute)), Duration: 200},
			{Phase: "D", PiT: pipedrive.NewTime(then.Add(100 * time.Minute)), Duration: 11},
			{Phase: "E", PiT: pipedrive.NewTime(then.Add(110 * time.Minute)), Duration: 11},
		}}}

	expected := []pipedrive.PipelineChangeResult{{
		Deal: pipedrive.Deal{},
		Updates: []pipedrive.DealFlowUpdate{
			{Phase: "A", PiT: pipedrive.NewTime(then.Add(1 * time.Minute)), Duration: 0},
			{Phase: "B", PiT: pipedrive.NewTime(then.Add(1 * time.Minute)), Duration: 66},
			{Phase: "C", PiT: pipedrive.NewTime(then.Add(10 * time.Minute)), Duration: 200},
			{Phase: "D", PiT: pipedrive.NewTime(then.Add(110 * time.Minute)), Duration: 0},
			{Phase: "E", PiT: pipedrive.NewTime(then.Add(110 * time.Minute)), Duration: 22},
		}}}

	pt.zeroShortPhases(changes, 12)
	for i, wanted := range expected[0].Updates {
		got := changes[0].Updates[i]
		if wanted.Duration != got.Duration || wanted.PiT != got.PiT {
			t.Errorf("wanted: %v, got: %v", wanted, got)
		}
	}
}

func TestZeroShortPhasesEmpty(t *testing.T) {
	pt := Pipetimer{}
	then := time.Now().Round(time.Hour).Add(-2 * time.Hour)
	changes := []pipedrive.PipelineChangeResult{{
		Deal: pipedrive.Deal{},
		Updates: []pipedrive.DealFlowUpdate{
			{Phase: "A", PiT: pipedrive.NewTime(then), Duration: 11},
		}}}

	expected := []pipedrive.PipelineChangeResult{{
		Deal: pipedrive.Deal{},
		Updates: []pipedrive.DealFlowUpdate{
			{Phase: "A", PiT: pipedrive.NewTime(then), Duration: 11},
		}}}

	pt.zeroShortPhases(changes, 12)
	for i, wanted := range expected[0].Updates {
		got := changes[0].Updates[i]
		if wanted.Duration != got.Duration || wanted.PiT != got.PiT {
			t.Errorf("wanted: %v, got: %v", wanted, got)
		}
	}
}

func TestProcessPipelineChanges(t *testing.T) {
	pt := Pipetimer{}
	then := time.Now().Round(time.Hour).Add(-2 * time.Hour)
	lostAt := pipedrive.NewTime(then.Add(5 * time.Minute))
	bAt := then.Add(1 * time.Minute)
	changes := []pipedrive.PipelineChangeResult{{
		Deal: pipedrive.Deal{
			Status: "lost",
			LostAt: &lostAt,
		},
		Updates: []pipedrive.DealFlowUpdate{
			{Phase: "A", PiT: pipedrive.NewTime(then)},
			{Phase: "B", PiT: pipedrive.NewTime(then.Add(1 * time.Minute))},
		}}}

	expected := []pipedrive.PipelineChangeResult{{
		Deal: pipedrive.Deal{},
		Updates: []pipedrive.DealFlowUpdate{
			{Phase: "A", PiT: pipedrive.NewTime(then), Duration: 60, PhaseTouchdowns: 1},
			{Phase: "B", PiT: pipedrive.NewTime(bAt), Duration: lostAt.Sub(bAt).Seconds(), PhaseTouchdowns: 1},
		}}}

	pt.processPipelineChanges(changes)
	for i, wanted := range expected[0].Updates {
		got := changes[0].Updates[i]
		if wanted.Duration != got.Duration || wanted.PiT != got.PiT {
			t.Errorf("wanted: %v, got: %v", wanted, got)
		}
	}
}
