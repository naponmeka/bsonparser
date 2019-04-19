package bsonparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/naponmeka/bsonparser/bsonparse"
)

// BsonToJson converts MongoDB "Shell mode" into "Strict mode" MongoDB extended json representation
func BsonToJson(bsonStr string) (output string, err error) {
	return BsonToJsonIndent(bsonStr, "", "    ")
}

// BsonToJsonIndent is like JsonToBson but applies Indent to format the output. Each JSON element in the output will begin on a new line beginning with prefix followed by one or more copies of indent according to the indentation nesting.
func BsonToJsonIndent(bsonStr, prefix, indent string) (output string, err error) {
	var out bytes.Buffer
	bsonStr = strings.TrimSpace(bsonStr)
	isArray := strings.HasPrefix(bsonStr, "[")
	if !isArray {
		bsonStr = fmt.Sprintf("[%s]", bsonStr)
	}
	res, err := bsonparse.Parse([]byte(bsonStr))
	if err != nil {
		return "", err
	}
	var outputB []byte
	if !isArray {
		outputB, err = json.Marshal(&res[0])
	} else {
		outputB, err = json.Marshal(&res)
	}
	err = json.Indent(&out, outputB, prefix, indent)
	return out.String(), err
}
