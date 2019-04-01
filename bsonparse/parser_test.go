package bsonparse

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	testcases := []struct {
		input   string
		output  map[string]interface{}
		wantErr string
	}{{
		input:  `{}`,
		output: map[string]interface{}{},
	}, {
		input: `{"a": 1}`,
		output: map[string]interface{}{
			"a": float64(1),
		},
	}, {
		input: `{"a": 1, "_id": ObjectId("5c99f90cf1c077b8fbb76089")}`,
		output: map[string]interface{}{
			"a":   float64(1),
			"_id": map[string]interface{}{"$id": "5c99f90cf1c077b8fbb76089"},
		},
	}, {
		input: `{"a": 1, "date": ISODate("xxxx")}`,
		output: map[string]interface{}{
			"a":    float64(1),
			"date": map[string]interface{}{"$date": "xxxx"},
		},
	}, {
		input: `{"a": 1, "age": undefined}`,
		output: map[string]interface{}{
			"a":   float64(1),
			"age": map[string]interface{}{"$undefined": true},
		},
	}, {
		input: `{"a": 1, "age": MinKey}`,
		output: map[string]interface{}{
			"a":   float64(1),
			"age": map[string]interface{}{"$minKey": true},
		},
	}, {
		input: `{"a": 1, "age": MaxKey}`,
		output: map[string]interface{}{
			"a":   float64(1),
			"age": map[string]interface{}{"$maxKey": true},
		},
	}, {
		input: `{"a": 1, "b": ["c", 2]}`,
		output: map[string]interface{}{
			"a": float64(1),
			"b": []interface{}{"c", float64(2)},
		},
	}, {
		input: `{"a": []}`,
		output: map[string]interface{}{
			"a": []interface{}{},
		},
	}, {
		input: `{"a": [1.2]}`,
		output: map[string]interface{}{
			"a": []interface{}{float64(1.2)},
		},
	}, {
		input: `{"a": [1.2, 2.3]}`,
		output: map[string]interface{}{
			"a": []interface{}{float64(1.2), float64(2.3)},
		},
	}, {
		input: `{"a": true, "b": false, "c": null}`,
		output: map[string]interface{}{
			"a": true,
			"b": false,
			"c": nil,
		},
	},
	// {
	// 	input:   `.1`,
	// 	wantErr: `syntax error`,
	// }, {
	// 	input:   `invalid`,
	// 	wantErr: `syntax error`,
	// }
	}
	for _, tc := range testcases {
		got, err := Parse([]byte(tc.input))
		var gotErr string
		if err != nil {
			gotErr = err.Error()
		}
		if gotErr != tc.wantErr {
			t.Errorf(`%s err: %v, want %v`, tc.input, gotErr, tc.wantErr)
		}
		if !reflect.DeepEqual(got, tc.output) {
			t.Errorf(`%s: %#v want %#v`, tc.input, got, tc.output)
		}
	}
}
