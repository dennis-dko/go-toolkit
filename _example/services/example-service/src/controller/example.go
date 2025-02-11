package controller

import (
	"fmt"
	"net/http"

	"github.com/dennis-dko/go-toolkit/errorhandler"
	"github.com/dennis-dko/go-toolkit/example/src/service"

	"github.com/dennis-dko/go-toolkit/example/src/model"
	"github.com/labstack/echo/v4"
)

//go:generate moq -pkg controller -out example_mock_test.go ../service ExampleRepository

type ExampleController struct {
	ExampleService *service.ExampleService
}

func NewExampleController(service *service.ExampleService) *ExampleController {
	return &ExampleController{
		ExampleService: service,
	}
}

// CreateExample godoc
//
//	@Summary		Create an example
//	@Description	Create an entry of example
//	@ID				example-create
//	@Tags			Example Actions
//	@Accept			json
//	@Produce		json
//	@Param			params	body model.Example	true	"Example Data"
//	@Success		204	{string} http.StatusCreated
//	@Failure		400	{string} http.StatusBadRequest
//	@Failure		500	{string} http.StatusInternalServerError
//	@Router			/create [post]
func (e *ExampleController) CreateExample(c echo.Context) error {
	// Bind data
	data := model.Example{}
	if err := c.Bind(&data); err != nil {
		return fmt.Errorf("%s (%w)", err.Error(), errorhandler.ErrBindingFailed)
	}

	// Validate data
	if err := c.Validate(data); err != nil {
		return fmt.Errorf("%s (%w)", err.Error(), errorhandler.ErrValidationFailed)
	}

	// Create an example
	err := e.ExampleService.CreateExample(c.Request().Context(), data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, data)
}

// GetAllExamples godoc
//
//	@Summary		Get all examples
//	@Description	Get all entries of examples
//	@ID				examples-get
//	@Tags			Example Actions
//	@Produce		json
//	@Success		200	{array}     model.Example
//	@Failure		400	{string}	http.StatusBadRequest
//	@Failure		500	{string}	http.StatusInternalServerError
//	@Router			/examples [get]
func (e *ExampleController) GetAllExamples(c echo.Context) error {
	// Find all examples
	examples, err := e.ExampleService.FindExamples(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, examples)
}

// GetExampleCheckStatus godoc
//
//	@Summary		Get example check status
//	@Description	Get example check status
//	@ID				example-check-status-get
//	@Tags			Example Actions
//	@Produce		json
//	@Param			params	query	model.Example	true	"Example Check Status Filter Data"
//	@Success		200	{string}    http.StatusOK
//	@Failure		400	{string}	http.StatusBadRequest
//	@Failure		500	{string}	http.StatusInternalServerError
//	@Router			/example/check [get]
func (e *ExampleController) GetExampleCheckStatus(c echo.Context) error {
	// Bind filter
	filterRequest := model.Example{}
	if err := c.Bind(&filterRequest); err != nil {
		return fmt.Errorf("%s (%w)", err.Error(), errorhandler.ErrBindingFailed)
	}

	// Validate filter
	if err := c.Validate(filterRequest); err != nil {
		return fmt.Errorf("%s (%w)", err.Error(), errorhandler.ErrValidationFailed)
	}

	// Get example check status
	data, err := e.ExampleService.CheckExample(c.Request().Context(), &filterRequest)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}
