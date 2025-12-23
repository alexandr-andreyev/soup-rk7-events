package app

import (
	"github.com/pkg/errors"
)

// Run launches the service
func Run(svcName, sha1ver string) error {

	s, err := setup(svcName, sha1ver)
	if err != nil {
		return errors.Wrap(err, "setup")
	}

	// Your service should be launched as a GO routine
	go runApp(s)

	return nil
}
