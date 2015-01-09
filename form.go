package form

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type field struct {
	name     string
	setter   setter
	required bool
	defval   interface{}
}

// Decoder of form structs
type Decoder struct {
	typ    reflect.Type
	fields map[string]field
}

type setter func(reflect.Value, interface{}) error

// NewDecoder creates decoder for given form
func NewDecoder(v interface{}) *Decoder {
	if v == nil {
		panic("null argument")
	}
	var d = &Decoder{typ: typeOf(v)}
	d.reflectFields(d.typ)
	return d
}

// Decode form struct from given input
func (d *Decoder) Decode(form interface{}, input interface{}) error {
	switch input.(type) {
	case *http.Request:
		return d.decodeRequest(form, input.(*http.Request))
	case map[string]interface{}:
		return d.decodeMap(form, input.(map[string]interface{}))
	case map[string]string:
		var m = input.(map[string]string)
		var data = make(map[string]interface{})
		for key, val := range m {
			data[key] = val
		}
		return d.decodeMap(form, data)
	case url.Values:
		return d.decodeForm(form, input.(url.Values))
	case map[string][]string:
		return d.decodeForm(form, input.(map[string][]string))
	default:
		return errors.New("invalid input")
	}
}

func (d *Decoder) decodeRequest(form interface{}, r *http.Request) error {
	if r.Method == "GET" {
		return d.decodeForm(form, r.URL.Query())
	}

	var contentType = r.Header.Get("Content-Type")

	// support JSON
	if contentType == "application/json" {
		var data = make(map[string]interface{})
		var dec = json.NewDecoder(r.Body)
		var err = dec.Decode(&data)
		if err != nil {
			return err
		}
		return d.decodeMap(form, data)
	}

	// support XML
	if contentType == "text/xml" || contentType == "application/xml" {
		var data = make(map[string]interface{})
		var dec = xml.NewDecoder(r.Body)
		var err = dec.Decode(&data)
		if err != nil {
			return err
		}
		return d.decodeMap(form, data)
	}

	var err = r.ParseForm()
	if err != nil {
		return err
	}

	return d.decodeForm(form, r.Form)
}

func (d *Decoder) decodeForm(form interface{}, data url.Values) error {
	var f = make(map[string]interface{})
	for key := range data {
		f[key] = data.Get(key)
	}
	return d.decodeMap(form, f)
}

func (d *Decoder) decodeMap(form interface{}, data map[string]interface{}) error {
	value := reflect.ValueOf(form).Elem()

	required := make(map[string]bool)
	for k, f := range d.fields {
		if f.required {
			required[k] = true
		}
	}

	var errs = MultiError{}

	for k, v := range data {
		var key = strings.ToLower(k)

		f, ok := d.fields[key]
		if !ok {
			continue
		}

		err := f.setter(value, v)
		if err != nil {
			errs[key] = err
		}

		if _, ok = required[key]; ok && !isNullOrEmpty(v) {
			delete(required, key)
		}
	}

	// TODO set default values

	if len(errs) > 0 {
		return errs
	}

	if len(required) > 0 {
		var names = []string{}
		for k := range required {
			names = append(names, k)
		}
		return fmt.Errorf("Please fill required fields: %s", strings.Join(names, ", "))
	}

	return nil
}

func (d *Decoder) reflectFields(t reflect.Type) {
	d.fields = make(map[string]field)
	for i := 0; i < t.NumField(); i++ {
		var f = t.Field(i)
		var name = strings.ToLower(f.Name)

		var json = f.Tag.Get("json")
		if len(json) > 0 {
			var spec = strings.Split(json, ",")
			name = strings.ToLower(spec[0])
		}

		var required = false
		if req, err := strconv.ParseBool(f.Tag.Get("required")); err == nil {
			required = req
		}

		d.fields[name] = field{
			name:     name,
			setter:   fieldSetter(f),
			required: required,
		}
	}
}

func fieldSetter(field reflect.StructField) setter {
	return func(instance reflect.Value, v interface{}) error {
		cnv, ok := converters[reflect.TypeOf(v)]
		if !ok {
			return errors.New("unsupported value type")
		}

		var val = reflect.ValueOf(cnv(v))

		instance.FieldByName(field.Name).Set(val)

		return nil
	}
}

func typeOf(v interface{}) reflect.Type {
	var t reflect.Type

	switch v.(type) {
	case reflect.Type:
		t = v.(reflect.Type)
		break
	default:
		t = reflect.TypeOf(v)
		break
	}

	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}

	return t
}

func isNullOrEmpty(v interface{}) bool {
	switch v.(type) {
	case string:
		s := v.(string)
		return len(strings.TrimSpace(s)) == 0
	default:
		return v == nil
	}
}
