package util

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type MapSuite struct {
	suite.Suite
}

func TestMapSuite(t *testing.T) {
	suite.Run(t, new(MapSuite))
}

func (m *MapSuite) TestStringifyMap() {
	tests := map[string]struct {
		data     interface{}
		expected string
	}{
		"happy path - stringify map": {
			data:     map[string]interface{}{"key": "value"},
			expected: "map[key:value]",
		},
		"happy path - stringify not a map": {
			data:     5,
			expected: "",
		},
	}
	for name, tc := range tests {
		m.Run(name, func() {
			result := StringifyMap(tc.data)
			m.Equal(tc.expected, result)
		})
	}
}
