package datatype

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type CustomDateTest struct {
	Date        CustomDate
	DatePointer *CustomDate
}

type CustomDateTestSuite struct {
	suite.Suite
	dateString        string
	invalidDateString string
	customDateTest    CustomDateTest
	currentDate       CustomDate
	xmlData           []byte
	invalidXmlData    []byte
	jsonData          []byte
	invalidJsonData   []byte
}

func (cd *CustomDateTestSuite) SetupTest() {
	// Setup
	newDate, _ := NewDate(false)
	cd.currentDate = *newDate
	cd.invalidDateString = "invalid_date"
	cd.customDateTest = CustomDateTest{
		Date:        cd.currentDate,
		DatePointer: nil,
	}
	cd.dateString = cd.currentDate.String()
	cd.xmlData = []byte(fmt.Sprintf("<CustomDate>%s</CustomDate>", cd.dateString))
	cd.invalidXmlData = []byte(fmt.Sprintf("<CustomDate>%s</CustomDate>", cd.invalidDateString))
	cd.jsonData = []byte(fmt.Sprintf(`{"Date":"%s","DatePointer":null}`, cd.dateString))
	cd.invalidJsonData = []byte(fmt.Sprintf(`{"Date":"%s"}`, cd.invalidDateString))
}

func TestCustomDateTestSuite(t *testing.T) {
	suite.Run(t, new(CustomDateTestSuite))
}

func (cd *CustomDateTestSuite) TestMarshalDate() {

	cd.Run("happy path - date as xml is marshalled", func() {
		// Run
		data, err := xml.Marshal(cd.customDateTest.Date)

		// Assert
		cd.NoError(err)
		cd.Equal(cd.xmlData, data)
	})
	cd.Run("happy path - date as json is marshalled", func() {
		// Run
		data, err := json.Marshal(cd.customDateTest)

		// Assert
		cd.NoError(err)
		cd.Equal(cd.jsonData, data)
	})
}

func (cd *CustomDateTestSuite) TestUnmarshalDate() {

	cd.Run("happy path - date as xml is unmarshalled", func() {
		// Init
		var customDate CustomDate

		// Run
		err := xml.Unmarshal(cd.xmlData, &customDate)

		// Assert
		cd.NoError(err)
		cd.Equal(cd.currentDate, customDate)
	})
	cd.Run("should return an error while unmarshalling date as xml", func() {
		// Init
		var customDate CustomDate

		// Run
		err := xml.Unmarshal(cd.invalidXmlData, &customDate)

		// Assert
		cd.Error(err)
		cd.Empty(customDate.Time)
	})
	cd.Run("happy path - date as json is unmarshalled", func() {
		// Init
		var data CustomDateTest

		// Run
		err := json.Unmarshal(cd.jsonData, &data)

		// Assert
		cd.NoError(err)
		cd.Equal(cd.customDateTest, data)
	})
	cd.Run("should return an error while unmarshalling date as json", func() {
		// Init
		var data CustomDateTest

		// Run
		err := json.Unmarshal(cd.invalidJsonData, &data)

		// Assert
		cd.Error(err)
		cd.Empty(data.Date.Time)
	})
}

func (cd *CustomDateTestSuite) TestValueScanDate() {

	cd.Run("happy path - date is scanned and valued", func() {
		// Init
		var data CustomDate

		// Run
		err := data.Scan(time.Now().UTC())
		value, _ := data.Value()

		// Assert
		cd.NoError(err)
		cd.NotEmpty(value)
	})
	cd.Run("should return an error while scanning and valuing date", func() {
		// Init
		var data CustomDate

		// Run
		err := data.Scan(cd.invalidDateString)
		value, _ := data.Value()

		// Assert
		cd.Error(err)
		cd.Empty(data.Time)
		cd.Nil(value)
	})
}

func (cd *CustomDateTestSuite) TestStringDate() {

	cd.Run("happy path - date is stringified", func() {
		// Run
		dateString := cd.customDateTest.Date.String()

		// Assert
		cd.NotEmpty(dateString)
	})
}

func (cd *CustomDateTestSuite) TestSubDate() {

	cd.Run("happy path - date is subtracted", func() {
		// Run
		subDate := cd.customDateTest.Date.SubDate(&cd.currentDate)

		// Assert
		cd.Empty(subDate)
	})
}

func (cd *CustomDateTestSuite) TestNewDate() {

	cd.Run("happy path - date is new created", func() {
		// Run
		newDate, err := NewDate(false)

		// Assert
		cd.NoError(err)
		cd.NotNil(newDate)
	})
}

func (cd *CustomDateTestSuite) TestParseDate() {

	cd.Run("happy path - date is parsed", func() {
		// Run
		parsedDate, err := ParseDate(cd.dateString, false)

		// Assert
		cd.NoError(err)
		cd.Equal(cd.currentDate, *parsedDate)
	})
	cd.Run("should return an error while parsing date", func() {
		// Run
		parsedDate, err := ParseDate(cd.invalidDateString, false)

		// Assert
		cd.Error(err)
		cd.Nil(parsedDate)
	})
}
