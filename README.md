# bsonparser

bsonparser is bson to json parser written in Go. It converts bson string into MongoDB json extended syntax and the other way around.

## Installation
```
go get -u github.com/naponmeka/bsonparser
```

## Usage

bson -> json
```
jsonStr, err := bsonparser.BsonToJson(bsonStr)
```

json -> bson
```
bsonStr, err := bsonparser.JsonToBson(jsonStr)
```

## Example
```
bsonStr := `
{
    "numbers": [
        5,
        6,
        7
    ],
    "x": {
        "y": "z"
    },
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
jsonStr, err := bsonparser.BsonToJson(bsonStr)
if err != nil {
    log.Fatal(err)
}
fmt.Println(jsonStr)
===output===
{
    "_id": {
        "$id": "5c91e115214fc660b6ca650a"
    },
    "arr": [
        {
            "name": "bone"
        },
        {
            "name": "napon"
        }
    ],
    "name": "pi5",
    "numbers": [
        5,
        6,
        7
    ],
    "value": 3.14159,
    "x": {
        "y": "z"
    }
}
```

```
jsonStr := `
{
    "arr":[{"name":"bone"},{"name":"napon"}],
    "name": "pi5",
    "numbers": [5,6,7],
    "value": 3.14159,
    "x": {
        "y": "z"
    },
    "_id": {
        "$oid": "5c91e115214fc660b6ca650a"
    },
}`
bsonStr, err := bsonparser.JsonToBson(jsonStr)
if err != nil {
    log.Fatal(err)
}
fmt.Println(bsonStr)
===output===
{
    "arr": [
        {"name": "bone"},
        {"name": "napon"}
    ],
    "name": "pi5",
    "numbers": [
        5,
        6,
        7
    ],
    "value": 3.141590,
    "x": {
        "y": "z"
    },
    "_id": ObjectId("5c91e115214fc660b6ca650a")
}

```

