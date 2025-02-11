package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dennis-dko/go-toolkit/httphandler"

	"github.com/dennis-dko/go-toolkit/errorhandler"
	"github.com/dennis-dko/go-toolkit/example/src/service"
	"github.com/dennis-dko/go-toolkit/testhandler"
	"github.com/dennis-dko/go-toolkit/validation"

	"github.com/dennis-dko/go-toolkit/example/src/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type ExampleTestSuite struct {
	suite.Suite
	ctx                      context.Context
	mockResponse             model.Example
	mockExampleCheckResponse map[string]map[string]bool
}

func (e *ExampleTestSuite) SetupTest() {
	e.ctx = testhandler.Ctx(false, false)
	e.mockResponse = model.Example{
		Name:   "Example",
		Age:    25,
		Email:  "example@example.com",
		Active: true,
	}
	e.mockExampleCheckResponse = map[string]map[string]bool{
		"all": {"errors": false,
			"successful": true,
		},
	}
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ExampleTestSuite))
}

func (e *ExampleTestSuite) TestCreateExample() {
	tests := map[string]struct {
		createData model.Example
		status     int
		err        error
	}{
		"happy path - create example": {
			createData: e.mockResponse,
			status:     http.StatusCreated,
			err:        nil,
		},
		"should return an error while creating example": {
			createData: model.Example{},
			status:     http.StatusInternalServerError,
			err:        errorhandler.ErrDocumentNotCreate,
		},
	}
	for name, tc := range tests {
		e.Run(name, func() {
			repositoryMock := &ExampleRepositoryMock{
				InsertFunc: func(ctx context.Context, create model.Example) error {
					return tc.err
				},
			}
			exampleService := service.NewExampleService(repositoryMock)
			exampleController := NewExampleController(exampleService)
			body, _ := json.Marshal(tc.createData)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "http://localhost/example/create", bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ec := echo.New()
			ec.Validator = validation.New(e.ctx)
			c := ec.NewContext(req, rec)
			err := exampleController.CreateExample(c)
			if tc.err == nil {
				e.NoError(err)
				e.Equal(tc.status, rec.Code)
				if tc.status == rec.Code {
					var response *model.Example
					err = json.Unmarshal(rec.Body.Bytes(), &response)
					e.NoError(err)
				}
			} else {
				errorhandler.New(
					errorhandler.NewErrorStatusCodeMaps(),
				).Handler(err, c)
				e.ErrorIs(err, tc.err)
				e.Equal(tc.status, rec.Code)
			}
		})
	}
}

func (e *ExampleTestSuite) TestGetAllExamples() {
	tests := map[string]struct {
		result []model.Example
		status int
		err    error
	}{
		"happy path - find all examples": {
			result: []model.Example{
				e.mockResponse,
			},
			status: http.StatusOK,
			err:    nil,
		},
		"should return an error while getting all examples": {
			result: nil,
			status: http.StatusNotFound,
			err:    errorhandler.ErrDocumentsNotFound,
		},
	}
	for name, tc := range tests {
		e.Run(name, func() {
			repositoryMock := &ExampleRepositoryMock{
				FindAllFunc: func(ctx context.Context) ([]model.Example, error) {
					return tc.result, tc.err
				},
			}
			exampleService := service.NewExampleService(repositoryMock)
			exampleController := NewExampleController(exampleService)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost/examples", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ec := echo.New()
			ec.Validator = validation.New(e.ctx)
			c := ec.NewContext(req, rec)
			err := exampleController.GetAllExamples(c)
			if tc.err == nil {
				e.NoError(err)
				e.Equal(tc.status, rec.Code)
				if tc.status == rec.Code {
					var response []model.Example
					err = json.Unmarshal(rec.Body.Bytes(), &response)
					e.NoError(err)
					e.Equal(tc.result[0].Name, response[0].Name)
					e.Equal(tc.result[0].Age, response[0].Age)
					e.Equal(tc.result[0].Email, response[0].Email)
					e.Equal(tc.result[0].Active, response[0].Active)
				}
			} else {
				errorhandler.New(
					errorhandler.NewErrorStatusCodeMaps(),
				).Handler(err, c)
				e.ErrorIs(err, tc.err)
				e.Equal(tc.status, rec.Code)
			}
		})
	}
}

func (e *ExampleTestSuite) TestGetExampleCheckStatus() {
	tests := map[string]struct {
		filter model.Example
		result map[string]map[string]bool
		status int
		err    error
	}{
		"happy path - get example check status": {
			filter: model.Example{
				Name:   "Example",
				Age:    25,
				Email:  "example@example.com",
				Active: true,
			},
			result: e.mockExampleCheckResponse,
			status: http.StatusOK,
			err:    nil,
		},
	}
	for name, tc := range tests {
		e.Run(name, func() {
			repositoryMock := &ExampleRepositoryMock{
				ExampleCheckFunc: func(ctx context.Context, filter *model.Example) (*map[string]map[string]bool, error) {
					return &tc.result, tc.err
				},
			}
			exampleService := service.NewExampleService(repositoryMock)
			exampleController := NewExampleController(exampleService)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost/example/check", nil)
			q := req.URL.Query()
			for queryKey, queryValue := range httphandler.GetParams(tc.filter, false, httphandler.QueryTag).(map[string]string) {
				q.Add(queryKey, queryValue)
			}
			req.URL.RawQuery = q.Encode()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ec := echo.New()
			ec.Validator = validation.New(e.ctx)
			c := ec.NewContext(req, rec)
			err := exampleController.GetExampleCheckStatus(c)
			if tc.err == nil {
				e.NoError(err)
				e.Equal(tc.status, rec.Code)
			} else {
				errorhandler.New(
					errorhandler.NewErrorStatusCodeMaps(),
				).Handler(err, c)
				e.ErrorIs(err, tc.err)
			}
		})
	}
}
