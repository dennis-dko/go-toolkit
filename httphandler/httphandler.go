package httphandler

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/dennis-dko/go-toolkit/datatype"
	"github.com/dennis-dko/go-toolkit/logging"

	"github.com/labstack/echo/v4/middleware"

	"github.com/dennis-dko/go-toolkit/util"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
)

const (
	PathTag  = "param"
	QueryTag = "query"
)

type SlogAdapter struct {
	mu     sync.Mutex
	Ctx    context.Context
	Logger *slog.Logger
}

type Config struct {
	BaseURL       string        `env:"REST_CLIENT_BASE_URL,notEmpty"`
	Timeout       time.Duration `env:"REST_CLIENT_TIMEOUT" envDefault:"60s"`
	Username      string        `env:"REST_CLIENT_USERNAME,unset"`
	Password      string        `env:"REST_CLIENT_PASSWORD,unset"`
	Token         string        `env:"REST_CLIENT_TOKEN,unset"`
	ContentLength bool          `env:"REST_CLIENT_CONTENT_LENGTH"`
	TLSConfig     tls.Config
	Cookies       []*http.Cookie
}

type HttpHandler struct {
	Client *resty.Client
}

// New creates a new instance of HttpHandler
func New(ctx context.Context, cfg *Config) *HttpHandler {
	httpHandler := &HttpHandler{
		Client: resty.New(),
	}
	httpHandler.setConfig(ctx, cfg)
	return httpHandler
}

type HttpRequest struct {
	Method                string
	URL                   string
	ForceContentType      string
	Headers               map[string]string
	PathParams            map[string]string
	QueryParams           map[string]string
	QueryParamsFromValues map[string][]string
	FormData              map[string]string
	Body                  interface{}
	DestResult            interface{}
}

// DoHTTPRequest executes the http request
func (h *HttpHandler) DoHTTPRequest(data *HttpRequest) (*resty.Response, error) {
	var (
		err      error
		response *resty.Response
	)
	requestClient := h.buildRequest(data)
	switch data.Method {
	case http.MethodGet:
		response, err = requestClient.Get(data.URL)
		if err != nil {
			return nil, err
		}
	case http.MethodPost:
		response, err = requestClient.Post(data.URL)
		if err != nil {
			return nil, err
		}
	case http.MethodPut:
		response, err = requestClient.Put(data.URL)
		if err != nil {
			return nil, err
		}
	case http.MethodDelete:
		response, err = requestClient.Delete(data.URL)
		if err != nil {
			return nil, err
		}
	case http.MethodHead:
		response, err = requestClient.Head(data.URL)
		if err != nil {
			return nil, err
		}
	case http.MethodOptions:
		response, err = requestClient.Options(data.URL)
		if err != nil {
			return nil, err
		}
	case http.MethodPatch:
		response, err = requestClient.Patch(data.URL)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid method type to create http request")
	}
	return response, nil
}

// UseRequestID generate the request id
func UseRequestID(ctx context.Context, instance *echo.Echo) {
	instance.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			uuid := util.SetUUID()
			slog.DebugContext(ctx, "GENERATE_REQUEST_ID",
				slog.String("uuid", uuid),
			)
			return uuid
		},
	}))
	instance.Use(RequestIDtToContext)
}

// RequestIDtToContext sets the request id in the context
func RequestIDtToContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := c.Request().Header.Get(echo.HeaderXRequestID)
		if requestID == "" {
			requestID = c.Response().Header().Get(echo.HeaderXRequestID)
		}
		logCtx := logging.AppendCtx(
			c.Request().Context(),
			slog.String("id", requestID),
		)
		ctx := context.WithValue(
			logCtx,
			echo.HeaderXRequestID,
			requestID,
		)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

// Errorf logs an error message
func (s *SlogAdapter) Errorf(format string, v ...interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Logger.ErrorContext(s.Ctx, "error while using http client", slog.String("id", getRequestID(v...)),
		slog.String("data", fmt.Sprintf(format, v...)))
}

// Warnf logs a warning message
func (s *SlogAdapter) Warnf(format string, v ...interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Logger.WarnContext(s.Ctx, "Warning while using http client", slog.String("id", getRequestID(v...)),
		slog.String("data", fmt.Sprintf(format, v...)))
}

// Debugf logs a debug message
func (s *SlogAdapter) Debugf(format string, v ...interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Logger.DebugContext(s.Ctx, "Debugging while using http client", slog.String("id", getRequestID(v...)),
		slog.String("data", fmt.Sprintf(format, v...)))
}

// Close closes opened http body correctly
// logs an error if the body cannot be closed
// the response must not be nil
func Close(ctx context.Context, response *http.Response) {
	if response != nil && response.Body != nil {
		err := response.Body.Close()
		if err != nil {
			slog.ErrorContext(ctx, "error closing response body", slog.String("error", err.Error()))
		}
	}
}

