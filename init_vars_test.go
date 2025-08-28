package structlog_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

var (
	pid            = strconv.Itoa(os.Getpid())
	wd, _          = os.Getwd()
	unit           = filepath.Base(wd) // Differ on CI.
	osNotExistsMsg = func() string {
		if runtime.GOOS == "windows" {
			return "The system cannot find the path specified."
		}
		return "no such file or directory"
	}()
)
