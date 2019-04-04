package bsonparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/naponmeka/bsonparser/bsonparse"
)

func BsonToJson(bsonStr string) (output string, err error) {
	return BsonToJsonIndent(bsonStr, "", "    ")
}

func BsonToJsonIndent(bsonStr, prefix, indent string) (output string, err error) {
	var out bytes.Buffer
	bsonStr = strings.TrimSpace(bsonStr)
	isArray := strings.HasPrefix(bsonStr, "[")
	if !isArray {
		bsonStr = fmt.Sprintf("[%s]", bsonStr)
	}
	res, err := bsonparse.Parse([]byte(bsonStr))
	var outputB []byte
	if !isArray {
		outputB, err = json.Marshal(&res[0])
	} else {
		outputB, err = json.Marshal(&res)
	}
	err = json.Indent(&out, outputB, prefix, indent)
	return out.String(), err
}
