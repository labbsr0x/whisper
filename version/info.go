package version

import "fmt"

// Default build-time variable.
// These values are overridden via ldflags
var (
	Version   = "unknown-version"
	GitCommit = "unknown-commit"
	BuildTime = "unknown-buildtime"
)

const versionF = `Whisper
  Version: %s
  GitCommit: %s
  BuildTime: %s
`

// FormattedMessage gets the full formatted version message
func FormattedMessage() string {
	return fmt.Sprintf(versionF, Version, GitCommit, BuildTime)
}
