package bsonparser

import (
	"fmt"
	"log"
)

func ExampleJsonToBson() {
	jsonExt := `
	{
		"nickname": "p",
		"value": 3.14159,
		"name": "pi5",
		"r" :{
			"$ref": "<name>",
			"$id": "<id>"
		},
		"raw": { "$binary": "<bindata>", "$type": "<t>" },
		"reg": { "$regex": "<\"hello\"sRegex>", "$options": "<sOptions>" },
		"_id": {
			"$oid": "5c91e115214fc660b6ca650a"
		},
		"numbers": [5,6,7],
		"arr":[{"name":"bone"},{"name":"napon"}],
		"x": {
			"y": "z"
		},
		"a": {
			"b": {
				"c": "d"
			}
		},
		"meta" : {
			"country" : "TH",
			"reply_to_user_id" : "",
			"user_mention" : [],
			"account_type" : "user",
			"retweet_count" : 0,
			"reply_to_status_id" : "",
			"source" : "Facebook",
			"type" : "tweet",
			"favorite_count" : 0
		}
	}`
	bsonStr, err := JsonToBson(jsonExt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(bsonStr)
	// Output: {
	//"_id": ObjectId("5c91e115214fc660b6ca650a"),
	// "a": {
	// 	"b": {
	// 		"c": "d"
	// 	}
	// },
	// "arr": [
	// 	{"name": "bone"},
	// 	{"name": "napon"}
	// ],
	// "meta": {
	// 	"account_type": "user",
	// 	"country": "TH",
	// 	"favorite_count": 0,
	// 	"reply_to_status_id": "",
	// 	"reply_to_user_id": "",
	// 	"retweet_count": 0,
	// 	"source": "Facebook",
	// 	"type": "tweet",
	// 	"user_mention": [

	// 	]
	// },
	// "name": "pi5",
	// "nickname": "p",
	// "numbers": [
	// 	5,
	// 	6,
	// 	7
	// ],
	// "r": DBRef("<name>", "<id>"),
	// "raw": BinData(<t>, <bindata>),
	// "reg": /<"hello"sRegex>/<sOptions>,
	// "value": 3.141590,
	// "x": {
	// 	"y": "z"
	// }
	// }

}
