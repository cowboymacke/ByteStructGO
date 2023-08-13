package bytestruct

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

func Unmarshal(reader io.Reader, order binary.ByteOrder, v interface{}) error {

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	storedValues := make(map[string]reflect.Value)

	readData(reader, order, reflect.StructField{}, val, storedValues)

	return nil
}

func readData(reader io.Reader, order binary.ByteOrder, structField reflect.StructField, val reflect.Value, storedValues map[string]reflect.Value) error {

	if val.Kind() != reflect.Struct {
		storedValues[structField.Name] = val
	}

	switch val.Kind() {

	case reflect.Struct:
		t := val.Type()
		for i := 0; i < val.NumField(); i++ {
			structF := t.Field(i)
			if v := val.Field(i); v.CanSet() || structF.Name != "_" {
				readData(reader, order, structF, v, storedValues)
			}
		}

	case reflect.Bool:
		var value bool
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetBool(value)
		} else {
			return err
		}
	case reflect.Int:
		var value int
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetInt(int64(value))
		} else {
			return err
		}
	case reflect.Int8:
		var value int8
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetInt(int64(value))
		} else {
			return err
		}
	case reflect.Int16:
		var value int16
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetInt(int64(value))
		} else {
			return err
		}
	case reflect.Int32:
		var value int32
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetInt(int64(value))
		} else {
			return err
		}
	case reflect.Int64:
		var value int64
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetInt(int64(value))
		} else {
			return err
		}
	case reflect.Uint:
		var value uint
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetUint(uint64(value))
		} else {
			return err
		}
	case reflect.Uint8:
		var value uint8
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetUint(uint64(value))
		} else {
			return err
		}
	case reflect.Uint16:
		var value uint16
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetUint(uint64(value))
		} else {
			return err
		}
	case reflect.Uint32:
		var value uint32
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetUint(uint64(value))
		} else {
			return err
		}
	case reflect.Uint64:
		var value uint64
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetUint(uint64(value))
		} else {
			return err
		}
	case reflect.Float32:
		var value float32
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetFloat(float64(value))
		} else {
			return err
		}
	case reflect.Float64:
		var value float64
		if err := binary.Read(reader, order, &value); err == nil {
			val.SetFloat(float64(value))
		} else {
			return err
		}
	case reflect.String:

		if err := unmarshalString(reader, order, storedValues, structField, val); err != nil {
			return err
		}
	case reflect.Slice:
		if err := unmarshalArray(reader, order, storedValues, structField, val); err != nil {
			return err
		}
	}

	return nil
}

func unmarshalArray(reader io.Reader, order binary.ByteOrder, storedValues map[string]reflect.Value, field reflect.StructField, value reflect.Value) error {

	if v, ok := field.Tag.Lookup("byte"); ok {

		size, err := strconv.Atoi(v)

		if err != nil {
			if mapValue, ok := storedValues[v]; ok {
				size = int(mapValue.Uint())
			}

		}

		data := make([]byte, size)

		switch field.Type.Elem().Kind() {

		case reflect.Uint8:
			if err := binary.Read(reader, order, &data); err == nil {
				value.SetBytes(data)
			}

		default:
			return fmt.Errorf("does not support type with array: %s ", field.Type.Elem().Kind())
		}
	} else {
		return fmt.Errorf("Field: %s is missing byte tag", field.Name)
	}

	return nil
}

func unmarshalString(reader io.Reader, order binary.ByteOrder, storedValues map[string]reflect.Value, field reflect.StructField, value reflect.Value) error {
	if v, ok := field.Tag.Lookup("byte"); ok {

		size, err := strconv.Atoi(v)

		if err != nil {
			if mapValue, ok := storedValues[v]; ok {
				size = int(mapValue.Uint())
			}

		}

		data := make([]byte, size)
		if err := binary.Read(reader, order, &data); err == nil {
			value.SetString(string(data))
		}
	} else {
		return fmt.Errorf("missing byte tag")
	}

	return nil
}

func Marshal(order binary.ByteOrder, v interface{}) ([]byte, error) {

	val := reflect.ValueOf(v)
	ty := reflect.TypeOf(v)

	m := make(map[string]reflect.Value)

	var buf bytes.Buffer

	for i := 0; i < val.NumField(); i++ {

		f := val.Field(i)
		t := ty.Field(i)

		if v, ok := t.Tag.Lookup("byteLength"); ok {
			length := uint8(len(val.FieldByName(v).String()))
			if err := binary.Write(&buf, order, length); err != nil {
				return nil, err
			}
		} else {
			m[t.Name] = f
			switch t.Type.Kind() {
			case reflect.Bool:
				if err := binary.Write(&buf, order, f.Bool()); err != nil {
					return nil, err
				}
			case reflect.Int:
				if err := binary.Write(&buf, order, f.Int()); err != nil {
					return nil, err
				}
			case reflect.Int8:
				if err := binary.Write(&buf, order, int8(f.Int())); err != nil {
					return nil, err
				}
			case reflect.Int16:
				if err := binary.Write(&buf, order, int16(f.Int())); err != nil {
					return nil, err
				}
			case reflect.Int32:
				if err := binary.Write(&buf, order, int32(f.Int())); err != nil {
					return nil, err
				}
			case reflect.Int64:
				if err := binary.Write(&buf, order, int64(f.Int())); err != nil {
					return nil, err
				}
			case reflect.Uint:
				if err := binary.Write(&buf, order, f.Uint()); err != nil {
					return nil, err
				}
			case reflect.Uint8:
				if err := binary.Write(&buf, order, uint8(f.Uint())); err != nil {
					return nil, err
				}
			case reflect.Uint16:
				if err := binary.Write(&buf, order, uint16(f.Uint())); err != nil {
					return nil, err
				}
			case reflect.Uint32:
				if err := binary.Write(&buf, order, uint32(f.Uint())); err != nil {
					return nil, err
				}
			case reflect.Uint64:
				if err := binary.Write(&buf, order, uint64(f.Uint())); err != nil {
					return nil, err
				}
			case reflect.Float32:
				if err := binary.Write(&buf, order, float32(f.Float())); err != nil {
					return nil, err
				}
			case reflect.Float64:
				if err := binary.Write(&buf, order, float64(f.Float())); err != nil {
					return nil, err
				}
			case reflect.String:
				if err := binary.Write(&buf, order, []byte(f.String())); err != nil {
					return nil, err
				}
			case reflect.Slice:
				if err := binary.Write(&buf, order, f.Bytes()); err != nil {
					return nil, err
				}
			case reflect.Struct:
				data, err := Marshal(order, f.Interface())
				if err != nil {
					return nil, err
				}

				if err := binary.Write(&buf, order, data); err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("does not support type : %s ", t.Type.Kind())
			}
		}

	}

	return buf.Bytes(), nil
}
