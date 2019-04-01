package main

import (
	"fmt"
	"log"

	"github.com/naponmeka/bsonparser"
)

func main() {
	jsonExt := `
	{
		"_id": {
			"$oid": "5c91e115214fc660b6ca650a"
		},
		"numbers": [5,6,7],
		"arr":[{"name":"bone"},{"name":"napon"}],
		"r" :{
			"$ref": "<name>",
			"$id": "<id>"
		},
		"x": {
			"y": "z"
		},
		"a": {
			"b": {
				"c": "d"
			}
		},
		"nickname": "p",
		"value": 3.14159,
		"name": "pi5"
	  }`
	bsonStr, err := bsonparser.JsonToBson(jsonExt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(bsonStr)
	fmt.Println("==========")
	bsonStrInput := `
	{
		"numbers": [
		  5,
		  6,
		  7
		],
		"x": {
		  "y": "z"
		},
		"a": {
		  "b": {
			"c": "d"
		  }
		},
		"nickname": "p",
		"value": 3.141590,
		"_id": ObjectId("5c91e115214fc660b6ca650a"),
		"arr": [
		  {
			"name": "bone"
		  },
		  {
			"name": "napon"
		  }
		],
		"name": "pi5"
	  }`
	jsonStrOut, err := bsonparser.BsonToJson(bsonStrInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(jsonStrOut)

}
