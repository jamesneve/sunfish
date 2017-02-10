package parser

import (
	"strconv"
	"reflect"
	"encoding/csv"
	"io"
	"os"
	"errors"
)

type Parser interface {
	ReadCsvFromFileInOrder(fileName string, ents interface{}) error
	ReadCsvFromFileWithHeaders(fileName string, ents interface{}) error
	ReadCsvInOrder(file io.Reader, ents interface{}) error
	ReadCsvWithHeaders(file io.Reader, ents interface{}) error
}

type ParserImpl struct{}

func NewParser() Parser {
	c := ParserImpl{}

	var reader Parser = &c
	return reader
}

func (r *ParserImpl) ReadCsvFromFileInOrder(fileName string, ents interface{}) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	return r.ReadCsvInOrder(f, ents)
}

func (r *ParserImpl) ReadCsvFromFileWithHeaders(fileName string, ents interface{}) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	return r.ReadCsvWithHeaders(f, ents)
}

func (r *ParserImpl) ReadCsvInOrder(file io.Reader, ents interface{}) error {
	reader := csv.NewReader(file)

	err := r.readDataInOrder(reader, ents)

	return err
}

func (r *ParserImpl) ReadCsvWithHeaders(file io.Reader, ents interface{}) error {
	reader := csv.NewReader(file)

	err := r.readDataWithNames(reader, ents)

	return err
}

func (r *ParserImpl) readDataInOrder(reader *csv.Reader, ents interface{}) error {
	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, line := range lines {
		r.readLine(line, ents)
	}

	return err
}

func (r *ParserImpl) readDataWithNames(reader *csv.Reader, ents interface{}) error {
	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(lines) <= 1 {
		return errors.New("No data")
	}
	headers := lines[0]

	for i, line := range lines {
		if i == 0 { continue }
		r.readLineWithNames(line, ents, headers)
	}

	return err
}

func (r *ParserImpl) readLine(line []string, ents interface{}) error {
	sliceValue := reflect.Indirect(reflect.ValueOf(ents))

	sliceElement := sliceValue.Type().Elem()

	newValue := makeNewElemFunction(sliceElement)().Elem()
	structIndex := 0
	for _, elem := range line {
		for ; structIndex < newValue.NumField() && sliceElement.Field(structIndex).Tag.Get("csv") != "parse"; structIndex++ {}
		if structIndex >= newValue.NumField() {
			break
		}

		err := r.makeStructField(&newValue, elem, structIndex)
		if err != nil {
			return err
		}

		structIndex++
	}

	makeSliceValueSetFunc(sliceElement, sliceValue)(&newValue)
	return nil
}

func (r *ParserImpl) readLineWithNames(line []string, ents interface{}, headers []string) error {
	sliceValue := reflect.Indirect(reflect.ValueOf(ents))

	sliceElement := sliceValue.Type().Elem()

	newValue := makeNewElemFunction(sliceElement)().Elem()
	for i, header := range headers {
		for j := 0; j < newValue.NumField(); j++ {
			if sliceElement.Field(j).Tag.Get("csv") == header {
				err := r.makeStructField(&newValue, line[i], j)
				if err != nil { return err }
			}
		}
	}

	makeSliceValueSetFunc(sliceElement, sliceValue)(&newValue)
	return nil
}

func (r *ParserImpl) makeStructField(value *reflect.Value, elem string, structIndex int) error {
	field := value.Field(structIndex)
	err := r.setField(&field, elem)

	return err
}

func (r *ParserImpl) setField(field *reflect.Value, elem string) error {
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