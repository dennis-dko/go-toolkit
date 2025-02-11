package datatype

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type CustomTimeTest struct {
	Time        CustomTime
	TimePointer *CustomTime
}

type CustomTimeTestSuite struct {
	suite.Suite
	timeString        string
	invalidTimeString string
	customTimeTest    CustomTimeTest
	currentTime       CustomTime
	xmlData           []byte
	invalidXmlData    []byte
	jsonData          []byte
	invalidJsonData   []byte
}

func (ct *CustomTimeTestSuite) SetupTest() {
	// Setup
	newTime, _ := NewTime(false)
	ct.currentTime = *newTime
	ct.invalidTimeString = "invalid_time"
	ct.customTimeTest = CustomTimeTest{
		Time:        ct.currentTime,
		TimePointer: nil,
	}
	ct.timeString = ct.currentTime.String()
	ct.xmlData = []byte(fmt.Sprintf("<CustomTime>%s</CustomTime>", ct.timeString))
	ct.invalidXmlData = []byte(fmt.Sprintf("<CustomTime>%s</CustomTime>", ct.invalidTimeString))
	ct.jsonData = []byte(fmt.Sprintf(`{"Time":"%s","TimePointer":null}`, ct.timeString))
	ct.invalidJsonData = []byte(fmt.Sprintf(`{"Time":"%s"}`, ct.invalidTimeString))
}

func TestCustomTimeTestSuite(t *testing.T) {
	suite.Run(t, new(CustomTimeTestSuite))
}

func (ct *CustomTimeTestSuite) TestMarshalTime() {

	ct.Run("happy path - time as xml is marshalled", func() {
		// Run
		data, err := xml.Marshal(ct.customTimeTest.Time)

		// Assert
		ct.NoError(err)
		ct.Equal(ct.xmlData, data)
	})
	ct.Run("happy path - time as json is marshalled", func() {
		// Run
		data, err := json.Marshal(ct.customTimeTest)

		// Assert
		ct.NoError(err)
		ct.Equal(ct.jsonData, data)
	})
}

func (ct *CustomTimeTestSuite) TestUnmarshalTime() {

	ct.Run("happy path - time as xml is unmarshalled", func() {
		// Init
		var customTime CustomTime

		// Run
		err := xml.Unmarshal(ct.xmlData, &customTime)

		// Assert
		ct.NoError(err)
		ct.Equal(ct.currentTime, customTime)
	})
	ct.Run("should return an error while unmarshalling time as xml", func() {
		// Init
		var customTime CustomTime

		// Run
		err := xml.Unmarshal(ct.invalidXmlData, &customTime)

		// Assert
		ct.Error(err)
		ct.Empty(customTime.Time)
	})
	ct.Run("happy path - time as json is unmarshalled", func() {
		// Init
		var data CustomTimeTest

		// Run
		err := json.Unmarshal(ct.jsonData, &data)

		// Assert
		ct.NoError(err)
		ct.Equal(ct.customTimeTest, data)
	})
	ct.Run("should return an error while unmarshalling time as json", func() {
		// Init
		var data CustomTimeTest

		// Run
		err := json.Unmarshal(ct.invalidJsonData, &data)

		// Assert
		ct.Error(err)
		ct.Empty(data.Time.Time)
	})
}

func (ct *CustomTimeTestSuite) TestValueScanTime() {

	ct.Run("happy path - time is scanned and valued", func() {
		// Init
		var data CustomTime

		// Run
		err := data.Scan(time.Now().UTC())
		value, _ := data.Value()

		// Assert
		ct.NoError(err)
		ct.NotEmpty(value)
	})
	ct.Run("should return an error while scanning and valuing time", func() {
		// Init
		var data CustomTime

		// Run
		err := data.Scan(ct.invalidTimeString)
		value, _ := data.Value()

		// Assert
		ct.Error(err)
		ct.Empty(data.Time)
		ct.Nil(value)
	})
}

func (ct *CustomTimeTestSuite) TestStringTime() {

	ct.Run("happy path - time is stringified", func() {
		// Run
		timeString := ct.customTimeTest.Time.String()

		// Assert
		ct.NotEmpty(timeString)
	})
}

func (ct *CustomTimeTestSuite) TestSubTime() {

	ct.Run("happy path - time is subtracted", func() {
		// Run
		subTime := ct.customTimeTest.Time.SubTime(&ct.currentTime)

		// Assert
		ct.Empty(subTime)
	})
}

func (ct *CustomTimeTestSuite) TestNewTime() {

	ct.Run("happy path - time is new created", func() {
		// Run
		newTime, err := NewTime(false)

		// Assert
		ct.NoError(err)
		ct.NotNil(newTime)
	})
}

func (ct *CustomTimeTestSuite) TestParseTime() {

	ct.Run("happy path - time is parsed", func() {
		// Run
		parsedTime, err := ParseTime(ct.timeString, false)

		// Assert
		ct.NoError(err)
		ct.Equal(ct.currentTime, *parsedTime)
	})
	ct.Run("should return an error while parsing time", func() {
		// Run
		parsedTime, err := ParseTime(ct.invalidTimeString, false)

		// Assert
		ct.Error(err)
		ct.Nil(parsedTime)
	})
}
