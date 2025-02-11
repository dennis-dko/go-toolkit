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
	DefaultTimeFormat      = "15:04:05Z"
	CustomTimeFormat       = "2006-01-02T15:04:05Z"
	DefaultTimeFormatNoUtc = "15:04:05"
	CustomTimeFormatNoUtc  = "2006-01-02T15:04:05"
)

var UtcTimeList = []string{
	CustomTimeFormat,
	DefaultTimeFormat,
}

var TimeList = []string{
	CustomTimeFormatNoUtc,
	DefaultTimeFormatNoUtc,
}

type CustomTime struct {
	time.Time
	NoUtc  bool
	Format string
}

func (ct *CustomTime) UnmarshalParam(param string) error {
	if param != "" {
		if IsUtcTime(param) {
			ct.NoUtc = false
		} else {
			ct.NoUtc = true
		}
		parsedTime, err := ParseTime(param, ct.NoUtc)
		if err != nil {
			return err
		}
		ct.Time = parsedTime.Time
		ct.Format = parsedTime.Format
	}
	return nil
}

func (ct CustomTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(ct.String(), start)
}

func (ct *CustomTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s *string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	if s != nil && *s != "" {
		if IsUtcTime(*s) {
			ct.NoUtc = false
		} else {
			ct.NoUtc = true
		}
		parsedTime, err := ParseTime(*s, ct.NoUtc)
		if err != nil {
			return err
		}
		ct.Time = parsedTime.Time
		ct.Format = parsedTime.Format
	}
	return nil
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.String())
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil && *s != "" {
		if IsUtcTime(*s) {
			ct.NoUtc = false
		} else {
			ct.NoUtc = true
		}
		parsedTime, err := ParseTime(*s, ct.NoUtc)
		if err != nil {
			return err
		}
		ct.Time = parsedTime.Time
		ct.Format = parsedTime.Format
	}
	return nil
}

func (ct *CustomTime) Scan(value interface{}) error {
	if data, ok := value.(time.Time); ok {
		ct.Time = data
		ct.NoUtc = false
		ct.Format = CustomTimeFormat
		return nil
	}
	if data, ok := value.(string); ok {
		if data != "" {
			if IsUtcTime(data) {
				ct.NoUtc = false
			} else {
				ct.NoUtc = true
			}
			parsedTime, err := ParseTime(data, ct.NoUtc)
			if err != nil {
				return err
			}
			ct.Time = parsedTime.Time
			ct.Format = parsedTime.Format
		}
		return nil
	}
	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", value, ct.Time)
}

func (ct CustomTime) Value() (driver.Value, error) {
	if ct.Time.IsZero() {
		return nil, nil
	}
	return ct.String(), nil
}

// String returns the time in the custom format
func (ct *CustomTime) String() string {
	return ct.Time.Format(ct.Format)
}

// SubTime returns the duration t-u. If the result exceeds the maximum (or minimum) value that can be stored in a Duration,
// the maximum (or minimum) duration will be returned. To compute t-d for a duration d, use t.Add(-d).
func (ct *CustomTime) SubTime(value *CustomTime) time.Duration {
	return ct.Time.Sub(value.Time)
}

// NewTime returns a new CustomTime object with actual utc time.
func NewTime(noTimezone bool) (*CustomTime, error) {
	receivedFormat := CustomTimeFormat
	if noTimezone {
		receivedFormat = CustomTimeFormatNoUtc
	}
	parsedTime, err := time.Parse(receivedFormat, time.Now().UTC().Format(receivedFormat))
	if err != nil {
		return nil, err
	}
	return &CustomTime{
		Time:   parsedTime,
		NoUtc:  noTimezone,
		Format: receivedFormat,
	}, nil
}

// ParseTime parses the given time string and returns a CustomTime object.
// It returns an error if the parsing fails.
func ParseTime(timeString string, noTimezone bool) (*CustomTime, error) {
	list := UtcTimeList
	if noTimezone {
		list = TimeList
	} else {
		timeString = SetAsUtc(timeString)
	}
	for _, format := range list {
		parsedTime, err := time.Parse(format, timeString)
		if err == nil {
			return &CustomTime{
				Time:   parsedTime,
				NoUtc:  noTimezone,
				Format: format,
			}, nil
		}
	}
	return nil, fmt.Errorf(`cannot parse the given time string "%s" to one of these formats (%s)`, timeString, strings.Join(list, ", "))
}
