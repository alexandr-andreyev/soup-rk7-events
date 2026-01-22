package app

import (
	"fmt"

	"github.com/pkg/errors"
)

// BuildInfo contains version information
type BuildInfo struct {
	Version   string
	Commit    string
	BuildTime string
}

// String returns formatted version string
func (b BuildInfo) String() string {
	return fmt.Sprintf("%s (commit: %s, built: %s)", b.Version, b.Commit, b.BuildTime)
}

// Run launches the service
func Run(svcName string, build BuildInfo) error {

	s, err := setup(svcName, build)
	if err != nil {
		return errors.Wrap(err, "setup")
	}

	// Your service should be launched as a GO routine
	go runApp(s)

	return nil
}
