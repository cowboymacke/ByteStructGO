package bytestruct

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

func Unmarshal(reader io.Reader, order binary.ByteOrder, v interface{}) error {

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	storedValues := make(map[string]reflect.Value)

	if err := readData(reader, order, reflect.StructField{}, val, storedValues); err != nil {
		return err
	}

	return nil
}

func readData(reader io.Reader, order binary.ByteOrder, structField reflect.StructField, val reflect.Value, storedValues map[string]reflect.Value) error {

	if val.Kind() != reflect.Struct {
		storedValues[structField.Name] = val
	}

	switch val.Kind() {

	case reflect.Struct:
		//We always enter here first since we want to unmarshel a struct
		t := val.Type()
		for i := 0; i < val.NumField(); i++ {
			structF := t.Field(i)
			if v := val.Field(i); v.CanSet() {
				if err := readData(reader, order, structF, v, storedValues); err != nil {
					return err
				}
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

func test() {

}

func unmarshalArray(reader io.Reader, order binary.ByteOrder, storedValues map[string]reflect.Value, field reflect.StructField, value reflect.Value) error {

	if v, ok := field.Tag.Lookup("byteSize"); ok {

		var size int

		if mapValue, ok := storedValues[v]; ok {
			size = int(mapValue.Uint())
		}

		//Do we need to loop value based on struct
		data := make([]byte, size)

		switch field.Type.Elem().Kind() {

		case reflect.String:
			return fmt.Errorf("does not support type with array: %s ", field.Type.Elem().Kind())

		case reflect.Struct:
			if err := binary.Read(reader, order, &data); err != nil {
				return err
			}

			slice := reflect.MakeSlice(value.Type(), size, size)

			sliceReader := bytes.NewBuffer(data)
			index := 0
			for ; sliceReader.Len() != 0; index++ {
				sliceStruct := reflect.New(value.Type().Elem()).Elem()
				t := value.Type().Elem()
				for i := 0; i < t.NumField(); i++ {
					structF := t.Field(i)
					if v := sliceStruct.Field(i); v.CanSet() {
						if err := readData(sliceReader, order, structF, v, storedValues); err != nil {
							return err
						}
					}

				}
				v := slice.Index(index)
				v.Set(sliceStruct)
			}
			value.Set(slice.Slice(0, index))
			break
		default:
			if err := binary.Read(reader, order, &data); err == nil {
				value.SetBytes(data)
			}
			break

		}
	} else {
		return fmt.Errorf("Field: %s is missing byteSize tag", field.Name)
	}

	return nil
}

func unmarshalString(reader io.Reader, order binary.ByteOrder, storedValues map[string]reflect.Value, field reflect.StructField, value reflect.Value) error {
	if v, ok := field.Tag.Lookup("byteSize"); ok {
		var size int

		if mapValue, ok := storedValues[v]; ok {
			size = int(mapValue.Uint())
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

	var buf bytes.Buffer

	handled := make(map[string]interface{}, 0)

	for i := 0; i < val.NumField(); i++ {

		f := val.Field(i)
		t := ty.Field(i)

		if v, ok := t.Tag.Lookup("bytePayload"); ok {
			field := val.FieldByName(v)

			handled[v] = 0

			unknownSizeData, err := writeByteData(order, field.Kind(), field)

			if err != nil {
				return nil, err
			}

			length, err := intCaster(len(unknownSizeData), t.Type.Kind())

			if err != nil {
				return nil, err
			}

			unknownDataLength, err := writeByteData(order, t.Type.Kind(), reflect.ValueOf(length))

			if err != nil {
				return nil, err
			}

			if err := binary.Write(&buf, order, unknownDataLength); err != nil {
				return nil, err
			}

			if err := binary.Write(&buf, order, unknownSizeData); err != nil {
				return nil, err
			}

		} else if handled[t.Name] == nil {
			data, err := writeByteData(order, t.Type.Kind(), f)
			if err != nil {
				return nil, err
			}
			if err := binary.Write(&buf, order, data); err != nil {
				return nil, err
			}
		}

	}

	return buf.Bytes(), nil
}

func intCaster(value int, to reflect.Kind) (interface{}, error) {
	switch to {
	case reflect.Int:
		return value, nil
	case reflect.Int8:
		return int8(value), nil
	case reflect.Int16:
		return int16(value), nil
	case reflect.Int32:
		return int32(value), nil
	case reflect.Int64:
		return int64(value), nil
	case reflect.Uint:
		return uint(value), nil
	case reflect.Uint8:
		return uint8(value), nil
	case reflect.Uint16:
		return uint16(value), nil
	case reflect.Uint32:
		return uint32(value), nil
	case reflect.Uint64:
		return uint64(value), nil
	case reflect.Float32:
		return float32(value), nil
	case reflect.Float64:
		return float64(value), nil
	}

	return nil, fmt.Errorf("you cant cast int to type : %s ", to)
}

func writeByteData(order binary.ByteOrder, kind reflect.Kind, f reflect.Value) ([]byte, error) {
	var buf bytes.Buffer
	switch kind {
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
		for i := 0; i < f.Len(); i++ {

			if reflect.String == f.Index(i).Kind() {
				return nil, fmt.Errorf("cant handle string slices")
			}

			data, err := writeByteData(order, f.Index(i).Kind(), f.Index(i))

			if err != nil {
				return nil, err
			}
			if err := binary.Write(&buf, order, data); err != nil {
				return nil, err
			}
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
		return nil, fmt.Errorf("does not support type : %s ", kind)
	}

	return buf.Bytes(), nil
}
