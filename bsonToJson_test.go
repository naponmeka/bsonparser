package bsonparser

import (
	"fmt"
	"log"
)

func ExampleBsonToJson() {
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
	jsonStrOut, err := BsonToJson(bsonStrInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(jsonStrOut)
	// Output:{
	//     "_id": {
	//         "$oid": "5c91e115214fc660b6ca650a"
	//     },
	//     "a": {
	//         "b": {
	//             "c": "d"
	//         }
	//     },
	//     "arr": [
	//         {
	//             "name": "bone"
	//         },
	//         {
	//             "name": "napon"
	//         }
	//     ],
	//     "name": "pi5",
	//     "nickname": "p",
	//     "numbers": [
	//         5,
	//         6,
	//         7
	//     ],
	//     "value": 3.14159,
	//     "x": {
	//         "y": "z"
	//     }
	// }

}
