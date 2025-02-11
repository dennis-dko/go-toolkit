package validation

import (
	"context"
	"database/sql/driver"
	"log/slog"
	"os"
	"reflect"
	"strings"

	"github.com/dennis-dko/go-toolkit/datatype"

	"github.com/go-playground/validator/v10"
)

type RequestValidator struct {
	ctx       context.Context
	validator *validator.Validate
}

// New creates a new instance of RequestValidator
func New(ctx context.Context) *RequestValidator {
	validate := &RequestValidator{
		ctx: ctx,
		validator: validator.New(
			validator.WithRequiredStructEnabled(),
		),
	}
	validate.register()
	return validate
}

// Validate validates the given struct
func (r RequestValidator) Validate(i interface{}) error {
	if err := r.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func (r RequestValidator) register() {
	r.validator.RegisterCustomTypeFunc(
		validateValuer,
		datatype.NullBool{},
		datatype.NullFloat64{},
		datatype.NullInt64{},
		datatype.NullString{},
		datatype.NullTime{},
		datatype.NullDate{},
		datatype.CustomTime{},
		datatype.CustomDate{},
	)
	err := r.validator.RegisterValidationCtx("depends_on", validateDependsOn)
	if err != nil {
		slog.ErrorContext(r.ctx, "error while register depends on validation", slog.String("error", err.Error()))
		os.Exit(1)
	}
	err = r.validator.RegisterValidationCtx("depends_one_of", validateDependsOneOf)
	if err != nil {
		slog.ErrorContext(r.ctx, "error while register one of validation", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func validateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
	}
	return nil
}

func validateDependsOn(ctx context.Context, fl validator.FieldLevel) bool {
	parent := fl.Parent()
	tag := fl.Param()
	fields := strings.Split(tag, " ")
	for _, fieldName := range fields {
		field := parent.FieldByName(fieldName)
		if isFieldEmpty(field) {
			slog.InfoContext(ctx, "This field needs to be set", slog.String("fieldName", fieldName))
			return false
		}
	}
	return true
}

func validateDependsOneOf(ctx context.Context, fl validator.FieldLevel) bool {
	parent := fl.Parent()
	tag := fl.Param()
	fields := strings.Split(tag, " ")
	for _, fieldName := range fields {
		field := parent.FieldByName(fieldName)
		if !isFieldEmpty(field) {
			return true
		}
	}
	slog.InfoContext(ctx, "At least one of these fields needs to be set", slog.String("fields", strings.Join(fields, ", ")))
	return false
}

func isFieldEmpty(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.String:
		return field.String() == ""
	case reflect.Ptr:
		return field.IsNil() || isFieldEmpty(field.Elem())
	case reflect.Slice, reflect.Array:
		return field.Len() == 0
	default:
		return field.IsZero()
	}
}
