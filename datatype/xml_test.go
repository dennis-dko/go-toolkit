package datatype

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ParseStruct struct {
	XMLName    string `json:"-" nxml:"//users/user"`
	FirstName  string `nxml:"//user/@firstName" json:"first_name"`
	LastName   string `nxml:"//user/@lastName" json:"last_name"`
	Age        uint8  `nxml:"//user/@age" json:"age"`
	Gender     string `nxml:"//user/@gender" json:"gender"`
	Company    string `nxml:"//user/@company" json:"company"`
	Email      string `nxml:"//user/postal/street/email" json:"email"`
	Pet        string `nxml:"//user/postal/street/animal/@pet" json:"pet"`
	Street     string `nxml:"//user/postal/street/@address" json:"street"`
	PostalCode int    `nxml:"//user/postal/@code" json:"postal_code"`
}

type XMLTestSuite struct {
	suite.Suite
	expectedStructList []ParseStruct
	xml                []byte
}

func (x *XMLTestSuite) SetupTest() {
	// Setup
	x.expectedStructList = []ParseStruct{
		{
			FirstName:  "Walter",
			LastName:   "White",
			Age:        50,
			Gender:     "male",
			Company:    "T",
			Email:      "walter.white@example.com",
			Pet:        "cat",
			Street:     "Villa Gaeta",
			PostalCode: 91764,
		},
		{
			FirstName:  "James",
			LastName:   "McGill",
			Age:        45,
			Gender:     "male",
			Company:    "S",
			Email:      "james.mcgill@example.com",
			Pet:        "dog",
			Street:     "Saul Street",
			PostalCode: 65782,
		},
	}
	x.xml = []byte(`<users count="2"><user firstName="Walter" lastName="White" age="50" gender="male" company="T"><postal code="91764"><street address="Villa Gaeta"><email>walter.white@example.com</email><animal pet="cat"></animal></street></postal></user><user firstName="James" lastName="McGill" age="45" gender="male" company="S"><postal code="65782"><street address="Saul Street"><email>james.mcgill@example.com</email><animal pet="dog"></animal></street></postal></user></users>`)
}

func TestXMLTestSuite(t *testing.T) {
	suite.Run(t, new(XMLTestSuite))
}

func (x *XMLTestSuite) TestParseXMLToStruct() {

	x.Run("happy path - all data in xml could be parsed to struct", func() {
		// Init
		var parseStructList []ParseStruct

		// Run
		parseErr := ParseXMLToStruct(string(x.xml), &parseStructList)

		// Assert
		x.NoError(parseErr)
		x.Equal(x.expectedStructList, parseStructList)
	})
}

func (x *XMLTestSuite) TestGetXMLValue() {

	x.Run("happy path - get correct xml value", func() {
		// Run
		rawValue, getErr := GetXMLValue(string(x.xml), "//users/@count")
		convValue, convErr := strconv.Atoi(rawValue)

		// Assert
		x.NoError(getErr)
		x.NoError(convErr)
		x.Equal(2, convValue)
	})
}
