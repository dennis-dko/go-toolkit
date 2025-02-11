package service

import (
	"context"
	"testing"

	"github.com/dennis-dko/go-toolkit/errorhandler"

	"github.com/dennis-dko/go-toolkit/testhandler"

	"github.com/dennis-dko/go-toolkit/example/src/model"
	"github.com/stretchr/testify/suite"
)

type ExampleTestSuite struct {
	suite.Suite
	ctx                      context.Context
	mockListResponse         model.Example
	mockExampleCheckResponse map[string]map[string]bool
}

func (e *ExampleTestSuite) SetupTest() {
	e.ctx = testhandler.Ctx(false, false)
	e.mockListResponse = model.Example{
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
		create model.Example
		err    error
	}{
		"happy path - create example": {
			create: e.mockListResponse,
			err:    nil,
		},
		"should return an error while creating an example": {
			create: model.Example{},
			err:    errorhandler.ErrDocumentNotCreate,
		},
	}
	for name, tc := range tests {
		e.Run(name, func() {
			repositoryMock := &ExampleRepositoryMock{
				InsertFunc: func(ctx context.Context, create model.Example) error {
					return tc.err
				},
			}
			service := NewExampleService(repositoryMock)
			err := service.CreateExample(e.ctx, tc.create)
			if tc.err == nil {
				e.NoError(err)
			} else {
				e.ErrorIs(err, tc.err)
			}
		})
	}
}

func (e *ExampleTestSuite) TestFindExamples() {
	tests := map[string]struct {
		result []model.Example
		err    error
	}{
		"happy path - return all examples": {
			result: []model.Example{
				e.mockListResponse,
			},
			err: nil,
		},
		"should return an error while getting all examples": {
			result: nil,
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
			service := NewExampleService(repositoryMock)
			examples, err := service.FindExamples(e.ctx)
			if tc.err == nil {
				e.NoError(err)
				e.Equal(tc.result, examples)
			} else {
				e.ErrorIs(err, tc.err)
				e.Equal(tc.result, examples)
			}
		})
	}
}

func (e *ExampleTestSuite) TestCheckExample() {
	tests := map[string]struct {
		filter         model.Example
		expectedResult map[string]map[string]bool
		err            error
	}{
		"happy path - checking example successfully": {
			filter: model.Example{
				Name:   "Example",
				Age:    25,
				Email:  "example@example.com",
				Active: true,
			},
			expectedResult: e.mockExampleCheckResponse,
			err:            nil,
		},
		"should return an error while checking example": {
			filter:         model.Example{},
			expectedResult: nil,
			err:            errorhandler.ErrRequestFailed,
		},
	}
	for name, tc := range tests {
		e.Run(name, func() {
			repositoryMock := &ExampleRepositoryMock{
				ExampleCheckFunc: func(ctx context.Context, filter *model.Example) (*map[string]map[string]bool, error) {
					return &tc.expectedResult, tc.err
				},
			}
			service := NewExampleService(repositoryMock)
			data, err := service.CheckExample(e.ctx, &tc.filter)
			if tc.err == nil {
				e.NoError(err)
				e.Equal(tc.expectedResult, data)
			} else {
				e.ErrorIs(err, tc.err)
			}
		})
	}
}
