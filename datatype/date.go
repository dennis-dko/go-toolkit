package datatype

import (
	"database/sql/driver"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

const (
	CustomDateFormat       = "20060102Z"
	DefaultDateFormat      = "2006-01-02Z"
	CustomDateFormatNoUtc  = "20060102"
	DefaultDateFormatNoUtc = "2006-01-02"
)

var UtcDateList = []string{
	DefaultDateFormat,
	CustomDateFormat,
}

var DateList = []string{
	DefaultDateFormatNoUtc,
	CustomDateFormatNoUtc,
}

type CustomDate struct {
	time.Time
	NoUtc  bool
	Format string
}

func (cd *CustomDate) UnmarshalParam(param string) error {
	if param != "" {
		if IsUtcTime(param) {
			cd.NoUtc = false
		} else {
			cd.NoUtc = true
		}
		parsedDate, err := ParseDate(param, cd.NoUtc)
		if err != nil {
			return err
		}
		cd.Time = parsedDate.Time
		cd.Format = parsedDate.Format
	}
	return nil
}

func (cd CustomDate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(cd.String(), start)
}

func (cd *CustomDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s *string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	if s != nil && *s != "" {
		if IsUtcTime(*s) {
			cd.NoUtc = false
		} else {
			cd.NoUtc = true
		}
		parsedDate, err := ParseDate(*s, cd.NoUtc)
		if err != nil {
			return err
		}
		cd.Time = parsedDate.Time
		cd.Format = parsedDate.Format
	}
	return nil
}

func (cd CustomDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(cd.String())
}

func (cd *CustomDate) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil && *s != "" {
		if IsUtcTime(*s) {
			cd.NoUtc = false
		} else {
			cd.NoUtc = true
		}
		parsedDate, err := ParseDate(*s, cd.NoUtc)
		if err != nil {
			return err
		}
		cd.Time = parsedDate.Time
		cd.Format = parsedDate.Format
	}
	return nil
}

func (cd *CustomDate) Scan(value interface{}) error {
	if data, ok := value.(time.Time); ok {
		cd.Time = data
		cd.NoUtc = false
		cd.Format = DefaultDateFormat
		return nil
	}
	if data, ok := value.(string); ok {
		if data != "" {
			if IsUtcTime(data) {
				cd.NoUtc = false
			} else {
				cd.NoUtc = true
			}
			parsedDate, err := ParseDate(data, cd.NoUtc)
			if err != nil {
				return err
			}
			cd.Time = parsedDate.Time
			cd.Format = parsedDate.Format
		}
		return nil
	}
	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", value, cd.Time)
}

func (cd CustomDate) Value() (driver.Value, error) {
	if cd.Time.IsZero() {
		return nil, nil
	}
	return cd.String(), nil
}

// String returns the date in the custom format
func (cd *CustomDate) String() string {
	return cd.Time.Format(cd.Format)
}

// SubDate returns the duration t-u. If the result exceeds the maximum (or minimum) value that can be stored in a Duration,
// the maximum (or minimum) duration will be returned. To compute t-d for a duration d, use t.Add(-d).
func (cd *CustomDate) SubDate(value *CustomDate) time.Duration {
	return cd.Time.Sub(value.Time)
}

// NewDate returns a new CustomDate object with actual utc date.
func NewDate(noTimezone bool) (*CustomDate, error) {
	receivedFormat := DefaultDateFormat
	if noTimezone {
		receivedFormat = DefaultDateFormatNoUtc
	}
	parsedDate, err := time.Parse(receivedFormat, time.Now().UTC().Format(receivedFormat))
	if err != nil {
		return nil, err
	}
	return &CustomDate{
		Time:   parsedDate,
		NoUtc:  noTimezone,
		Format: receivedFormat,
	}, nil
}

// ParseDate parses the given date string and returns a CustomDate object.
// It returns an error if the parsing fails.
func ParseDate(dateString string, noTimezone bool) (*CustomDate, error) {
	if dateString == "" {
		return nil, nil
	}
	list := UtcDateList
	if noTimezone {
		list = DateList
	} else {
		dateString = SetAsUtc(dateString)
	}
	for _, format := range list {
		parsedDate, err := time.Parse(format, dateString)
		if err == nil {
			return &CustomDate{
				Time:   parsedDate,
				NoUtc:  noTimezone,
				Format: format,
			}, nil
		}
	}
	return nil, fmt.Errorf(`cannot parse the given date string "%s" to one of these formats (%s)`, dateString, strings.Join(list, ", "))
}