// GetHeaderCtxValue returns the header value of the context
func GetHeaderCtxValue(ctx context.Context, ctxKey string) string {
	if value, ok := ctx.Value(ctxKey).(string); ok {
		return value
	}
	return ""
}

// GetParams returns the parameters of a struct
func GetParams(rawStruct interface{}, onlySlice bool, tags ...string) interface{} {
	structData := reflect.ValueOf(rawStruct)
	if structData.Kind() == reflect.Ptr {
		structData = reflect.ValueOf(rawStruct).Elem()
	}
	var (
		stringValues = make(map[string]string)
		sliceValues  = make(map[string][]string)
	)
	if structData.Kind() == reflect.Struct {
		if tags == nil {
			tags = []string{
				PathTag,
				QueryTag,
			}
		}
		for i := 0; i < structData.NumField(); i++ {
			if !structData.Field(i).CanInterface() {
				continue
			}
			if (structData.Field(i).Kind() == reflect.Ptr && structData.Field(i).IsNil()) || structData.Field(i).IsZero() {
				continue
			}
			fieldName := getFieldNameByTag(structData, i, tags...)
			if fieldName == "" {
				continue
			}
			if onlySlice {
				slice := getSlice(structData.Field(i).Interface())
				if slice != nil {
					sliceValues[fieldName] = slice
					continue
				}
			} else {
				nullType := getNullTypes(structData.Field(i).Interface())
				if nullType != nil {
					stringValues[fieldName] = *nullType
					continue
				}
				dateTime := getDateTime(structData.Field(i).Interface())
				if dateTime != nil {
					stringValues[fieldName] = *dateTime
					continue
				}
				if structData.Field(i).Kind() == reflect.Ptr {
					ptrValue := structData.Field(i).Elem()
					if ptrValue.Kind() == reflect.Slice {
						continue
					}
					stringValues[fieldName] = fmt.Sprint(ptrValue)
				} else {
					value := structData.Field(i)
					if value.Kind() == reflect.Slice {
						continue
					}
					stringValues[fieldName] = fmt.Sprint(value)
				}
			}
		}
	}
	if onlySlice {
		return sliceValues
	}
	return stringValues
}

func getNullTypes(data interface{}) *string {
	var nullType *string
	switch value := data.(type) {
	case datatype.NullBool:
		nullType = datatype.StringPtr(fmt.Sprintf("%v", value.Bool))
	case *datatype.NullBool:
		nullType = datatype.StringPtr(fmt.Sprintf("%v", value.Bool))
	case datatype.NullFloat64:
		nullType = datatype.StringPtr(fmt.Sprintf("%v", value.Float64))
	case *datatype.NullFloat64:
		nullType = datatype.StringPtr(fmt.Sprintf("%v", value.Float64))
	case datatype.NullInt64:
		nullType = datatype.StringPtr(fmt.Sprintf("%v", value.Int64))
	case *datatype.NullInt64:
		nullType = datatype.StringPtr(fmt.Sprintf("%v", value.Int64))
	case datatype.NullString:
		nullType = datatype.StringPtr(fmt.Sprintf("%v", value.String))
	case *datatype.NullString:
		nullType = datatype.StringPtr(fmt.Sprintf("%v", value.String))
	default:
		return nil
	}
	return nullType
}

func getDateTime(data interface{}) *string {
	var dateTime *string
	switch value := data.(type) {
	case datatype.NullTime:
		dateTime = datatype.StringPtr(value.String())
	case *datatype.NullTime:
		dateTime = datatype.StringPtr(value.String())
	case datatype.NullDate:
		dateTime = datatype.StringPtr(value.String())
	case *datatype.NullDate:
		dateTime = datatype.StringPtr(value.String())
	case datatype.CustomTime:
		dateTime = datatype.StringPtr(value.String())
	case *datatype.CustomTime:
		dateTime = datatype.StringPtr(value.String())
	case datatype.CustomDate:
		dateTime = datatype.StringPtr(value.String())
	case *datatype.CustomDate:
		dateTime = datatype.StringPtr(value.String())
	default:
		return nil
	}
	return dateTime
}

