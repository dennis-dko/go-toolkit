package datatype

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"time"
)

// NullBool nullable bool that overrides sql.NullBool
type NullBool struct {
	sql.NullBool
}

// NewNullBool returns a new NullBool with the given bool value.
func NewNullBool(value *bool) NullBool {
	data := NullBool{
		sql.NullBool{
			Bool:  false,
			Valid: false,
		},
	}
	if value != nil {
		data.Bool = *value
		data.Valid = true
	}
	return data
}

// NullBoolPtrToBoolPtr converts a *NullBool to *bool
func NullBoolPtrToBoolPtr(nb *NullBool) *bool {
	if nb != nil && nb.Valid {
		return &nb.Bool
	}
	return nil
}

func (nb NullBool) MarshalJSON() ([]byte, error) {
	if nb.Valid {
		return json.Marshal(nb.Bool)
	}
	return json.Marshal(nil)
}

func (nb *NullBool) UnmarshalJSON(data []byte) error {
	var b *bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	if b != nil {
		nb.Valid = true
		nb.Bool = *b
	} else {
		nb.Valid = false
	}
	return nil
}

func (nb NullBool) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if nb.Valid {
		return e.EncodeElement(nb.Bool, start)
	}
	return e.EncodeElement(nil, start)
}

func (nb *NullBool) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var b *bool
	if err := d.DecodeElement(&b, &start); err != nil {
		return err
	}
	if b != nil {
		nb.Valid = true
		nb.Bool = *b
	} else {
		nb.Valid = false
	}
	return nil
}

// NullFloat64 nullable float64 that overrides sql.NullFloat64
type NullFloat64 struct {
	sql.NullFloat64
}

// NewNullFloat64 creates a NullFloat64
func NewNullFloat64(value *float64) NullFloat64 {
	data := NullFloat64{
		sql.NullFloat64{
			Float64: 0,
			Valid:   false,
		},
	}
	if value != nil {
		data.Float64 = *value
		data.Valid = true
	}
	return data
}

// NullFloat64PtrToFloat64Ptr converts a *NullFloat64 to *float64
func NullFloat64PtrToFloat64Ptr(nf *NullFloat64) *float64 {
	if nf != nil && nf.Valid {
		return &nf.Float64
	}
	return nil
}

func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if nf.Valid {
		return json.Marshal(nf.Float64)
	}
	return json.Marshal(nil)
}

func (nf *NullFloat64) UnmarshalJSON(data []byte) error {
	var f *float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	if f != nil {
		nf.Valid = true
		nf.Float64 = *f
	} else {
		nf.Valid = false
	}
	return nil
}

func (nf NullFloat64) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if nf.Valid {
		return e.EncodeElement(nf.Float64, start)
	}
	return e.EncodeElement(nil, start)
}

func (nf *NullFloat64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var f *float64
	if err := d.DecodeElement(&f, &start); err != nil {
		return err
	}
	if f != nil {
		nf.Valid = true
		nf.Float64 = *f
	} else {
		nf.Valid = false
	}
	return nil
}

// NullInt64 nullable int64 that overrides sql.NullInt64
type NullInt64 struct {
	sql.NullInt64
}

// NewNullInt64 creates a NullInt64
func NewNullInt64(value *int64) NullInt64 {
	data := NullInt64{
		sql.NullInt64{
			Int64: 0,
			Valid: false,
		},
	}
	if value != nil {
		data.Int64 = *value
		data.Valid = true
	}
	return data
}

// NullInt64PtrToInt64Ptr converts a *NullInt64 to *int64
func NullInt64PtrToInt64Ptr(ni *NullInt64) *int64 {
	if ni != nil && ni.Valid {
		return &ni.Int64
	}
	return nil
}

func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if ni.Valid {
		return json.Marshal(ni.Int64)
	}
	return json.Marshal(nil)
}

func (ni *NullInt64) UnmarshalJSON(data []byte) error {
	var i *int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		ni.Valid = true
		ni.Int64 = *i
	} else {
		ni.Valid = false
	}
	return nil
}

func (ni NullInt64) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if ni.Valid {
		return e.EncodeElement(ni.Int64, start)
	}
	return e.EncodeElement(nil, start)
}

func (ni *NullInt64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var i *int64
	if err := d.DecodeElement(&i, &start); err != nil {
		return err
	}
	if i != nil {
		ni.Valid = true
		ni.Int64 = *i
	} else {
		ni.Valid = false
	}
	return nil
}

// NullString nullable string that overrides sql.NullString
type NullString struct {
	sql.NullString
}

// NewNullString creates a NullString
func NewNullString(value *string) NullString {
	data := NullString{
		sql.NullString{
			String: "",
			Valid:  false,
		},
	}
	if value != nil {
		data.String = *value
		data.Valid = true
	}
	return data
}

// NullStringPtrToStringPtr converts a *NullString to *string
func NullStringPtrToStringPtr(ns *NullString) *string {
	if ns != nil && ns.Valid {
		return &ns.String
	}
	return nil
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}

func (ns NullString) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if ns.Valid {
		return e.EncodeElement(ns.String, start)
	}
	return e.EncodeElement(nil, start)
}

func (ns *NullString) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s *string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}

// NullTime nullable time that overrides sql.NullTime
type NullTime struct {
	sql.NullTime
	NoUtc  bool
	Format string
}

