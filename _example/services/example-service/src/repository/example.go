package repository

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dennis-dko/go-toolkit/database"
	"github.com/dennis-dko/go-toolkit/errorhandler"
	"github.com/dennis-dko/go-toolkit/example/src/model"
	"github.com/dennis-dko/go-toolkit/httphandler"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

const (
	ExampleTable = "example"
)

type ExampleRepository struct {
	ExampleHandler *httphandler.HttpHandler
	postgresClient *gorm.DB
	mongoDb        *database.MongoDBData
}

func NewExampleRepository(ctx context.Context, exampleCfg *httphandler.Config, postgresClient *gorm.DB, mongoDb *database.MongoDBData) *ExampleRepository {
	return &ExampleRepository{
		mongoDb:        mongoDb,
		postgresClient: postgresClient,
		ExampleHandler: httphandler.New(ctx, exampleCfg),
	}
}

// Insert inserts a new example object
func (e *ExampleRepository) Insert(ctx context.Context, create model.Example) error {
	_, err := e.mongoDb.Collections[fmt.Sprintf("%s-objects", ExampleTable)].InsertOne(ctx, create)
	if err != nil {
		return err
	}
	return nil
}

// Find returns all examples
func (e *ExampleRepository) FindAll(ctx context.Context) ([]model.Example, error) {
	var example []model.Example
	err := e.postgresClient.WithContext(ctx).Table(ExampleTable).Find(&example).Error
	if err != nil {
		return nil, err
	}
	return example, nil
}

// ExampleCheck returns the example status
func (e *ExampleRepository) ExampleCheck(ctx context.Context, filter *model.Example) (*map[string]map[string]bool, error) {
	// Send the request to another service
	req := httphandler.HttpRequest{
		Method:           http.MethodGet,
		URL:              "/example/check",
		ForceContentType: echo.MIMEApplicationJSON,
		QueryParams:      httphandler.GetParams(filter, false, httphandler.QueryTag).(map[string]string),
		Headers: map[string]string{
			echo.HeaderXRequestID: httphandler.GetHeaderCtxValue(ctx, echo.HeaderXRequestID),
		},
		DestResult: map[string]map[string]bool{},
	}
	response, err := e.ExampleHandler.DoHTTPRequest(&req)
	if err != nil {
		slog.ErrorContext(ctx, "error while executing example status request", slog.String("error", err.Error()))
		return nil, errorhandler.ErrRequestFailed
	}
	if response.StatusCode() != http.StatusOK {
		slog.ErrorContext(ctx, "error while executing example status request - unexpected status code", slog.Int("status code", response.StatusCode()))
		return nil, errorhandler.ErrRequestFailed
	}
	return response.Result().(*map[string]map[string]bool), nil
}
