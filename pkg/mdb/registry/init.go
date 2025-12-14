package registry

import (
	"reflect"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func CreateRegistry() *bson.Registry {
	newReg := bson.NewMgoRegistry()

	uuidType := reflect.TypeFor[uuid.UUID]()
	uuidCodec := uuidCodec{}

	newReg.RegisterTypeDecoder(uuidType, uuidCodec)
	newReg.RegisterTypeEncoder(uuidType, uuidCodec)

	return newReg
}
