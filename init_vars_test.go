package structlog_test

import (
	"os"
	"path"
	"strconv"
)

var (
	pid   = strconv.Itoa(os.Getpid())
	wd, _ = os.Getwd()
	unit  = path.Base(wd) // Differ on CI.
)
