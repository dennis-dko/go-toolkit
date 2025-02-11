package repository

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/dennis-dko/go-toolkit/database"
	"github.com/dennis-dko/go-toolkit/errorhandler"
	"github.com/dennis-dko/go-toolkit/example/src/model"
	"github.com/dennis-dko/go-toolkit/httphandler"
	"github.com/dennis-dko/go-toolkit/testhandler"
	"gorm.io/gorm"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

type ExampleTestSuite struct {
	suite.Suite
	ctx                      context.Context
	repository               *ExampleRepository
	mockExampleCheckResponse map[string]map[string]bool
}

func (e *ExampleTestSuite) SetupTest() {
	// Setup
	e.ctx = testhandler.Ctx(true, false)
	config := httphandler.Config{
		BaseURL: "http://localhost",
		Timeout: 1 * time.Minute,
	}
	e.repository = NewExampleRepository(e.ctx, &config, &gorm.DB{}, &database.MongoDBData{})
	httpmock.ActivateNonDefault(e.repository.ExampleHandler.Client.GetClient())
	e.mockExampleCheckResponse = map[string]map[string]bool{
		"all": {"errors": false,
			"successful": true,
		},
	}
}

func (e *ExampleTestSuite) SetupSubTest() {
	// Sub setup
	httpmock.Reset()
}

func (e *ExampleTestSuite) TearDownTest() {
	// Teardown
	httpmock.DeactivateAndReset()
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ExampleTestSuite))
}

func (e *ExampleTestSuite) TestExampleCheck() {
	httpmock.ActivateNonDefault(e.repository.ExampleHandler.Client.GetClient())
	tests := map[string]struct {
		data         model.Example
		responseCode int
		responseBody map[string]map[string]bool
		err          error
	}{
		"happy path - example check": {
			data: model.Example{
				Name:   "Example",
				Age:    25,
				Email:  "example@example.com",
				Active: true,
			},
			responseBody: e.mockExampleCheckResponse,
			responseCode: http.StatusOK,
			err:          nil,
		},
		"should return an error while example check": {
			data:         model.Example{},
			responseBody: map[string]map[string]bool{},
			responseCode: http.StatusInternalServerError,
			err:          errorhandler.ErrRequestFailed,
		},
		"should return an error while example check - wrong status code": {
			data:         model.Example{},
			responseBody: map[string]map[string]bool{},
			responseCode: http.StatusBadRequest,
			err:          nil,
		},
	}
	for name, tc := range tests {
		e.Run(name, func() {
			httpmock.RegisterResponder("GET", "/example/check",
				func(req *http.Request) (*http.Response, error) {
					if tc.err != nil {
						return nil, tc.err
					}
					resp, _ := httpmock.NewJsonResponse(tc.responseCode, tc.responseBody)
					return resp, nil
				})
			response, err := e.repository.ExampleCheck(
				e.ctx,
				&tc.data)
			if tc.err == nil {
				if tc.responseCode != http.StatusOK {
					e.ErrorIs(err, errorhandler.ErrRequestFailed)
				} else {
					e.NoError(err)
					e.Equal(tc.responseBody, *response)
				}
			} else {
				e.ErrorIs(err, tc.err)
			}
		})
	}
}
