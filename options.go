package pipetimer

import "errors"

// PipeName is an package option to provide the name of the pipe to be timed
func PipeName(name string) Option {
	return func(t *Pipetimer) error {
		if name == "" {
			return errors.New("Pipename must not be empty")
		}
		t.name = name
		return nil
	}
}
