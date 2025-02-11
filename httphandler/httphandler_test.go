package httphandler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/dennis-dko/go-toolkit/datatype"
	"github.com/dennis-dko/go-toolkit/testhandler"

	"github.com/jarcoal/httpmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type StructTest struct {
	FirstName  string                `query:"firstname" json:"first_name"`
	LastName   string                `query:"lastname" json:"last_name"`
	Birthday   *datatype.CustomDate  `query:"birthday" json:"birthday"`
	Hobbies    []string              `query:"hobbies" json:"hobbies"`
	Animals    []datatype.NullString `query:"animals" json:"animals"`
	Answers    []bool                `query:"answers" json:"answers"`
	Priorities []int64               `query:"priorities" json:"priorities"`
	Coins      *[]float64            `query:"coins" json:"coins"`
	Age        uint8                 `param:"age" json:"age"`
	Email      string                `param:"email" json:"email"`
}

type HttpHandlerTestSuite struct {
	suite.Suite
	ctx                       context.Context
	httpHandler               *HttpHandler
	getRequest                *HttpRequest
	postRequest               *HttpRequest
	putRequest                *HttpRequest
	deleteRequest             *HttpRequest
	pathParams                StructTest
	queryParams               StructTest
	defaultParams             StructTest
	expectedPathParams        map[string]string
	expectedStringQueryParams map[string]string
	expectedSliceQueryParams  map[string][]string
	expectedDefaultParams     map[string]string
	jsonResponse              string
}

func (h *HttpHandlerTestSuite) SetupTest() {
	// Setup
	h.ctx = testhandler.Ctx(true, false)
	config := Config{
		BaseURL: "http://localhost",
		Timeout: 1 * time.Minute,
	}
	h.httpHandler = New(h.ctx, &config)
	httpmock.ActivateNonDefault(h.httpHandler.Client.GetClient())
}

func (h *HttpHandlerTestSuite) SetupSubTest() {
	// Sub setup
	currentDate, _ := datatype.NewDate(false)
	h.jsonResponse = `
		{
			"first_name": "Walter", 
			"last_name": "White", 
			"age": 50,
			"email": "walther.white@example.com"
		}
	`
	h.getRequest = &HttpRequest{
		Method:           http.MethodGet,
		URL:              "/users/{userId}/get",
		ForceContentType: echo.MIMEApplicationJSON,
		PathParams: map[string]string{
			"userId": "1",
		},
		DestResult: StructTest{},
	}
	h.postRequest = &HttpRequest{
		Method:     http.MethodPost,
		URL:        "/users/post",
		Body:       []byte(h.jsonResponse),
		DestResult: StructTest{},
	}
	h.putRequest = &HttpRequest{
		Method:     http.MethodPut,
		URL:        "/users/put",
		Body:       []byte(h.jsonResponse),
		DestResult: StructTest{},
	}
	h.deleteRequest = &HttpRequest{
		Method: http.MethodDelete,
		URL:    "/users/delete",
		QueryParams: map[string]string{
			"userId": "1",
		},
	}
	h.pathParams = StructTest{
		Age:   50,
		Email: "walther.white@example.com",
	}
	h.expectedPathParams = map[string]string{
		"age":   "50",
		"email": "walther.white@example.com",
	}
	h.queryParams = StructTest{
		FirstName: "Walter",
		LastName:  "White",
		Birthday:  currentDate,
		Hobbies:   []string{"cooking", "chemistry"},
		Animals: []datatype.NullString{
			datatype.NewNullString(datatype.StringPtr("dog")),
			datatype.NewNullString(datatype.StringPtr("cat")),
		},
		Answers:    []bool{true, false},
		Priorities: []int64{1, 2, 3, 4, 5},
		Coins:      &[]float64{5.6, 3.2},
	}
	h.expectedStringQueryParams = map[string]string{
		"firstname": "Walter",
		"lastname":  "White",
		"birthday":  currentDate.String(),
	}

	h.expectedSliceQueryParams = map[string][]string{
		"hobbies": {
			"cooking",
			"chemistry",
		},
		"animals": {
			"dog",
			"cat",
		},
		"answers": {
			"true",
			"false",
		},
		"priorities": {
			"1",
			"2",
			"3",
			"4",
			"5",
		},
		"coins": {
			"5.6",
			"3.2",
		},
	}

	h.defaultParams = StructTest{
		FirstName: "Walter",
		LastName:  "White",
		Age:       50,
		Email:     "walther.white@example.com",
	}
	h.expectedDefaultParams = map[string]string{
		"firstname": "Walter",
		"lastname":  "White",
		"age":       "50",
		"email":     "walther.white@example.com",
	}
	httpmock.Reset()
}

func (h *HttpHandlerTestSuite) TearDownTest() {
	// Teardown
	httpmock.DeactivateAndReset()
}

func TestHttpHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HttpHandlerTestSuite))
}

func (h *HttpHandlerTestSuite) TestGetRequest() {

	h.Run("happy path - get request", func() {
		// Init
		data := new(bytes.Buffer)
		expectedData := new(bytes.Buffer)
		responder := httpmock.NewStringResponder(http.StatusOK, h.jsonResponse)
		httpmock.RegisterResponder("GET", fmt.Sprintf("/users/%s/get", h.getRequest.PathParams["userId"]), responder)
		_ = json.Compact(expectedData, []byte(h.jsonResponse))

		// Run
		response, err := h.httpHandler.DoHTTPRequest(h.getRequest)
		_ = json.Compact(data, []byte(response.String()))
		structResult := response.Result().(*StructTest)
		_, requestIDExists := response.Request.Header[echo.HeaderXRequestID]

		// Assert
		h.NoError(err)
		h.NotNil(structResult)
		h.True(requestIDExists)
		h.Equal(http.StatusOK, response.StatusCode())
		h.Equal(expectedData.String(), data.String())
	})

	h.Run("should return an error while no responder exist", func() {
		// Run
		response, err := h.httpHandler.DoHTTPRequest(h.getRequest)

		// Assert
		h.Error(err)
		h.ErrorContains(err, "no responder found")
		h.Nil(response)
	})
}

