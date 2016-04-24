package pipetimer

import (
	"strconv"
	"sync"
	"time"

	"github.com/vitraum/golang-pipedrive"
)

type sourceMap map[string]string

type sourceCache struct {
	smap   sourceMap
	filled bool
	sync.Mutex
}

var cache sourceCache

type ChangeResultConverter struct {
	pipedrive.PipelineChangeResult
	sourceMapping sourceMap
}

func NewChangeResultConverter(cr pipedrive.PipelineChangeResult, api *pipedrive.API) *ChangeResultConverter {
	return &ChangeResultConverter{
		cr,
		fetchSourceMapping(api),
	}
}

func (cr *ChangeResultConverter) ID() string {
	return strconv.Itoa(cr.Deal.Id)
}

func (cr *ChangeResultConverter) Status() string {
	return cr.Deal.Status
}

func (cr *ChangeResultConverter) Value() float64 {
	return cr.Deal.Value
}

func (cr *ChangeResultConverter) Source() string {
	return cr.Deal.Source
}

func (cr *ChangeResultConverter) Age() time.Duration {
	return cr.DecisionTime().Sub(cr.Deal.Added.Time)
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
