package router

import (
	"github.com/dennis-dko/go-toolkit/acl"
	"github.com/dennis-dko/go-toolkit/database"
	"github.com/dennis-dko/go-toolkit/errorhandler"
	"github.com/dennis-dko/go-toolkit/example/src/config"
	"github.com/dennis-dko/go-toolkit/example/src/constant"
	"github.com/dennis-dko/go-toolkit/example/src/controller"
	"github.com/dennis-dko/go-toolkit/example/src/repository"
	"github.com/dennis-dko/go-toolkit/example/src/service"
	"github.com/dennis-dko/go-toolkit/httphandler"
	"github.com/dennis-dko/go-toolkit/logging"
	"github.com/dennis-dko/go-toolkit/recoverhandler"
	"github.com/dennis-dko/go-toolkit/secure"
	s "github.com/dennis-dko/go-toolkit/server"
	"github.com/dennis-dko/go-toolkit/server/health"
	"github.com/dennis-dko/go-toolkit/validation"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

func Init(server *s.Server, cfg *config.Config, postgres *gorm.DB, mongoDb *database.MongoDBData) {
	// Initialize error handler
	server.Echo.HTTPErrorHandler = errorhandler.New(
		errorhandler.NewErrorStatusCodeMaps(),
	).Handler

	// Initialize validator
	server.Echo.Validator = validation.New(server.Context)

	// Initialize middleware
	logging.UseRequestLog(server.Context, server.Echo)
	logging.UseBodyDump(server.Context, server.Echo)
	recoverhandler.UseRecover(server.Context, server.Echo)
	httphandler.UseRequestID(server.Context, server.Echo)
	secure.UseSecure(server.Context, server.Echo)
	acl.UseAuthEnforcer(server.Context, server.Echo)

	// Initialize example dependencies
	exampleRepository := repository.NewExampleRepository(server.Context, &cfg.Client.ExampleService, postgres, mongoDb)
	exampleService := service.NewExampleService(exampleRepository)
	exampleController := controller.NewExampleController(exampleService)

	// Initialize health dependencies
	healthController := health.NewHealthController()

	// Initialize example
	ex := server.Echo.Group("/example")
	ex.GET("/check", exampleController.GetExampleCheckStatus)
	ex.POST("/create", exampleController.CreateExample)

	// Initialize examples
	exls := server.Echo.Group("/examples")
	exls.GET("", exampleController.GetAllExamples)

	// Initialize swagger
	server.Echo.GET(constant.SwaggerRoute, echoSwagger.WrapHandler)

	// Initialize health
	server.Echo.GET(constant.HealthRoute, healthController.HandleHealth)
}