func getSlice(data interface{}) []string {
	var stringSlice []string
	switch value := data.(type) {
	case []datatype.NullBool:
		for _, v := range value {
			stringSlice = append(stringSlice, fmt.Sprintf("%v", v.Bool))
		}
	case *[]datatype.NullBool:
		for _, v := range *value {
			stringSlice = append(stringSlice, fmt.Sprintf("%v", v.Bool))
		}
	case []bool:
		for _, v := range value {
			stringSlice = append(stringSlice, fmt.Sprintf("%v", v))
		}
	case *[]bool:
		for _, v := range *value {
			stringSlice = append(stringSlice, fmt.Sprintf("%v", v))
		}
	case []datatype.NullFloat64:
		for _, v := range value {
			stringSlice = append(stringSlice, fmt.Sprintf("%v", v.Float64))
		}
	case *[]datatype.NullFloat64:
		for _, v := range *value {
			stringSlice = append(stringSlice, fmt.Sprintf("%v", v.Float64))
		}
	case []float64:
		for _, v := range value {
			stringSlice = append(stringSlice, fmt.Sprintf("%v", v))
		}
	case *[]float64:
		for _, v := range *value {
			stringSlice = append(stringSlice, fmt.Sprintf("%v", v))
		}
	case []datatype.NullInt64:
		for _, v := range value {
			stringSlice = append(stringSlice, fmt.Sprintf("%v", v.Int64))
		}
	case *[]datatype.NullInt64:
		for _, v := range *value {
			stringSlice = append(stringSlice, fmt.Sprintf("%v", v.Int64))
		}
	case []int64:
		for _, v := range value {
			stringSlice = append(stringSlice, fmt.Sprintf("%v", v))
		}
	case *[]int64:
		for _, v := range *value {
			stringSlice = append(stringSlice, fmt.Sprintf("%v", v))
		}
	case []datatype.NullString:
		for _, v := range value {
			stringSlice = append(stringSlice, fmt.Sprintf("%v", v.String))
		}
	case *[]datatype.NullString:
		for _, v := range *value {
			stringSlice = append(stringSlice, fmt.Sprintf("%v", v.String))
		}
	case []string:
		stringSlice = value
	case *[]string:
		stringSlice = *value
	default:
		return nil
	}
	return stringSlice
}

func (h *HttpHandler) setConfig(ctx context.Context, cfg *Config) {
	h.Client.SetLogger(
		&SlogAdapter{
			Ctx:    ctx,
			Logger: slog.Default(),
		},
	)
	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		h.Client.SetDebug(true)
	}
	if cfg.BaseURL != "" {
		h.Client.SetBaseURL(cfg.BaseURL)
	}
	if cfg.Timeout > 0 {
		h.Client.SetTimeout(cfg.Timeout)
	}
	if cfg.Username != "" && cfg.Password != "" {
		h.Client.SetBasicAuth(cfg.Username, cfg.Password)
	}
	if cfg.Token != "" {
		h.Client.SetAuthToken(cfg.Token)
	}
	if cfg.ContentLength {
		h.Client.SetContentLength(cfg.ContentLength)
	}
	if cfg.TLSConfig.ServerName != "" {
		h.Client.SetTLSClientConfig(&cfg.TLSConfig)
	}
	if len(cfg.Cookies) > 0 {
		h.Client.SetCookies(cfg.Cookies)
	}
}

func (h *HttpHandler) buildRequest(data *HttpRequest) *resty.Request {
	instance := h.Client.R()
	if data.ForceContentType != "" {
		instance.ForceContentType(data.ForceContentType)
	}
	if data.Headers == nil {
		data.Headers = make(map[string]string)
	}
	if _, ok := data.Headers[echo.HeaderXRequestID]; !ok {
		data.Headers[echo.HeaderXRequestID] = util.SetUUID()
	}
	instance.SetHeaders(data.Headers)
	if data.PathParams != nil {
		instance.SetPathParams(data.PathParams)
	}
	if data.QueryParams != nil {
		instance.SetQueryParams(data.QueryParams)
	}
	if data.QueryParamsFromValues != nil {
		instance.SetQueryParamsFromValues(data.QueryParamsFromValues)
	}
	if (data.Method == http.MethodPost || data.Method == http.MethodPut) && data.FormData != nil {
		instance.SetFormData(data.FormData)
	}
	if data.Body != nil {
		instance.SetBody(data.Body)
	}
	if data.DestResult != nil {
		instance.SetResult(data.DestResult)
	}
	return instance
}

func getFieldNameByTag(data reflect.Value, fieldIndex int, tags ...string) string {
	var fieldName string
	for _, tag := range tags {
		fieldName = strings.SplitN(data.Type().Field(fieldIndex).Tag.Get(tag), ",", 2)[0]
		if fieldName == "" || fieldName == "-" {
			continue
		}
		break
	}
	return fieldName
}

func getRequestID(v ...interface{}) string {
	for _, data := range v {
		if msg, ok := data.(string); ok {
			m := regexp.MustCompile(
				fmt.Sprintf(`%s: ([a-f0-9\-]{36})`, echo.HeaderXRequestID),
			)
			requestID := m.FindStringSubmatch(msg)
			if requestID == nil {
				continue
			}
			return requestID[1]
		}
	}
	return ""
}
