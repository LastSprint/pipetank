package mdb

import (
	"bytes"
	"encoding/json"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func JSONtoBSON(input json.RawMessage) (bson.Raw, error) {
	if len(input) == 0 {
		return nil, nil
	}

	var vl any
	err := bson.UnmarshalExtJSON(input, false, &vl)
	if err != nil {
		return nil, err
	}

	raw, err := MarshalBson(vl)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func BsonToJSON(input bson.Raw) (json.RawMessage, error) {
	if len(input) == 0 {
		return nil, nil
	}

	var vl any
	err := UnmarshalBson(input, &vl)
	if err != nil {
		return nil, err
	}

	return bson.MarshalExtJSON(vl, false, true)
}

func MarshalBson(input any) ([]byte, error) {
	w := bytes.NewBuffer(make([]byte, 0))
	vw := bson.NewDocumentWriter(w)

	enc := bson.NewEncoder(vw)
	enc.SetRegistry(_registry)

	return w.Bytes(), enc.Encode(input)
}

func UnmarshalBson(input []byte, output any) error {
	r := bson.NewDocumentReader(bytes.NewBuffer(input))
	dec := bson.NewDecoder(r)
	dec.SetRegistry(_registry)

	return dec.Decode(output)
}
