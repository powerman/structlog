package structlog_test

import (
	"context"
	"testing"

	"github.com/powerman/check"

	"github.com/powerman/structlog"
)

func TestContext(tt *testing.T) {
	t := check.T(tt)
	log1 := structlog.New()
	log2 := structlog.New()
	log3 := structlog.New()
	ctx := context.Background()
	t.NotNil(structlog.FromContext(ctx, nil))
	t.HasType(structlog.FromContext(ctx, nil), log1)
	t.Equal(structlog.FromContext(ctx, log1), log1)
	ctx = structlog.NewContext(ctx, log2)
	t.Equal(structlog.FromContext(ctx, nil), log2)
	t.Equal(structlog.FromContext(ctx, log1), log2)
	ctx = structlog.NewContext(ctx, log3)
	t.Equal(structlog.FromContext(ctx, nil), log3)
	t.Equal(structlog.FromContext(ctx, log1), log3)
}
