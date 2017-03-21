package pipetimer

import (
	"strconv"
	"sync"
	"time"

	"github.com/vitraum/golang-pipedrive"
)

type sourceMap map[string]string

type stageMap map[int]string

type sourceCache struct {
	smap   sourceMap
	filled bool
	sync.Mutex
}

type AgeCalculator func(a, b time.Time) time.Duration

var cache sourceCache

type ChangeResultConverter struct {
	pipedrive.PipelineChangeResult
	sourceMapping sourceMap
	ageCalculator AgeCalculator
	stages        stageMap
}

func NewChangeResultConverter(cr pipedrive.PipelineChangeResult, api *pipedrive.API, age AgeCalculator, stages pipedrive.Stages) *ChangeResultConverter {
	stageMapping := stageMap{}
	for _, stage := range stages {
		stageMapping[stage.Id] = stage.Name
	}

	return &ChangeResultConverter{
		cr,
		fetchSourceMapping(api),
		age,
		stageMapping,
	}
}

func (cr *ChangeResultConverter) ID() string {
	return strconv.Itoa(cr.Deal.Id)
}

func (cr *ChangeResultConverter) Status() string {
	return cr.Deal.Status
}

func (cr *ChangeResultConverter) Added() time.Time {
	return cr.Deal.Added.Time
}

func (cr *ChangeResultConverter) Value() float64 {
	return cr.Deal.Value
}

func (cr *ChangeResultConverter) Source() string {
	return cr.Deal.Source
}

func (cr *ChangeResultConverter) LastStage() string {
	stageName, ok := cr.stages[cr.Deal.Stage]
	if ok {
		return stageName
	} else {
		return ""
	}
}

func (cr *ChangeResultConverter) Age() time.Duration {
	return cr.ageCalculator(cr.DecisionTime(), cr.Deal.Added.Time)
}

func (cr *ChangeResultConverter) DealUpdates() pipedrive.DealFlowUpdates {
	return cr.Updates
}

func fetchSourceMapping(api *pipedrive.API) sourceMap {
	cache.Lock()
	defer cache.Unlock()
	if cache.filled {
		return cache.smap
	}

	result := sourceMap{}

	return result
}
