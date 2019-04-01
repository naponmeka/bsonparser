package bsonparser

import (
	"bytes"
	"encoding/json"
)

func BsonToJson(bsonStr string) (output string, err error) {
	return BsonToJsonIndent(bsonStr, "", "  ")
}

func BsonToJsonIndent(bsonStr, prefix, indent string) (output string, err error) {
	var out bytes.Buffer
	res, err := Parse([]byte(bsonStr))
	b, err := json.Marshal(&res)
	err = json.Indent(&out, b, prefix, indent)
	return out.String(), err
}
