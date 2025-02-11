package datatype

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PointerTest struct {
	Bool        bool
	NullBool    NullBool
	Float64     float64
	NullFloat64 NullFloat64
	Int64       int64
	NullInt64   NullInt64
	Int         int
	String      string
	NullString  NullString
	Time        CustomTime
	NullTime    NullTime
	Date        CustomDate
	NullDate    NullDate
}

type PointerTestSuite struct {
	suite.Suite
	pointerTest    PointerTest
	nilPointerTest *PointerTest
}

func (p *PointerTestSuite) SetupTest() {
	// Setup
	currentTime, _ := NewTime(false)
	currentDate, _ := NewDate(false)
	p.pointerTest = PointerTest{
		Bool: true,
		NullBool: NullBool{
			sql.NullBool{
				Valid: true,
				Bool:  true,
			},
		},
		Float64: 5.7,
		NullFloat64: NullFloat64{
			sql.NullFloat64{
				Valid:   true,
				Float64: 5.7,
			},
		},
		Int64: 56,
		NullInt64: NullInt64{
			sql.NullInt64{
				Valid: true,
				Int64: 56,
			},
		},
		Int:    5,
		String: "pointer",
		NullString: NullString{
			sql.NullString{
				Valid:  true,
				String: "pointer",
			},
		},
		Time: *currentTime,
		NullTime: NullTime{
			NullTime: sql.NullTime{
				Valid: true,
				Time:  currentTime.Time,
			},
		},
		Date: *currentDate,
		NullDate: NullDate{
			NullTime: sql.NullTime{
				Valid: true,
				Time:  currentDate.Time,
			},
		},
	}
	p.nilPointerTest = &PointerTest{
		Bool: true,
	}
}

func TestPointerTestSuite(t *testing.T) {
	suite.Run(t, new(PointerTestSuite))
}

func (p *PointerTestSuite) TestBoolPtr() {

	p.Run("happy path - return pointer bool", func() {
		// Init
		testData := p.pointerTest.Bool

		// Run
		pointerBool := BoolPtr(testData)

		// Assert
		p.Equal(&testData, pointerBool)
	})
}

func (p *PointerTestSuite) TestNullBoolPtr() {

	p.Run("happy path - return pointer null bool", func() {
		// Init
		testData := p.pointerTest.NullBool

		// Run
		pointerNullBool := NullBoolPtr(testData)

		// Assert
		p.Equal(&testData, pointerNullBool)
	})
}

func (p *PointerTestSuite) TestFloat64Ptr() {

	p.Run("happy path - return pointer float64", func() {
		// Init
		testData := p.pointerTest.Float64

		// Run
		pointerFloat64 := Float64Ptr(testData)

		// Assert
		p.Equal(&testData, pointerFloat64)
	})
}

func (p *PointerTestSuite) TestNullFloat64Ptr() {

	p.Run("happy path - return pointer null float64", func() {
		// Init
		testData := p.pointerTest.NullFloat64

		// Run
		pointerNullFloat64 := NullFloat64Ptr(testData)

		// Assert
		p.Equal(&testData, pointerNullFloat64)
	})
}

func (p *PointerTestSuite) TestInt64Ptr() {

	p.Run("happy path - return pointer int64", func() {
		// Init
		testData := p.pointerTest.Int64

		// Run
		pointerInt64 := Int64Ptr(testData)

		// Assert
		p.Equal(&testData, pointerInt64)
	})
}

func (p *PointerTestSuite) TestNullInt64Ptr() {

	p.Run("happy path - return pointer null int64", func() {
		// Init
		testData := p.pointerTest.NullInt64

		// Run
		pointerNullInt64 := NullInt64Ptr(testData)

		// Assert
		p.Equal(&testData, pointerNullInt64)
	})
}

func (p *PointerTestSuite) TestIntPtr() {

	p.Run("happy path - return pointer int", func() {
		// Init
		testData := p.pointerTest.Int

		// Run
		pointerInt := IntPtr(testData)

		// Assert
		p.Equal(&testData, pointerInt)
	})
}

func (p *PointerTestSuite) TestStringPtr() {

	p.Run("happy path - return pointer string", func() {
		// Init
		testData := p.pointerTest.String

		// Run
		pointerString := StringPtr(testData)

		// Assert
		p.Equal(&testData, pointerString)
	})
}

func (p *PointerTestSuite) TestNullStringPtr() {

	p.Run("happy path - return pointer null string", func() {
		// Init
		testData := p.pointerTest.NullString

		// Run
		pointerNullString := NullStringPtr(testData)

		// Assert
		p.Equal(&testData, pointerNullString)
	})
}

func (p *PointerTestSuite) TestTimePtr() {

	p.Run("happy path - return pointer time", func() {
		// Init
		testData := p.pointerTest.Time

		// Run
		pointerTime := TimePtr(testData)

		// Assert
		p.Equal(&testData, pointerTime)
	})
}

func (p *PointerTestSuite) TestNullTimePtr() {

	p.Run("happy path - return pointer time", func() {
		// Init
		testData := p.pointerTest.NullTime

		// Run
		pointerNullTime := NullTimePtr(testData)

		// Assert
		p.Equal(&testData, pointerNullTime)
	})
}

func (p *PointerTestSuite) TestDatePtr() {

	p.Run("happy path - return pointer date", func() {
		// Init
		testData := p.pointerTest.Date

		// Run
		pointerDate := DatePtr(testData)

		// Assert
		p.Equal(&testData, pointerDate)
	})
}

func (p *PointerTestSuite) TestNullDatePtr() {

	p.Run("happy path - return pointer null date", func() {
		// Init
		testData := p.pointerTest.NullDate

		// Run
		pointerNullDate := NullDatePtr(testData)

		// Assert
		p.Equal(&testData, pointerNullDate)
	})
}

func (p *PointerTestSuite) TestCheckPtrFieldValues() {

	p.Run("happy path - all field values are nil or zero value except Bool", func() {
		// Init
		testData := p.nilPointerTest

		// Run
		allSet, err := CheckPtrFieldValues(testData, "Bool")

		// Assert
		p.NoError(err)
		p.True(*allSet)
	})
}
