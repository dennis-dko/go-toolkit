package datatype

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/dennis-dko/go-toolkit/constant"

	"github.com/antchfx/xmlquery"
)

// Only the lower version of the value should be set
var skipValues = []string{
	"t",
	"f",
}

// ParseXMLToStruct parses xml data into given struct
// by using defined field tags like XPath e.g.
// Value: get (nested) value use nxml:"//|xmltag|" or nxml:"//|xmltag|/|xmltag|"
// Attribute: get (nested) attribute use nxml:"//|xmltag|/@|attrkey|" or nxml:"//|xmltag|/|xmltag|/@|attrkey|"
// It will be use the struct field names for json keys as default, change it by adding a json tag like this: json:"|jsontag|"
// It's important to define a struct field named "XMLName" with the outer xml node name e.g.
// XMLName string `json:"-" nxml:"//tests/test"` for "<tests><test>value</test></tests>"
// The other fields will be used for the inner xml nodes or attributes e.g.
// FieldName *string `json:"fieldname,omitempty" nxml:"//test/@attrkey"`
func ParseXMLToStruct(xmlData string, rawStruct interface{}) error {
	var (
		isSlice     bool
		marshalData interface{}
	)
	structData := reflect.ValueOf(rawStruct)
	if structData.Kind() == reflect.Ptr {
		structData = reflect.ValueOf(rawStruct).Elem()
	}
	if structData.Kind() == reflect.Slice {
		isSlice = true
		structData = reflect.New(structData.Type().Elem()).Elem()
	}
	if structData.Kind() != reflect.Struct {
		return errors.New("invalid interface type or no data found")
	}
	if structData.FieldByName(constant.XMLField) == (reflect.Value{}) {
		return fmt.Errorf(`missing field "%s" in struct`, constant.XMLField)
	}
	data := make(map[string]string)
	for i := 0; i < structData.NumField(); i++ {
		nxmlSelector := strings.SplitN(structData.Type().Field(i).Tag.Get("nxml"), ",", 2)[0]
		if !strings.HasPrefix(nxmlSelector, "//") {
			continue
		}
		structField := strings.SplitN(structData.Type().Field(i).Tag.Get("json"), ",", 2)[0]
		if structField == "" || structField == "-" {
			structField = structData.Type().Field(i).Name
		}
		data[structField] = nxmlSelector
	}
	if len(data) > 0 {
		xmlDoc, err := xmlquery.Parse(
			strings.NewReader(xmlData),
		)
		if err != nil {
			return err
		}
		nxmlList := make([]map[string]interface{}, 0)
		for field, nxml := range data {
			if field == constant.XMLField {
				continue
			}
			for index, parentNode := range xmlquery.Find(xmlDoc, data[constant.XMLField]) {
				childNode, err := xmlquery.Parse(
					strings.NewReader(
						parentNode.OutputXML(true),
					),
				)
				if err != nil {
					return err
				}
				var (
					value      string
					valueSlice []string
				)
				for _, content := range xmlquery.Find(childNode, nxml) {
					valueSlice = append(valueSlice, content.InnerText())
				}
				if valueSlice != nil {
					value = strings.Join(valueSlice, ",")
					if len(nxmlList) > index {
						nxmlList[index][field] = setDataType(value)
						continue
					}
					nxmlData := make(map[string]interface{})
					nxmlData[field] = setDataType(value)
					nxmlList = append(nxmlList, nxmlData)
				}
			}
		}
		if isSlice {
			marshalData = nxmlList
		} else {
			marshalData = nxmlList[0]
		}
		jsonData, err := json.Marshal(marshalData)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(jsonData, &rawStruct); err != nil {
			return err
		}
	}
	return nil
}

// GetXMLValue get the first value from xml data which found by given selector
func GetXMLValue(xmlData string, selector string) (string, error) {
	xmlDoc, err := xmlquery.Parse(strings.NewReader(xmlData))
	if err != nil {
		return "", err
	}
	node := xmlquery.FindOne(xmlDoc, selector)
	if node == nil {
		return "", errors.New("cannot find node by given selector")
	}
	return node.InnerText(), nil
}

func setDataType(value string) interface{} {
	if slices.Contains(skipValues, strings.ToLower(value)) {
		return value
	}
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				return value
			}
			return boolValue
		}
		return floatValue
	}
	return intValue
}
