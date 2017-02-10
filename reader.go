package sunfish

import (
	"strconv"
	"reflect"
	"encoding/csv"
	"io"
	"os"
	"errors"
)

type Reader interface {
	ReadCsvFromFile(fileName string, ents interface{}) error
	ReadCsv(file io.Reader, ents interface{}) error
}

type ReaderImpl struct{}

func NewReader() Reader {
	c := ReaderImpl{}

	var reader Reader = &c
	return reader
}

func (r *ReaderImpl) ReadCsvFromFile(fileName string, ents interface{}) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	return r.ReadCsv(f, ents)
}

func (r *ReaderImpl) ReadCsv(file io.Reader, ents interface{}) error {
	reader := csv.NewReader(file)

	err := r.readData(reader, ents)

	return err
}

func (r *ReaderImpl) readData(reader *csv.Reader, ents interface{}) error {
	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}

	sliceValue := reflect.Indirect(reflect.ValueOf(ents))
	sliceElementType := sliceValue.Type()
	sliceElement := sliceElementType.Elem()

	newElemFunc := MakeNewElemFunction(sliceElement)
	sliceValueSetFunc := MakeSliceValueSetFunc(sliceElement, sliceValue)

	for _, line := range lines {
		newValue := newElemFunc().Elem()
		structIndex := 0
		for _, elem := range line {
			for ; structIndex < newValue.NumField() && sliceElement.Field(structIndex).Tag.Get("csv") != "parse"; structIndex++ {}
			if structIndex >= newValue.NumField() {
				break
			}

			field := newValue.Field(structIndex)
			err := r.setField(&field, elem)
			if err != nil {
				return err
			}
			structIndex++
		}
		sliceValueSetFunc(&newValue)
	}

	return err
}

func (r *ReaderImpl) setField(field *reflect.Value, elem string) error {
	if field.IsValid() && field.CanSet() {
		switch field.Kind() {
		case reflect.String:
			field.SetString(elem)
		case reflect.Bool:
			b, err := strconv.ParseBool(elem)
			if err != nil {
				return err
			}
			field.SetBool(b)
		case reflect.Int, reflect.Int32, reflect.Int64:
			ii, err := strconv.Atoi(elem)
			if err != nil {
				return err
			}
			field.SetInt(int64(ii))
		case reflect.Float64:
			f, err := strconv.ParseFloat(elem, 64)
			if err != nil {
				return err
			}
			field.SetFloat(f)
		case reflect.Float32:
			f, err := strconv.ParseFloat(elem, 32)
			if err != nil {
				return err
			}
			field.SetFloat(f)
		default:
			errors.New("Invalid type. CSVs can only be parsed into primitive types.")
		}
	} else {
		return errors.New("Field is either invalid or cannot be set.")
	}
	return nil
}