// Package version contains version information for this app
package version

import (
	"fmt"
	"runtime"
)

var Version = "was not built properly"

func String() string {
	return fmt.Sprintf(`Version: %s, Go version: %s, Go OS/ARCH: %s %s`,
		Version,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH)
}
