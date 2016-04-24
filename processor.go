package pipetimer

import (
	"errors"
	"sort"

	"github.com/vitraum/golang-pipedrive"
)

// FetchDeals uses the pipedrive API to fetch all deals from the pipeline, using an optional filter
// filterID is either 0 or a valid filterID
func (t *Pipetimer) FetchDeals(filterName string, filterID int) (pipedrive.Deals, error) {
	if filterName != "" && filterID == 0 {
		var err error
		filterID, err = t.API.GetFilterIDByName(filterName)
		if err != nil {
			return nil, err
		}
	} else if filterName != "" && filterID != 0 {
		return nil, errors.New("Don't provide filter name and id simultainously")
	}

	res, err := t.API.FetchDealsFromPipeline(t.id, filterID)
	return res, err
}

// FilterPipelineChanges fetches the updates for the given deals and processes them
func (t *Pipetimer) FilterPipelineChanges(deals []pipedrive.Deal) (pipedrive.PipelineChangeResults, error) {
	changes, err := t.API.FetchPipelineChanges(deals, t.Stages)
	if err != nil {
		return nil, err
	}
	t.processPipelineChanges(changes)
	t.zeroShortPhases(changes, 30)

	return changes, err
}

// zeroShortPhases will add the duration of phases < cutoff to the next phase
func (t *Pipetimer) zeroShortPhases(changes pipedrive.PipelineChangeResults, cutoff float64) {
	for _, dealFlow := range changes {
		if len(dealFlow.Updates) <= 1 {
			continue
		}

		for i := range dealFlow.Updates[0 : len(dealFlow.Updates)-1] {
			next := &dealFlow.Updates[i+1]
			cur := &dealFlow.Updates[i]
			if cur.Duration < cutoff {
				next.Duration += cur.Duration
				cur.Duration = 0
				cur.PiT = next.PiT
			}
		}
	}
}

func (t *Pipetimer) processPipelineChanges(changes pipedrive.PipelineChangeResults) {
	for _, dealFlow := range changes {
		if len(dealFlow.Updates) > 0 {
			sort.Sort(dealFlow.Updates)

			if len(dealFlow.Updates) > 1 {
				for i, update := range dealFlow.Updates[1:] {
					last := dealFlow.Updates[i]
					dealFlow.Updates[i].Duration += update.PiT.Time.Sub(last.PiT.Time).Seconds()
					dealFlow.Updates[i].PhaseTouchdowns++
				}
			}

			j := len(dealFlow.Updates) - 1
			last := dealFlow.Updates[j]
			dealFlow.Updates[j].Duration += dealFlow.DecisionTime().Sub(last.PiT.Time).Seconds()
			dealFlow.Updates[j].PhaseTouchdowns++
		}
	}
}
