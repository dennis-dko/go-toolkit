package datatype

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NullSQLTest struct {
	Bool           NullBool
	BoolPointer    *NullBool
	Float64        NullFloat64
	Float64Pointer *NullFloat64
	Int64          NullInt64
	Int64Pointer   *NullInt64
	String         NullString
	StringPointer  *NullString
	Time           NullTime
	TimePointer    *NullTime
	Date           NullDate
	DatePointer    *NullDate
}

type NullSQLTestSuite struct {
	suite.Suite
	nullSQLTest     NullSQLTest
	jsonData        []byte
	invalidJsonData []byte
	xmlData         []byte
	invalidXmlData  []byte
}

func (n *NullSQLTestSuite) SetupTest() {
	// Setup
	currentTime, _ := NewTime(false)
	currentDate, _ := NewDate(false)
	n.nullSQLTest = NullSQLTest{
		Bool: NullBool{
			NullBool: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		},
		Float64: NullFloat64{
			NullFloat64: sql.NullFloat64{
				Float64: 5.6,
				Valid:   true,
			},
		},
		Int64: NullInt64{
			NullInt64: sql.NullInt64{
				Int64: 5,
				Valid: true,
			},
		},
		String: NullString{
			NullString: sql.NullString{
				String: "String",
				Valid:  true,
			},
		},
		Time: NullTime{
			NullTime: sql.NullTime{
				Time:  currentTime.Time,
				Valid: true,
			},
			NoUtc:  false,
			Format: DefaultTimeFormat,
		},
		Date: NullDate{
			NullTime: sql.NullTime{
				Time:  currentDate.Time,
				Valid: true,
			},
			NoUtc:  false,
			Format: DefaultDateFormat,
		},
	}
	n.jsonData = []byte(fmt.Sprintf(`{"Bool":true,"BoolPointer":null,"Float64":5.6,"Float64Pointer":null,"Int64":5,"Int64Pointer":null,"String":"String","StringPointer":null,"Time":"%s","TimePointer":null,"Date":"%s","DatePointer":null}`, currentTime, currentDate))
	n.invalidJsonData = []byte(`{"Bool":"invalid_bool","BoolPointer":null,"Float64":"invalid_float","Float64Pointer":null,"Int64":"invalid_int64","Int64Pointer":null,"String":0,"StringPointer":null,"Time":"invalid_time","TimePointer":null,"Date":"invalid_date","DatePointer":null}`)
	n.xmlData = []byte(fmt.Sprintf(`<NullSQLTest><Bool>true</Bool><Float64>5.6</Float64><Int64>5</Int64><String>String</String><Time>%s</Time><Date>%s</Date></NullSQLTest>`, currentTime, currentDate))
	n.invalidXmlData = []byte(`<NullSQLTest><Bool>true</Bool><Float64>5.6</Float64><Int64>5</Int64><String>String</String><Time>invalid_time</Time><Date>invalid_date</Date></NullSQLTest>`)
}

func TestNullSQLTestSuite(t *testing.T) {
	suite.Run(t, new(NullSQLTestSuite))
}

func (n *NullSQLTestSuite) TestNullSQLJSONMarshal() {

	n.Run("happy path - all struct fields in struct could be marshalled to json", func() {
		// Run
		data, err := json.Marshal(n.nullSQLTest)

		// Assert
		n.NoError(err)
		n.Equal(n.jsonData, data)
	})
}

func (n *NullSQLTestSuite) TestNullSQLJSONUnmarshal() {

	n.Run("happy path - all json keys could be unmarshalled", func() {
		// Init
		var data NullSQLTest

		// Run
		err := json.Unmarshal(n.jsonData, &data)

		// Assert
		n.NoError(err)
		n.Equal(n.nullSQLTest, data)
	})
	n.Run("should return an error while unmarshalling all json keys", func() {
		// Init
		var data NullSQLTest

		// Run
		err := json.Unmarshal(n.invalidJsonData, &data)

		// Assert
		n.Error(err)
		n.Empty(data)
	})
}

func (n *NullSQLTestSuite) TestNullSQLXMLMarshal() {

	n.Run("happy path - all struct fields in struct could be marshalled to xml", func() {
		// Run
		data, err := xml.Marshal(n.nullSQLTest)

		// Assert
		n.NoError(err)
		n.Equal(n.xmlData, data)
	})
}

func (n *NullSQLTestSuite) TestNullSQLXMLUnmarshal() {

	n.Run("happy path - all xml tags could be unmarshalled", func() {
		// Init
		var data NullSQLTest

		// Run
		err := xml.Unmarshal(n.xmlData, &data)

		// Assert
		n.NoError(err)
		n.Equal(n.nullSQLTest, data)
	})
	n.Run("should return an error while unmarshalling all xml tags", func() {
		// Init
		var data NullSQLTest

		// Run
		err := xml.Unmarshal(n.invalidXmlData, &data)

		// Assert
		n.Error(err)
		n.NotEmpty(data)
	})
}

