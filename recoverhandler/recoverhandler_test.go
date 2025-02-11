package recoverhandler

import (
	"context"
	"testing"

	"github.com/dennis-dko/go-toolkit/testhandler"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type PanicHandlerTestSuite struct {
	suite.Suite
	ctx      context.Context
	instance *echo.Echo
	config   Config
}

func (p *PanicHandlerTestSuite) SetupTest() {
	// Setup
	p.instance = echo.New()
	p.config = Config{
		StackSize:         4096,
		DisableStackAll:   false,
		DisablePrintStack: false,
	}
}

func (p *PanicHandlerTestSuite) SetupSubTest() {
	// Sub setup
	p.ctx = testhandler.Ctx(false, false)
}

func TestPanicHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(PanicHandlerTestSuite))
}

func (p *PanicHandlerTestSuite) TestUseRecover() {

	p.Run("happy path - use recover", func() {
		// Run
		p.config.Provide()
		UseRecover(p.ctx, p.instance)

		// Assert
		p.Equal(p.config, *config)
	})
}
