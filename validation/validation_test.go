package validation

import (
	"testing"

	"github.com/dennis-dko/go-toolkit/datatype"
	"github.com/dennis-dko/go-toolkit/testhandler"

	"github.com/stretchr/testify/suite"
)

type StructTest struct {
	FirstName      string              `validate:"required"`
	LastName       string              `validate:"required"`
	Age            uint8               `validate:"gte=0,lte=130"`
	Email          string              `validate:"required,email"`
	Gender         string              `validate:"oneof=male female prefer_not_to"`
	FavouriteColor string              `validate:"iscolor"`
	Hobby          datatype.NullString `validate:"required,oneof=Cooking Driving"`
	Street         string              `validate:"depends_on=PostalCode"`
	PostalCode     int                 `validate:"depends_on=Street"`
	Birth          datatype.CustomDate `validate:"required"`
}

type ValidationTestSuite struct {
	suite.Suite
	validator  *RequestValidator
	structTest StructTest
}

func (v *ValidationTestSuite) SetupTest() {
	// Setup
	v.validator = New(testhandler.Ctx(false, false))
}

func (v *ValidationTestSuite) SetupSubTest() {
	// Sub setup
	currentDate, _ := datatype.NewDate(false)
	v.structTest = StructTest{
		FirstName:      "Walther",
		LastName:       "White",
		Age:            54,
		Email:          "test@example.com",
		Gender:         "male",
		FavouriteColor: "#348feb",
		Hobby: datatype.NewNullString(
			datatype.StringPtr("Cooking"),
		),
		Street:     "example",
		PostalCode: 1234,
		Birth:      *currentDate,
	}
}

func TestValidationTestSuite(t *testing.T) {
	suite.Run(t, new(ValidationTestSuite))
}

func (v *ValidationTestSuite) TestValidation() {

	v.Run("happy path - struct is valid", func() {
		// Run
		err := v.validator.Validate(v.structTest)

		// Assert
		v.NoError(err)
	})
	v.Run("should return an error while street is missing", func() {
		// Init
		v.structTest.Birth = datatype.CustomDate{}

		// Run
		err := v.validator.Validate(v.structTest)

		// Assert
		v.Error(err)
		v.ErrorContains(err, "Field validation for 'Birth' failed")
	})
}