func (n *NullSQLTestSuite) TestNewNullBool() {

	n.Run("happy path - new null bool created", func() {
		// Init
		testData := n.nullSQLTest.Bool.NullBool.Bool

		// Run
		nullBool := NewNullBool(&testData)

		// Assert
		n.Equal(testData, nullBool.Bool)
	})
}

func (n *NullSQLTestSuite) TestNullBoolPtrToBoolPtr() {

	n.Run("happy path - null bool to bool pointer converted", func() {
		// Init
		testData := n.nullSQLTest.Bool.NullBool.Bool

		// Run
		boolPointer := NullBoolPtrToBoolPtr(&n.nullSQLTest.Bool)

		// Assert
		n.Equal(BoolPtr(testData), boolPointer)
	})
}

func (n *NullSQLTestSuite) TestNewNullFloat64() {

	n.Run("happy path - new null float64 created", func() {
		// Init
		testData := n.nullSQLTest.Float64.NullFloat64.Float64

		// Run
		nullFloat64 := NewNullFloat64(&testData)

		// Assert
		n.Equal(testData, nullFloat64.Float64)
	})
}

func (n *NullSQLTestSuite) TestNullFloat64PtrToFloat64Ptr() {

	n.Run("happy path - null float64 to float64 pointer converted", func() {
		// Init
		testData := n.nullSQLTest.Float64.NullFloat64.Float64

		// Run
		float64Pointer := NullFloat64PtrToFloat64Ptr(&n.nullSQLTest.Float64)

		// Assert
		n.Equal(Float64Ptr(testData), float64Pointer)
	})
}

func (n *NullSQLTestSuite) TestNewNullInt64() {

	n.Run("happy path - new null int64 created", func() {
		// Init
		testData := n.nullSQLTest.Int64.NullInt64.Int64

		// Run
		nullInt64 := NewNullInt64(&testData)

		// Assert
		n.Equal(testData, nullInt64.Int64)
	})
}

func (n *NullSQLTestSuite) TestNullInt64PtrToInt64Ptr() {

	n.Run("happy path - null int64 to int64 pointer converted", func() {
		// Init
		testData := n.nullSQLTest.Int64.NullInt64.Int64

		// Run
		int64Pointer := NullInt64PtrToInt64Ptr(&n.nullSQLTest.Int64)

		// Assert
		n.Equal(Int64Ptr(testData), int64Pointer)
	})
}

func (n *NullSQLTestSuite) TestNewNullString() {

	n.Run("happy path - new null string created", func() {
		// Init
		testData := n.nullSQLTest.String.NullString.String

		// Run
		nullString := NewNullString(&testData)

		// Assert
		n.Equal(testData, nullString.String)
	})
}

func (n *NullSQLTestSuite) TestNullStringPtrToStringPtr() {

	n.Run("happy path - null string to string pointer converted", func() {
		// Init
		testData := n.nullSQLTest.String.NullString.String

		// Run
		stringPointer := NullStringPtrToStringPtr(&n.nullSQLTest.String)

		// Assert
		n.Equal(StringPtr(testData), stringPointer)
	})
}

func (n *NullSQLTestSuite) TestNewNullTime() {

	n.Run("happy path - new null time created", func() {
		// Init
		testData := CustomTime{
			Time: n.nullSQLTest.Time.NullTime.Time,
		}

		// Run
		nullTime := NewNullTime(&testData)

		// Assert
		n.Equal(testData.Time, nullTime.Time)
	})
}

func (n *NullSQLTestSuite) TestNullTimePtrToCustomTimePtr() {

	n.Run("happy path - null time to custom time pointer converted", func() {
		// Init
		testData := CustomTime{
			Time:   n.nullSQLTest.Time.NullTime.Time,
			NoUtc:  n.nullSQLTest.Time.NoUtc,
			Format: n.nullSQLTest.Time.Format,
		}

		// Run
		customTimePointer := NullTimePtrToCustomTimePtr(&n.nullSQLTest.Time)

		// Assert
		n.Equal(&testData, customTimePointer)
	})
}

func (n *NullSQLTestSuite) TestNewNullDate() {

	n.Run("happy path - new null date created", func() {
		// Init
		testData := CustomDate{
			Time: n.nullSQLTest.Date.NullTime.Time,
		}

		// Run
		nullDate := NewNullDate(&testData)

		// Assert
		n.Equal(testData.Time, nullDate.Time)
	})
}

func (n *NullSQLTestSuite) TestNullDatePtrToCustomDatePtr() {

	n.Run("happy path - null date to custom date pointer converted", func() {
		// Init
		testData := CustomDate{
			Time:   n.nullSQLTest.Date.NullTime.Time,
			NoUtc:  n.nullSQLTest.Date.NoUtc,
			Format: n.nullSQLTest.Date.Format,
		}

		// Run
		customDatePointer := NullDatePtrToCustomDatePtr(&n.nullSQLTest.Date)

		// Assert
		n.Equal(&testData, customDatePointer)
	})
}
