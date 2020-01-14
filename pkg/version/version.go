package version

import "fmt"

const (
	Major = 1
	Minor = 0
	Patch = 6
)

func Version() string {
	return fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)
}