func (h *HttpHandlerTestSuite) TestPostRequest() {

	h.Run("happy path - post request", func() {
		// Init
		data := new(bytes.Buffer)
		expectedData := new(bytes.Buffer)
		responder := httpmock.NewStringResponder(http.StatusCreated, h.jsonResponse)
		httpmock.RegisterResponder("POST", "/users/post", responder)
		_ = json.Compact(expectedData, []byte(h.jsonResponse))

		// Run
		response, err := h.httpHandler.DoHTTPRequest(h.postRequest)
		_ = json.Compact(data, []byte(response.String()))
		structResult := response.Result().(*StructTest)
		_, requestIDExists := response.Request.Header[echo.HeaderXRequestID]

		// Assert
		h.NoError(err)
		h.NotNil(structResult)
		h.True(requestIDExists)
		h.Equal(http.StatusCreated, response.StatusCode())
		h.Equal(expectedData.String(), data.String())
	})

	h.Run("should return an error while no responder exist", func() {
		// Run
		response, err := h.httpHandler.DoHTTPRequest(h.postRequest)

		// Assert
		h.Error(err)
		h.ErrorContains(err, "no responder found")
		h.Nil(response)
	})
}

func (h *HttpHandlerTestSuite) TestPutRequest() {

	h.Run("happy path - put request", func() {
		// Init
		data := new(bytes.Buffer)
		expectedData := new(bytes.Buffer)
		responder := httpmock.NewStringResponder(http.StatusOK, h.jsonResponse)
		httpmock.RegisterResponder("PUT", "/users/put", responder)
		_ = json.Compact(expectedData, []byte(h.jsonResponse))

		// Run
		response, err := h.httpHandler.DoHTTPRequest(h.putRequest)
		_ = json.Compact(data, []byte(response.String()))
		structResult := response.Result().(*StructTest)
		_, requestIDExists := response.Request.Header[echo.HeaderXRequestID]

		// Assert
		h.NoError(err)
		h.NotNil(structResult)
		h.True(requestIDExists)
		h.Equal(http.StatusOK, response.StatusCode())
		h.Equal(expectedData.String(), data.String())
	})

	h.Run("should return an error while no responder exist", func() {
		// Run
		response, err := h.httpHandler.DoHTTPRequest(h.putRequest)

		// Assert
		h.Error(err)
		h.ErrorContains(err, "no responder found")
		h.Nil(response)
	})
}

func (h *HttpHandlerTestSuite) TestDeleteRequest() {

	h.Run("happy path - delete request", func() {
		// Init
		responder := httpmock.NewStringResponder(http.StatusNoContent, "")
		httpmock.RegisterResponder("DELETE", "/users/delete", responder)

		// Run
		response, err := h.httpHandler.DoHTTPRequest(h.deleteRequest)
		_, requestIDExists := response.Request.Header[echo.HeaderXRequestID]

		// Assert
		h.NoError(err)
		h.True(requestIDExists)
		h.Equal(http.StatusNoContent, response.StatusCode())
		h.Equal("", response.String())
	})

	h.Run("should return an error while no responder exist", func() {
		// Run
		response, err := h.httpHandler.DoHTTPRequest(h.deleteRequest)

		// Assert
		h.Error(err)
		h.ErrorContains(err, "no responder found")
		h.Nil(response)
	})
}

func (h *HttpHandlerTestSuite) TestBodyClose() {

	h.Run("happy path - body close", func() {
		// Init
		response, respErr := http.Get("http://example.com")

		// Run
		Close(h.ctx, response)
		body, bodyErr := response.Body.Read(nil)

		// Assert
		h.NoError(respErr)
		h.Equal(0, body)
		h.ErrorContains(bodyErr, "closed response body")
	})
}

func (h *HttpHandlerTestSuite) TestGetParams() {

	h.Run("happy path - get all path parameters", func() {
		// Run
		pathParams := GetParams(h.pathParams, false, PathTag).(map[string]string)

		// Assert
		h.Equal(h.expectedPathParams, pathParams)
	})

	h.Run("happy path - get all query parameters", func() {
		// Run
		stringQueryParams := GetParams(h.queryParams, false, QueryTag).(map[string]string)
		sliceQueryParams := GetParams(h.queryParams, true, QueryTag).(map[string][]string)

		// Assert
		h.Equal(h.expectedStringQueryParams, stringQueryParams)
		h.Equal(h.expectedSliceQueryParams, sliceQueryParams)
	})

	h.Run("happy path - get all default parameters", func() {
		// Run
		defaultParams := GetParams(h.defaultParams, false).(map[string]string)
		emptyParams := GetParams(nil, false).(map[string]string)

		// Assert
		h.Equal(h.expectedDefaultParams, defaultParams)
		h.Empty(emptyParams)
	})
}

func (h *HttpHandlerTestSuite) TestGetHeaderCtxValue() {

	h.Run("happy path - get header value from context", func() {
		// Run
		value := GetHeaderCtxValue(h.ctx, echo.HeaderXRequestID)

		// Assert
		h.NotEmpty(value)
	})

	h.Run("happy path - get empty header value from context", func() {
		// Run
		emptyValue := GetHeaderCtxValue(h.ctx, "test")

		// Assert
		h.Empty(emptyValue)
	})
}
