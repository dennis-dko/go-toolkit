package service

import (
	"context"

	"github.com/dennis-dko/go-toolkit/example/src/model"
)

//go:generate moq -out example_mock_test.go . ExampleRepository

type ExampleRepository interface {
	Insert(ctx context.Context, create model.Example) error
	FindAll(ctx context.Context) ([]model.Example, error)
	ExampleCheck(ctx context.Context, filter *model.Example) (*map[string]map[string]bool, error)
}

type ExampleService struct {
	ExampleRepository ExampleRepository
}

func NewExampleService(exampleRepository ExampleRepository) *ExampleService {
	return &ExampleService{
		ExampleRepository: exampleRepository,
	}
}

func (e *ExampleService) CreateExample(ctx context.Context, create model.Example) error {
	err := e.ExampleRepository.Insert(ctx, create)
	if err != nil {
		return err
	}
	return nil
}

func (e *ExampleService) FindExamples(ctx context.Context) ([]model.Example, error) {
	data, err := e.ExampleRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (e *ExampleService) CheckExample(ctx context.Context, filter *model.Example) (map[string]map[string]bool, error) {
	data, err := e.ExampleRepository.ExampleCheck(ctx, filter)
	if err != nil {
		return nil, err
	}
	return *data, nil
}
