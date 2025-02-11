package datatype

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConverterTestSuite struct {
	suite.Suite
	currentTime CustomTime
	currentDate CustomDate
}

func TestConverterTestSuite(t *testing.T) {
	suite.Run(t, new(ConverterTestSuite))
}

func (c *ConverterTestSuite) SetupTest() {
	// Setup
	newDate, _ := ParseDate("2024-01-26", false)
	c.currentDate = *newDate
	newTime, _ := ParseTime("2024-01-26T10:55:00", false)
	c.currentTime = *newTime
}

func (c *ConverterTestSuite) TestStringToInt64Ptr() {
	tests := map[string]struct {
		input    string
		expected *int64
		fail     bool
	}{
		"happy path - Valid Input": {
			input:    "123",
			expected: Int64Ptr(123),
			fail:     false,
		},
		"should return an error while converting to int64 pointer": {
			input:    "abc",
			expected: nil,
			fail:     true,
		},
	}
	for name, tc := range tests {
		c.Run(name, func() {
			result, err := StringToInt64Ptr(tc.input)
			if tc.fail && err != nil {
				c.Nil(result)
			} else {
				c.NoError(err)
				c.Equal(*tc.expected, *result)
			}
		})
	}
}

func (c *ConverterTestSuite) TestStringToBoolPtr() {
	tests := map[string]struct {
		input    string
		expected *bool
		fail     bool
	}{
		"happy path - Valid Input": {
			input:    "true",
			expected: BoolPtr(true),
			fail:     false,
		},
		"should return an error while converting to bool pointer": {
			input:    "abc",
			expected: nil,
			fail:     true,
		},
	}
	for name, tc := range tests {
		c.Run(name, func() {
			result, err := StringToBoolPtr(tc.input)
			if tc.fail && err != nil {
				c.Nil(result)
			} else {
				c.NoError(err)
				c.Equal(*tc.expected, *result)
			}
		})
	}
}

func (c *ConverterTestSuite) TestCustomTimeToFormat() {
	tests := map[string]struct {
		input    *CustomTime
		format   *string
		expected *string
		fail     bool
	}{
		"happy path - Valid Input": {
			input:    &c.currentTime,
			format:   StringPtr(DefaultTimeFormat),
			expected: StringPtr("10:55:00Z"),
			fail:     false,
		},
		"should return an error while converting to invalid time format": {
			input:    &c.currentTime,
			format:   StringPtr("invalid"),
			expected: nil,
			fail:     true,
		},
		"should return an error while converting to specified time format": {
			input:    nil,
			format:   nil,
			expected: nil,
			fail:     true,
		},
	}
	for name, tc := range tests {
		c.Run(name, func() {
			result, err := CustomTimeToFormat(tc.input, tc.format)
			if tc.fail && err != nil {
				c.Nil(result)
			} else {
				c.NoError(err)
				c.Equal(*tc.expected, result.String())
			}
		})
	}
}

func (c *ConverterTestSuite) TestCustomDateToFormat() {
	tests := map[string]struct {
		input    *CustomDate
		format   *string
		expected *string
		fail     bool
	}{
		"happy path - Valid Input": {
			input:    &c.currentDate,
			format:   StringPtr(CustomDateFormat),
			expected: StringPtr("20240126Z"),
			fail:     false,
		},
		"should return an error while converting to invalid date format": {
			input:    &c.currentDate,
			format:   StringPtr("invalid"),
			expected: nil,
			fail:     true,
		},
		"should return an error while converting to specified date format": {
			input:    nil,
			format:   nil,
			expected: nil,
			fail:     true,
		},
	}
	for name, tc := range tests {
		c.Run(name, func() {
			result, err := CustomDateToFormat(tc.input, tc.format)
			if tc.fail && err != nil {
				c.Nil(result)
			} else {
				c.NoError(err)
				c.Equal(*tc.expected, result.String())
			}
		})
	}
}
