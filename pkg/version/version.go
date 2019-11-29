package version

import (
	"fmt"
	"strings"
)

var (
	version   = "v1.0.0"
	branch    = "unkown"
	buildTime = ""
	commitID  = ""
)

func GetVersion() string {
	return version
}

func FullVersion() string {
	v := []string{
		"Version", version,
		"Branch", branch,
		"CommitID", commitID,
		"BuildTime", buildTime,
	}

	maxKeySize := 0
	for ix := 0; ix < len(v); ix += 2 {
		if len(v[ix]) > maxKeySize {
			maxKeySize = len(v[ix])
		}
	}

	msg := ""
	for ix, val := range v {
		if ix%2 == 0 {
			msg += fmt.Sprintf("%s%s: ", val, strings.Repeat(" ", maxKeySize-len(val)))
		} else {
			msg += fmt.Sprint(val)
			if ix != len(v)-1 {
				msg += "\n"
			}
		}
	}
	return msg
}
