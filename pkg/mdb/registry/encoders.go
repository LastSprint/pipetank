package registry

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var errWrongType = errors.New("wrong type")

type uuidCodec struct{}

func (u uuidCodec) DecodeValue(
	_ bson.DecodeContext,
	reader bson.ValueReader,
	value reflect.Value,
) error {
	err := u.tryDecodeString(reader, value)
	if err == nil {
		return nil
	}

	if !errors.Is(err, errWrongType) {
		return bson.ValueDecoderError{
			Name:     "uuidCodec failed: " + err.Error(),
			Types:    []reflect.Type{reflect.TypeFor[uuid.UUID]()},
			Kinds:    []reflect.Kind{reflect.Struct},
			Received: value,
		}
	}

	err = u.tryDecodeBinary(reader, value)
	if err == nil {
		return nil
	}

	return bson.ValueDecoderError{
		Name:     "uuidCodec failed: " + err.Error(),
		Types:    []reflect.Type{reflect.TypeFor[uuid.UUID]()},
		Kinds:    []reflect.Kind{reflect.Struct},
		Received: value,
	}
}

func (u uuidCodec) tryDecodeString(reader bson.ValueReader, value reflect.Value) error {
	strVal, err := reader.ReadString()
	if err != nil {
		if strings.Contains(err.Error(), "but attempted to read %s") {
			return errors.Join(err, errWrongType)
		}

		return err
	}

	val, err := uuid.Parse(strVal)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(val))
	return nil
}

func (u uuidCodec) tryDecodeBinary(reader bson.ValueReader, value reflect.Value) error {
	bval, bytesType, err := reader.ReadBinary()
	if err != nil {
		if strings.Contains(err.Error(), "but attempted to read %s") {
			return errors.Join(err, errWrongType)
		}

		return err
	}

	switch bytesType {
	case bson.TypeBinaryUUID:
		if len(bval) != len(uuid.Nil) {
			return fmt.Errorf("wrong binary length: %d", len(bval))
		}

		val, err := uuid.Parse(string(bval))
		if err != nil {
			return err
		}

		value.Set(reflect.ValueOf(val))
		return nil
	}

	return fmt.Errorf("wrong binary type: %v", bytesType)
}

func (u uuidCodec) EncodeValue(
	_ bson.EncodeContext,
	writer bson.ValueWriter,
	value reflect.Value,
) error {
	if value.IsZero() || value.IsNil() {
		return writer.WriteNull()
	}

	err := u.tryEncode(writer, value)
	if err == nil {
		return nil
	}

	return bson.ValueEncoderError{
		Name:     "uuidCodec failed: " + err.Error(),
		Types:    []reflect.Type{reflect.TypeFor[uuid.UUID]()},
		Kinds:    []reflect.Kind{reflect.Array, reflect.String},
		Received: value,
	}
}

func (u uuidCodec) tryEncode(writer bson.ValueWriter, value reflect.Value) error {
	switch value.Kind() { //nolint:exhaustive
	case reflect.String:
		return writer.WriteString(value.String())
	case reflect.Array:
		if value.Len() != len(uuid.Nil) {
			return fmt.Errorf("wrong array length: %d", value.Len())
		}

		v, ok := value.Interface().([]byte)
		if !ok {
			vu, ok := value.Interface().(uuid.UUID)
			if !ok {
				return fmt.Errorf("wrong array type: %v", value.Kind())
			}

			return writer.WriteBinaryWithSubtype(vu[:], bson.TypeBinaryUUID)
		}

		return writer.WriteBinaryWithSubtype(v, bson.TypeBinaryUUID)
	default:
		return fmt.Errorf("wrong type: %v", value.Kind())
	}
}
