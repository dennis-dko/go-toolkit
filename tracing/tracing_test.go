package tracing

import (
	"context"
	"testing"
	"time"

	"github.com/dennis-dko/go-toolkit/testhandler"

	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

type TracingTestSuite struct {
	suite.Suite
	ctx            context.Context
	ctxCanceled    context.Context
	serviceName    string
	configEnabled  Config
	configDisabled Config
	configInvalid  Config
}

func (t *TracingTestSuite) SetupTest() {
	// Setup
	t.ctx = testhandler.Ctx(false, false)
	t.ctxCanceled = testhandler.Ctx(false, true)
	t.serviceName = "test-service"
	t.configEnabled = Config{
		Enabled:            true,
		Host:               "example.com",
		Port:               "4317",
		BatchTimeout:       5000 * time.Millisecond,
		MaxExportBatchSize: 512,
		HttpInsecure:       true,
	}
	t.configDisabled = Config{
		Enabled: false,
	}
	t.configInvalid = Config{
		Enabled:            t.configEnabled.Enabled,
		Host:               "",
		Port:               "",
		BatchTimeout:       t.configEnabled.BatchTimeout,
		MaxExportBatchSize: t.configEnabled.MaxExportBatchSize,
		HttpInsecure:       t.configEnabled.HttpInsecure,
	}
}

func TestTracingTestSuite(t *testing.T) {
	suite.Run(t, new(TracingTestSuite))
}

func (t *TracingTestSuite) TestTracingEnabledWithValidConfig() {
	err := t.configEnabled.Provide(t.ctx, t.serviceName)
	t.NoError(err)
	tp, ok := otel.GetTracerProvider().(*trace.TracerProvider)
	t.True(ok)
	t.NotEmpty(tp)
}

func (t *TracingTestSuite) TestTracingDisabled() {
	err := t.configDisabled.Provide(t.ctx, t.serviceName)
	t.NoError(err)
	t.Empty(otel.GetTracerProvider())
}

func (t *TracingTestSuite) TestTracingEnabledWithInvalidEndpoint() {
	err := t.configInvalid.Provide(t.ctxCanceled, t.serviceName)
	t.Error(err)
	t.Empty(otel.GetTracerProvider())
}
