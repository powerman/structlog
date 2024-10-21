//nolint:testpackage // To set global "now" variable.
package structlog

import (
	"log"
	"testing"
	"time"

	"github.com/powerman/check"
)

func TestMain(m *testing.M) {
	loc, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		log.Fatal(err)
	}
	now = func() time.Time { return time.Date(2020, 1, 2, 3, 4, 5, 123456789, loc) }

	check.TestMain(m)
}
