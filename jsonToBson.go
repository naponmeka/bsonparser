package bsonparser

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func castToString(x interface{}) string {
	switch val := x.(type) {
	case bson.M:
		return ""
	case float64:
		s := ""
		if val == float64(int64(val)) {
			s = fmt.Sprintf("%d", int(val))
		} else {
			s = fmt.Sprintf("%f", val)
		}
		return s
	default:
		return fmt.Sprint(`"`, x, `"`)
	}
}

func traverse(
	removeKey bool,
	current interface{},
	parent interface{},
	level int,
	prefix string,
	indent string,
) (bool, string) {
	needBracket := false
	nl := "\n"
	if prefix == "" && indent == "" {
		nl = ""
	}
	level++
	indents := strings.Repeat(indent, level-1)
	switch val := current.(type) {
	case []interface{}:
		results := []string{}
		res := ""
		for k, v := range val {
			needBracket, res = traverse(true, v, k, level, prefix, indent)
			results = append(results, res)
		}
		if needBracket {
			newResult := []string{}
			for _, r := range results {
				newResult = append(newResult, prefix+indents+indents+"{"+strings.TrimSpace(r)+"}")
			}
			results = newResult
		}
		glue := "," + nl
		childOutStr := strings.Join(results, glue)
		return false, prefix + indents +
			fmt.Sprintf("%s: [%s%s", castToString(parent), nl, childOutStr) +
			nl + prefix + indents + "]"
	case map[string]interface{}:
		results := []string{}
		hasRef := false
		res := ""
		for k, v := range val {
			if castToString(k) != `"$ref"` && castToString(k) != `"$id"` {
				needBracket, res = traverse(false, v, k, level, prefix, indent)
				results = append(results, res)
			} else {
				hasRef = true
			}
		}
		if hasRef {
			results = append(results, fmt.Sprintf(`DBRef(%s, %s)`, castToString(val["$ref"]), castToString(val["$id"])))
		}
		glue := "," + nl
		childOutStr := strings.Join(results, glue)
		if level == 1 {
			return false, fmt.Sprintf("{%s%s%s}", nl, childOutStr, nl)
		} else if removeKey {
			return true, prefix + indents + fmt.Sprintf("%s", childOutStr)
		} else {
			if hasRef {
				return false, prefix + indents + castToString(parent) + ": " + childOutStr
			} else if needBracket {
				return true, prefix + indents +
					castToString(parent) + ": {" + nl + childOutStr + nl + indents + "}"
			} else {
				return false, prefix + indents + castToString(parent) + ": " + childOutStr
			}
		}
	default:
		if castToString(parent) == `"$oid"` {
			return false, fmt.Sprintf("ObjectId(%s)", castToString(val))
		} else if castToString(parent) == `"$numberLong"` {
			return false, fmt.Sprintf(`NumberLong(%s)`, castToString(val))
		} else if castToString(parent) == `"$numberDecimal"` {
			return false, fmt.Sprintf(`NumberDecimal(%s)`, castToString(val))
		} else if castToString(parent) == `"$undefined"` {
			return false, "undefined"
		} else if castToString(parent) == `"$minKey"` {
			return false, "MinKey"
		} else if castToString(parent) == `"$maxKey"` {
			return false, "MaxKey"
		} else if removeKey {
			return false, prefix + indents + fmt.Sprintf(`%s`, castToString(val))
		} else if level == 2 {
			return false, prefix + indents + fmt.Sprintf(`%s: %s`, castToString(parent), castToString(val))
		}
		return true, prefix + indents + fmt.Sprintf("%s: %s", castToString(parent), castToString(val))
	}
}

func JsonToBsonIndent(jsonStr, prefix, indent string) (output string, err error) {
	b := []byte(jsonStr)
	var v interface{}
	err = json.Unmarshal(b, &v)
	_, output = traverse(false, v, v, 0, prefix, indent)
	return
}

func JsonToBson(jsonStr string) (output string, err error) {
	return JsonToBsonIndent(jsonStr, "", "    ")
}