// NewNullTime creates a NullTime
func NewNullTime(value *CustomTime) NullTime {
	data := NullTime{
		NullTime: sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
		NoUtc:  false,
		Format: "",
	}
	if value != nil {
		data.Time = value.Time
		data.Valid = true
		data.NoUtc = value.NoUtc
		data.Format = value.Format
	}
	return data
}

// NullTimePtrToCustomTimePtr converts a *NullTime to *CustomTime
func NullTimePtrToCustomTimePtr(nd *NullTime) *CustomTime {
	if nd != nil && nd.Valid {
		return &CustomTime{
			Time:   nd.Time,
			NoUtc:  nd.NoUtc,
			Format: nd.Format,
		}
	}
	return nil
}

// String returns the time in the custom format
func (nt *NullTime) String() string {
	return nt.Time.Format(nt.Format)
}

func (nt *NullTime) UnmarshalParam(param string) error {
	ct := new(CustomTime)
	err := ct.UnmarshalParam(param)
	if err != nil {
		return err
	}
	if ct != nil {
		nt.Valid = true
		nt.Time = ct.Time
		nt.NoUtc = ct.NoUtc
		nt.Format = ct.Format
	} else {
		nt.Valid = false
	}
	return nil
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return json.Marshal(nt.String())
	}
	return json.Marshal(nil)
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	var ct *CustomTime
	if err := json.Unmarshal(data, &ct); err != nil {
		return err
	}
	if ct != nil {
		nt.Valid = true
		nt.Time = ct.Time
		nt.NoUtc = ct.NoUtc
		nt.Format = ct.Format
	} else {
		nt.Valid = false
	}
	return nil
}

func (nt NullTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if nt.Valid {
		return e.EncodeElement(nt.String(), start)
	}
	return e.EncodeElement(nil, start)
}

func (nt *NullTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var ct *CustomTime
	if err := d.DecodeElement(&ct, &start); err != nil {
		return err
	}
	if ct != nil {
		nt.Valid = true
		nt.Time = ct.Time
		nt.NoUtc = ct.NoUtc
		nt.Format = ct.Format
	} else {
		nt.Valid = false
	}
	return nil
}

func (nt *NullTime) Scan(value interface{}) error {
	ct := new(CustomTime)
	err := ct.Scan(value)
	if err != nil {
		return err
	}
	if ct != nil {
		nt.Valid = true
		nt.Time = ct.Time
		nt.NoUtc = ct.NoUtc
		nt.Format = ct.Format
	} else {
		nt.Valid = false
	}
	return nil
}

// NullDate nullable date that overrides sql.NullTime
type NullDate struct {
	sql.NullTime
	NoUtc  bool
	Format string
}

// NewNullDate creates a NullDate
func NewNullDate(value *CustomDate) NullDate {
	data := NullDate{
		NullTime: sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
		NoUtc:  false,
		Format: "",
	}
	if value != nil {
		data.Time = value.Time
		data.Valid = true
		data.NoUtc = value.NoUtc
		data.Format = value.Format
	}
	return data
}

// NullDatePtrToCustomDatePtr converts a *NullDate to *CustomDate
func NullDatePtrToCustomDatePtr(nd *NullDate) *CustomDate {
	if nd != nil && nd.Valid {
		return &CustomDate{
			Time:   nd.Time,
			NoUtc:  nd.NoUtc,
			Format: nd.Format,
		}
	}
	return nil
}

// String returns the date in the custom format
func (nd *NullDate) String() string {
	return nd.Time.Format(nd.Format)
}

func (nd *NullDate) UnmarshalParam(param string) error {
	ct := new(CustomDate)
	err := ct.UnmarshalParam(param)
	if err != nil {
		return err
	}
	if ct != nil {
		nd.Valid = true
		nd.Time = ct.Time
		nd.NoUtc = ct.NoUtc
		nd.Format = ct.Format
	} else {
		nd.Valid = false
	}
	return nil
}

func (nd NullDate) MarshalJSON() ([]byte, error) {
	if nd.Valid {
		return json.Marshal(nd.String())
	}
	return json.Marshal(nil)
}

func (nd *NullDate) UnmarshalJSON(data []byte) error {
	var cd *CustomDate
	if err := json.Unmarshal(data, &cd); err != nil {
		return err
	}
	if cd != nil {
		nd.Valid = true
		nd.Time = cd.Time
		nd.NoUtc = cd.NoUtc
		nd.Format = cd.Format
	} else {
		nd.Valid = false
	}
	return nil
}

func (nd NullDate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if nd.Valid {
		return e.EncodeElement(nd.String(), start)
	}
	return e.EncodeElement(nil, start)
}

func (nd *NullDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var cd *CustomDate
	if err := d.DecodeElement(&cd, &start); err != nil {
		return err
	}
	if cd != nil {
		nd.Valid = true
		nd.Time = cd.Time
		nd.NoUtc = cd.NoUtc
		nd.Format = cd.Format
	} else {
		nd.Valid = false
	}
	return nil
}

func (nd *NullDate) Scan(value interface{}) error {
	ct := new(CustomDate)
	err := ct.Scan(value)
	if err != nil {
		return err
	}
	if ct != nil {
		nd.Valid = true
		nd.Time = ct.Time
		nd.NoUtc = ct.NoUtc
		nd.Format = ct.Format
	} else {
		nd.Valid = false
	}
	return nil
}
