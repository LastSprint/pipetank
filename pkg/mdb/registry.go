package mdb

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

func JSONtoBSON(input []byte) (bson.Raw, error) {
	if len(input) == 0 {
		return nil, nil
	}

	var vl any
	err := bson.UnmarshalExtJSON(input, false, &vl)
	if err != nil {
		return nil, err
	}

	raw, err := bson.Marshal(vl)
	if err != nil {
		return nil, err
	}

	return raw, nil
}
