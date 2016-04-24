package pipetimer

import "github.com/vitraum/golang-pipedrive"

// Pipetimer bundles all logic to generate CSV exports for a Pipedrive pipe
type Pipetimer struct {
	Stages pipedrive.Stages

	name string
	id   int
	API  *pipedrive.API
}

// Option is a type for options given to NewPipeTimer
type Option func(*Pipetimer) error

// NewPipeTimer constructs a new pipetimer object
// provide apiOptions (passed to pipedrive API and pipetimer options (see options.go)
func NewPipeTimer(apiOptions []pipedrive.Option, options ...Option) (*Pipetimer, error) {
	api, err := pipedrive.NewAPI(apiOptions...)
	if err != nil {
		return nil, err
	}
	timer := &Pipetimer{
		API: api,
	}

	for _, option := range options {
		err := option(timer)
		if err != nil {
			return nil, err
		}
	}

	plID, err := timer.API.GetPipelineIDByName(timer.name)
	if err != nil {
		return nil, err
	}
	timer.id = plID

	timer.Stages, err = timer.API.RetrieveStagesForPipeline(timer.id)
	if err != nil {
		return nil, err
	}

	return timer, nil
}
